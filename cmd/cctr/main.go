// cctr copies the standard input to the standard output with substitution or deletion of selected characters.
// A kind of find and replace for individual characters.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// classPattern matches "[:word:]" (quotes are optional), captures word.
// ^            - Start of string
// ["]?         - Optional opening quote (matches " if present)
// \[           - Literal opening bracket [ (escaped because [ is a regex metacharacter)
// :            - Literal colon : (start of class specifier)
// (\w+)        - Capture group 1: the class name (word characters only: a-z, A-Z, 0-9, _)
// :            - Literal colon : (end of class specifier)
// \]           - Literal closing bracket ] (escaped because ] is a regex metacharacter)
// ["]?         - Optional closing quote (matches " if present)
// $            - End of string
var classPattern = regexp.MustCompile(`^["]?\[:(\w+):\]["]?$`)

type config struct {
	in          io.Reader
	out         io.Writer
	deleteFlag  bool   // -d delete the target
	squeezeFlag bool   // -s squeeze the target
	target      string // the chars that tr shoudl target
	translation string // the chars tr shoudl translate target into
}

// This implementation supports ascii only, if not the class -> literal is hard to pull of
// Example : tr [:print:] 123 that class would expand to a lot off chars, and you have to map them all, like
// map first printable char " " -> 1, map second printable "!" -> 2 etc. Fine for ascii but weird for unicode.
var asciiAlpha = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz") // Alpha is all lower and uppercase letters
var asciiUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var asciiLower = []rune("abcdefghijklmnopqrstuvwxyz")
var asciiDigit = []rune("0123456789")
var asciiPrint = []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~") // All printable ascii chars
var asciiPunct = []rune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
var asciiSpace = []rune("\t\n\v\f\r ") // tab 9, newline 10, vertical tab 11, formfeed 12, carriage return 13, space 32

// exprKind is a label to know what kind of expression we're dealing with.
// needed to check for allowed expressions in string2/translation
type exprKind string

const (
	exprRegular exprKind = "regular"
	exprDelete  exprKind = "delete"
	exprSqueeze exprKind = "squeeze"
	exprAlpha   exprKind = "alpha"
	exprUpper   exprKind = "upper"
	exprLower   exprKind = "lower"
	exprDigit   exprKind = "digit"
	exprPrint   exprKind = "print"
	exprPunct   exprKind = "punct"
	exprSpace   exprKind = "space"
	exprUnknown exprKind = "unknown"
)

// main is a thin wrapper around run(), for easier integration testing.
func main() {
	err := run(os.Args, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tr: %v\n", err)
		os.Exit(1)
	}
}

// run parses flags, builds a translator and translates the input. Testable boundary.
func run(args []string, in io.Reader, out io.Writer) error {

	cfg, err := parseFlags(args[1:])
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	tr, err := buildTranslator(cfg)
	if err != nil {
		return fmt.Errorf("building translator: %v", err)
	}

	cfg.in = in
	cfg.out = out

	// translate returns an error, and so does run, so this works.
	return cfg.translate(tr)
}

// parseFlags takes a []string of args (minus the program name itself) and parses the flags,
// target and translation into a config struct.
// Returns error if -d and -s is both set, or there's a missing target or translation.
func parseFlags(args []string) (config, error) {
	cfg := config{}
	fs := flag.NewFlagSet("cctr", flag.ContinueOnError)

	fs.BoolVar(&cfg.deleteFlag, "d", false, "delete chosen chars")
	fs.BoolVar(&cfg.squeezeFlag, "s", false, "squeeze chosen chars")

	if err := fs.Parse(args); err != nil {
		return config{}, fmt.Errorf("parsing args: %v", err)
	}
	parsedArgs := fs.Args()

	// tr has 2 modes: 1 string mode with delete or squeeze, 2 string mode without those flags set
	switch {
	case cfg.deleteFlag || cfg.squeezeFlag:
		// if one arg and -d or -s flag, make that arg the target.
		if len(parsedArgs) == 1 {
			cfg.target = parsedArgs[0]
		} else {
			return config{}, fmt.Errorf("too many args: %q", parsedArgs)
		}
	case len(parsedArgs) < 2:
		// we need two strings (target + translation) bcs neither -d or -s was set
		return config{}, fmt.Errorf("missing operand after %q", parsedArgs[0])
	case len(parsedArgs) == 2 && parsedArgs[1] == "":
		// if two args but last one is empty
		return config{}, fmt.Errorf("string2 must be non-empty")
	case len(parsedArgs) == 2:
		// two args: first is target, second is translation
		cfg.target = parsedArgs[0]
		cfg.translation = parsedArgs[1]
	case len(parsedArgs) > 2:
		return config{}, fmt.Errorf("extra operand(s): %q", parsedArgs[2:])
	default:
		return config{}, fmt.Errorf("unknown operand configuration: %q", parsedArgs)
	}

	// if -d and -s, error, can't combine those
	if cfg.deleteFlag && cfg.squeezeFlag {
		return config{}, fmt.Errorf("cannot combine delete and squeeze: %v", args)

	}
	return cfg, nil
}

