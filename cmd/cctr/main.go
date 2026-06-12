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
var classPattern = regexp.MustCompile(`^["]?\[:](\w+):\]["?]$`)

type config struct {
	input       io.Reader
	output      io.Writer
	deleteFlag  bool   // -d delete the target
	squeezeFlag bool   // -s squeeze the target
	target      string // the chars that tr shoudl target
	translation string // the chars tr shoudl translate target into
}

// var specifierFuncMap = map[string]substFuncs{
// 	"alpha": {check: unicode.IsLetter, translate: ToLetter},
// 	"upper": {check: unicode.IsUpper, translate: unicode.ToUpper},
// 	"lower": {check: unicode.IsLower, translate: unicode.ToLower},
// 	"digit": {check: unicode.IsDigit, translate: ToDigit},
// 	"print": {check: unicode.IsPrint, translate: ToPrint},
// 	"punct": {check: unicode.IsPunct, translate: ToPunct},
// 	"space": {check: unicode.IsSpace, translate: ToSpace},
// }

func main() {
	var err error

	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("loading config: %v\n", err)
		os.Exit(1)
	}

	tr, err := buildTranslator(cfg)
	if err != nil {
		fmt.Printf("building translator: %v\n", err)
		os.Exit(1)
	}

	err = cfg.run(tr)
	if err != nil {
		fmt.Printf("translating: %v\n", err)
		os.Exit(1)
	}

}

// buildTranslator uses args classifies expressions, expands ranges, loads functions and returns a func(rune) rune
func buildTranslator(cfg config) (func(rune) rune, error) {
	var tr func(rune) rune

	// Check for function class specifier
	targetRunes, targetFuncName := parseExpression(cfg.target)
	translationRunes, translationFuncName := parseExpression(cfg.translation)

	switch {
	case targetFuncName == "" && translationFuncName == "": // reg to reg
		trMap := make(map[rune]rune)

		for i, r := range targetRunes {
			if len(translationRunes) >= i+1 {
				trMap[r] = translationRunes[i]
			} else {
				trMap[r] = translationRunes[len(translationRunes)-1]
			}
		}
		tr = func(r rune) rune {
			if val, exists := trMap[r]; exists {
				if cfg.deleteFlag {
					return -1
				} else {
					return val
				}
			} else {
				return r
			}
		}
	case targetFuncName == "" && translationFuncName != "": // reg to function
	case targetFuncName != "" && translationFuncName == "": // function to reg
	case targetFuncName != "" && translationFuncName != "": // function to function

	}

	return tr, nil
}

// parseExpression takes a string and evaluates if its a regular,
// range or function specifier. returns a []rune with chars if regular or range, or string with functionName if func. caller can check if functionname is "" to know if its reg/range
func parseExpression(s string) ([]rune, string) {
	var chars []rune

	// Handle function expression
	matches := classPattern.FindStringSubmatch(s)
	// matches is nil if no match
	// matches[0] is the full match, matches[1] os capture group
	if matches != nil {
		return []rune{}, matches[1]
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
		// but not including frist of last.
		// Rune is just an int, so we can increment like this
		r := first[len(first)-1]
		endRune := last[0]

		for r < endRune {
			chars = append(chars, r)
			r++
		}
		chars = append(chars, last...)

		return chars, ""
	}

	// Handle regular
	return []rune(s), ""
}

