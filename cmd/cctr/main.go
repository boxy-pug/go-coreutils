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

// var specifierFuncMap = map[string]substFuncs{
// 	"alpha": {check: unicode.IsLetter, translate: ToLetter},
// 	"upper": {check: unicode.IsUpper, translate: unicode.ToUpper},
// 	"lower": {check: unicode.IsLower, translate: unicode.ToLower},
// 	"digit": {check: unicode.IsDigit, translate: ToDigit},
// 	"print": {check: unicode.IsPrint, translate: ToPrint},
// 	"punct": {check: unicode.IsPunct, translate: ToPunct},
// 	"space": {check: unicode.IsSpace, translate: ToSpace},
// }

// main is a thin wrapper around run(), for easier integration testing.
func main() {
	err := run(os.Args, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "running tr with args %v: %v", os.Args, err)
		os.Exit(1)
	}
}

// run parses flags, builds a translator and translates the input.
func run(args []string, in io.Reader, out io.Writer) error {

	cfg, err := parseFlags(args[1:])
	if err != nil {
		return fmt.Errorf("parsing %v as flags: %v", args[1:], err)
	}

	cfg.in = in
	cfg.out = out

	tr, err := buildTranslator(cfg)
	if err != nil {
		return fmt.Errorf("building translator: %v\n", err)
	}

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

	switch {
	case len(parsedArgs) == 1 && (cfg.deleteFlag || cfg.squeezeFlag):
		// if one arg and -d or -s flag, make that arg the target.
		cfg.target = parsedArgs[0]
	case len(parsedArgs) < 2:
		// if less that two args and no flags -> error, ask for target + translation, two args
		return config{}, fmt.Errorf("please provide 'tr <target> <translation>': %v", parsedArgs)
	case len(parsedArgs) == 2:
		// two args: first is target, second is translation
		cfg.target = parsedArgs[0]
		cfg.translation = parsedArgs[1]
	default:
		return config{}, fmt.Errorf("please provide cmd <target> <translation>: %v", parsedArgs)
	}

	// if -d and -s, error, can't combine those
	if cfg.deleteFlag && cfg.squeezeFlag {
		return config{}, fmt.Errorf("cannot combine delete and squeeze: %v", args)

	}
	return cfg, nil
}

// expandExpression takes a string and evaluates if its a regular,
// range or function specifier, and expands it to a []rune.
func expandExpression(s string) ([]rune, error) {
	var chars []rune

	// Handle function expression
	matches := classPattern.FindStringSubmatch(s)
	// fmt.Printf("match found for %s: %s", matches[0], matches[1])
	// matches is nil if no match
	// matches[0] is the full match, matches[1] os capture group
	if matches != nil {
		switch matches[1] {
		case "alpha":
			return asciiAlpha, nil
		case "upper":
			return asciiUpper, nil
		case "lower":
			return asciiLower, nil
		case "digit":
			return asciiDigit, nil
		case "print":
			return asciiPrint, nil
		case "punct":
			return asciiPunct, nil
		case "space":
			return asciiSpace, nil
		default:
			return []rune{}, fmt.Errorf("%q doesnt match any class specifier", matches[1])
		}
	}

	// Handle range expression
	// simplified for now, only support one - in range
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

		return chars, nil
	}

	// Handle regular
	return []rune(s), nil
}

// buildTranslator uses args classifies expressions, expands ranges, loads functions and returns a func(rune) rune
func buildTranslator(cfg config) (func(rune) rune, error) {

	// expand target
	target, err := expandExpression(cfg.target)
	if err != nil {
		return nil, fmt.Errorf("expanding target expression %q: %w", cfg.target, err)
	}

	// expanding translation, if deleteFlag then it should just contain -1
	var translation []rune
	if cfg.deleteFlag {
		translation = []rune{-1}
	} else {
		translation, err = expandExpression(cfg.translation)
		if err != nil {
			return nil, fmt.Errorf("expanding translation expression %q: %w", cfg.translation, err)
		}
	}

	trMap := make(map[rune]rune)

	for i, r := range target {
		if len(translation) >= i+1 {
			trMap[r] = translation[i]
		} else {
			trMap[r] = translation[len(translation)-1]
		}
	}
	return func(r rune) rune {
		if val, exists := trMap[r]; exists {
			return val
		} else {
			return r
		}
	}, nil
}

func (cfg *config) translate(tr func(rune) rune) error {
	// Wrap output in buffered writer, gives you WriteRune()
	writer := bufio.NewWriter(cfg.out)
	defer writer.Flush() // Fluesh the buffer, even on unexpected errors

	reader := bufio.NewReader(cfg.in)

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		translated := tr(r)
		// -1 means deleted char
		if translated != -1 {
			if _, err := writer.WriteRune(translated); err != nil {
				return err
			}
		}
	}
	// return the flush error, if writer failed (disk full etc) we want to know.
	return writer.Flush()
}