// expandExpression takes a string and evaluates if its a regular,
// range or function specifier, and expands it to a []rune.
func expandExpression(s string) ([]rune, exprKind, error) {
	var chars []rune

	// Handle function expression
	matches := classPattern.FindStringSubmatch(s)
	// fmt.Printf("match found for %s: %s", matches[0], matches[1])
	// matches is nil if no match
	// matches[0] is the full match, matches[1] os capture group
	if matches != nil {
		switch matches[1] {
		case "alpha":
			return asciiAlpha, exprAlpha, nil
		case "upper":
			return asciiUpper, exprUpper, nil
		case "lower":
			return asciiLower, exprLower, nil
		case "digit":
			return asciiDigit, exprDigit, nil
		case "print":
			return asciiPrint, exprPrint, nil
		case "punct":
			return asciiPunct, exprPunct, nil
		case "space":
			return asciiSpace, exprSpace, nil
		default:
			return []rune{}, exprUnknown, fmt.Errorf("%q doesnt match any class specifier", matches[1])
		}
	}

	// Handle range expression
	// simplified for now, only support one '-' in range
	parts := strings.Split(s, "-")
	// If there's something before and after the - we treat it as real range
	if len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0 {
		first, last := []rune(parts[0]), []rune(parts[1])
		// Append first part minus last char to chars
		chars = append(chars, first[:len(first)-1]...)

		// We want to add any char from end of first up to
		// but not including first of last.
		// Rune is just an int, so we can increment like this
		r := first[len(first)-1]
		endRune := last[0]

		for r < endRune {
			chars = append(chars, r)
			r++
		}
		chars = append(chars, last...)

		return chars, exprRegular, nil
	}

	// Handle regular
	return []rune(s), exprRegular, nil
}

// buildTranslator uses args classifies expressions, expands ranges, loads functions and returns a func(rune) rune
func buildTranslator(cfg config) (func(rune) rune, error) {

	// expand target
	target, targetExprKind, err := expandExpression(cfg.target)
	if err != nil {
		return nil, fmt.Errorf("expanding target expression %q: %w", cfg.target, err)
	}

	// expanding translation, if deleteFlag then it should just contain -1
	var translation []rune
	var translationExprKind exprKind

	if cfg.deleteFlag {
		// for delete mode use -1 as sentinel value meaning "delete this"
		translation = []rune{-1}
		translationExprKind = exprDelete
	} else if cfg.squeezeFlag {
		// for squeeze mode use -2 as sentinel value meaning "squeeze this"
		translation = []rune{-2}
		translationExprKind = exprSqueeze
	} else {
		if cfg.translation == "" {
			return nil, fmt.Errorf("empty string2")

		}
		translation, translationExprKind, err = expandExpression(cfg.translation)
		if err != nil {
			return nil, fmt.Errorf("expanding translation expression %q: %w", cfg.translation, err)
		}
	}

	// handle invalid class expression combinations
	isValidCombination, err := isValidExprCombo(targetExprKind, translationExprKind)
	if err != nil {
		return nil, err
	}

	if !isValidCombination {
		return nil, fmt.Errorf("expanding translation: %w", err)
	}

	// Build the mapping, for each char in target, map to corresponding char in translation.
	// If len(target) > len(translation) just repeat last char of translation
	trMap := make(map[rune]rune)

	for i, r := range target {
		if len(translation) >= i+1 {
			trMap[r] = translation[i]
		} else {
			trMap[r] = translation[len(translation)-1]
		}
	}
	// return a closure that looks up translation in the map.
	// if it's not there, return original r unchanged
	return func(r rune) rune {
		if val, exists := trMap[r]; exists {
			return val
		} else {
			return r
		}
	}, nil
}

func isValidExprCombo(targetExpr, transExpr exprKind) (bool, error) {
	// All class specifiers in th translation position is illegal except the
	// four possible combos of upper and lower.
	switch transExpr {
	case exprUpper, exprLower:
		if targetExpr == exprUpper || targetExpr == exprLower {
			return true, nil
		} else {
			return false, fmt.Errorf("misaligned or invalid class: %q, %q", targetExpr, transExpr)
		}
	case exprAlpha, exprDigit, exprPrint, exprPunct, exprSpace:
		return false, fmt.Errorf("misaligned or invalid class: %q, %q", targetExpr, transExpr)
	default:
		return true, nil
	}
}

func (cfg *config) translate(tr func(rune) rune) error {
	// Wrap output in buffered writer, gives you WriteRune()
	writer := bufio.NewWriter(cfg.out)
	defer writer.Flush() // Fluesh the buffer, even on unexpected errors

	reader := bufio.NewReader(cfg.in)

	var prevRune rune

	// read one rune at a time, look it up in translator map
	// if it's a squeeze: -2 decide wether to skip it
	// if its a delete -1 just skip it
	// otherwise translate
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		translated := tr(r)

		// if translated is -2 that means r shoudl be squeezed
		if translated == -2 {
			// if r is same as prevRune squeeze it by marking it as -1 for delete
			// compare against original rune bcs in the map theres a sentinel val -2
			if r == prevRune {
				translated = -1
				// if this is the first occurence, print the original rune
			} else {
				translated = r
			}
		}
		// if -1 for delete
		if translated == -1 {
			continue
		}
		// else write the translated rune
		if _, err := writer.WriteRune(translated); err != nil {
			return err
		}
		prevRune = r // set current rune to prevRune
	}
	// return the flush error, if writer failed (disk full etc) we want to know.
	return writer.Flush()
}