func (cfg *config) run(tr func(rune) rune) error {
	// Wrap output in buffered writer, gives you WriteRune()
	writer := bufio.NewWriter(cfg.output)
	defer writer.Flush() // Fluesh the buffer, even on unexpected errors

	reader := bufio.NewReader(cfg.input)

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

// loadConfig parses flags and loads the config with flags, target and translation vals.
// Also sets default input and output and initializes the subs cache map.
func loadConfig() (config, error) {
	cfg := config{
		input:  os.Stdin,
		output: os.Stdout,
	}

	flag.BoolVar(&cfg.deleteFlag, "d", false, "delete chosen chars")
	flag.BoolVar(&cfg.squeezeFlag, "s", false, "squeeze chosen chars")

	flag.Parse()
	args := flag.Args()

	switch {
	case len(args) == 1 && (cfg.deleteFlag || cfg.squeezeFlag):
		// if one arg and -d or -s flag, make that arg the target.
		cfg.target = args[0]
	case len(args) < 2:
		// if less that two args and no flags -> error, ask for target + translation, two args
		return config{}, fmt.Errorf("please provide 'tr <target> <translation>': %v", args)
	case len(args) == 2:
		// two args: first is target, second is translation
		cfg.target = args[0]
		cfg.translation = args[1]
	default:
		return config{}, fmt.Errorf("please provide cmd <target> <translation>: %v", args)
	}

	// if -d and -s, error, can't combine thpose
	if cfg.deleteFlag && cfg.squeezeFlag {
		return config{}, fmt.Errorf("cannot combine delete and squeeze: %v", args)

	}
	return cfg, nil
}

// func (cfg *config) translateCmd() {
// 	// check if target and translation is regular, range or function
// 	// and expand range
// 	cfg.targetType, cfg.target = checkExpressionAndExpandRange(cfg.target)
// 	cfg.translationType, cfg.translation = checkExpressionAndExpandRange(cfg.translation)
//
// 	cfg.inputType = determineInputType(cfg.targetType, cfg.translationType)
//
// 	cfg.checkAndLoadExpression()
//
// 	scanner := bufio.NewScanner(cfg.input)
// 	scanner.Split(bufio.ScanLines)
//
// 	firstLine := true
//
// 	for scanner.Scan() {
// 		if !firstLine {
// 			fmt.Fprintln(cfg.output)
// 		} else {
// 			firstLine = false
// 		}
//
// 		line := scanner.Text()
// 		processedLine := ""
//
// 		processedLine = cfg.processRunes(line)
//
// 		fmt.Fprint(cfg.output, processedLine)
// 	}
// }
//
// func (cfg *config) processRunes(line string) string {
// 	scanner := bufio.NewScanner(strings.NewReader(line))
// 	scanner.Split(bufio.ScanRunes)
//
// 	// squeezeMap := make(map[rune]struct{})
// 	var res strings.Builder
//
// 	for scanner.Scan() {
// 		currentRune := []rune(scanner.Text())[0]
//
// 		// check cache first
// 		cachedRune, exists := cfg.subst[currentRune]
// 		if exists && cachedRune != 0 {
// 			if cfg.deleteFlag {
// 				continue
// 			}
// 			res.WriteRune(cachedRune)
// 			// if cfg.squeezeFlag {
// 			// squeezeMap[currentRune] = struct{}{}
// 			// }
// 		} else {
// 			processedRune := cfg.substitute(currentRune)
// 			if processedRune != 0 {
// 				res.WriteRune(processedRune)
// 				// if cfg.squeezeFlag && processedRune != currentRune {
// 				// squeezeMap[currentRune] = struct{}{}
// 				// }
// 			}
// 		}
// 	}
// 	return res.String()
// }
//
// func (cfg *config) checkAndLoadExpression() {
// 	// wanna get rid of this, subst is initialized in loadConfig
// 	// but have to fix testing first
// 	if cfg.subst == nil {
// 		cfg.subst = make(map[rune]rune)
// 	}
//
// 	switch cfg.inputType {
// 	case regToReg:
// 		cfg.subst = loadSubstitutionMap(cfg.target, cfg.translation)
// 		cfg.substitute = cfg.regToReg
// 	case regToFunc:
// 		cfg.subst = loadSubstitutionMap(cfg.target, nil)
//
// 		funcs, err := loadSubstFuncs(cfg.translation)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		cfg.translate = funcs.translate
// 		cfg.translation = nil
// 		cfg.substitute = cfg.regToFunc
// 	case funcToReg:
// 		funcs, err := loadSubstFuncs(cfg.target)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		cfg.check = funcs.check
//
// 		cfg.translationSlice = []rune(cfg.translation)
// 		cfg.target = nil
// 		cfg.substitute = cfg.funcToReg
// 	case funcToFunc:
// 		funcs, err := loadSubstFuncs(cfg.target)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		cfg.check = funcs.check
// 		cfg.target = nil
//
// 		funcs, err = loadSubstFuncs(cfg.translation)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		cfg.translate = funcs.translate
// 		cfg.translation = nil
// 		cfg.substitute = cfg.funcToFunc
// 	}
// }
