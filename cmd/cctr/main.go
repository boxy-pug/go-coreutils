// cctr copies the standard input to the standard output with substitution or deletion of selected characters.
// A kind of find and replace for individual characters.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const classPattern = `["]?[:](\w+)[:]["]?`

type config struct {
	input            io.Reader
	output           io.Writer
	deleteFlag       bool          // -d delete the target
	squeezeFlag      bool          // -s squeeze the target
	subst            map[rune]rune // caching layer, an optimization
	target           []rune        // the chars that tr shoudl target
	translation      []rune        // the chars tr shoudl translate target into
	targetType       expressionType
	translationType  expressionType
	translationSlice []rune
	inputType        inputType
	substFuncs
}

type expressionType string

const (
	Regular  expressionType = "regular"
	Range    expressionType = "range"
	Function expressionType = "function"
)

type (
	checkFunc        func(rune) bool
	translateFunc    func(rune) rune
	substitutionFunc func(rune) rune
)

// substFuncs is not really needed, overcomplicates things, will simplify
type substFuncs struct {
	check      checkFunc
	translate  translateFunc
	substitute substitutionFunc
}

type inputType int

const (
	regToReg inputType = iota
	regToFunc
	funcToReg
	funcToFunc
)

var specifierFuncMap = map[string]substFuncs{
	"alpha": {check: unicode.IsLetter, translate: ToLetter},
	"upper": {check: unicode.IsUpper, translate: unicode.ToUpper},
	"lower": {check: unicode.IsLower, translate: unicode.ToLower},
	"digit": {check: unicode.IsDigit, translate: ToDigit},
	"print": {check: unicode.IsPrint, translate: ToPrint},
	"punct": {check: unicode.IsPunct, translate: ToPunct},
	"space": {check: unicode.IsSpace, translate: ToSpace},
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("couldn't load config %v\n", err)
		os.Exit(1)
	}

	cfg.translateCmd()
}

// loadConfig parses flags and loads the config with flags, target and translation vals.
// Also sets default input and output and initializes the subs cache map.
func loadConfig() (config, error) {
	cfg := config{
		subst:  make(map[rune]rune),
		input:  os.Stdin,
		output: os.Stdout,
	}
	cfg.subst = map[rune]rune{}

	flag.BoolVar(&cfg.deleteFlag, "d", false, "delete chosen chars")
	flag.BoolVar(&cfg.squeezeFlag, "s", false, "squeeze chosen chars")

	flag.Parse()
	args := flag.Args()

	switch {
	case len(args) == 1 && (cfg.deleteFlag || cfg.squeezeFlag):
		// if one arg and -d or -s flag, make that arg the target.
		cfg.target = []rune(args[0])
	case len(args) < 2:
		// if less that two args and no flags -> error, ask for target + translation, two args
		return cfg, fmt.Errorf("please provide chars to translate and chars to translate into: %v", args)
	case len(args) == 2:
		// two args: first is target, second is translation
		cfg.target = []rune(args[0])
		cfg.translation = []rune(args[1])
	default:
		return cfg, fmt.Errorf("please provide cmd <target> <translation>: %v", args)
	}

	// should validate the flags better, now what happens when d & s is true?
	// if -d, then make the translation char empty ""
	if cfg.deleteFlag {
		cfg.translation = []rune("")
		cfg.translationSlice = []rune("")
	}
	return cfg, nil
}

func (cfg *config) translateCmd() {
	// check if target and translation is regular, range or function
	// and expand range
	cfg.targetType, cfg.target = checkExpressionAndExpandRange(cfg.target)
	cfg.translationType, cfg.translation = checkExpressionAndExpandRange(cfg.translation)

	cfg.inputType = determineInputType(cfg.targetType, cfg.translationType)

	cfg.checkAndLoadExpression()

	scanner := bufio.NewScanner(cfg.input)
	scanner.Split(bufio.ScanLines)

	firstLine := true

	for scanner.Scan() {
		if !firstLine {
			fmt.Fprintln(cfg.output)
		} else {
			firstLine = false
		}

		line := scanner.Text()
		processedLine := ""

		processedLine = cfg.processRunes(line)

		fmt.Fprint(cfg.output, processedLine)
	}
}

func (cfg *config) processRunes(line string) string {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanRunes)

	// squeezeMap := make(map[rune]struct{})
	var res strings.Builder

	for scanner.Scan() {
		currentRune := []rune(scanner.Text())[0]

		// check cache first
		cachedRune, exists := cfg.subst[currentRune]
		if exists && cachedRune != 0 {
			if cfg.deleteFlag {
				continue
			}
			res.WriteRune(cachedRune)
			// if cfg.squeezeFlag {
			// squeezeMap[currentRune] = struct{}{}
			// }
		} else {
			processedRune := cfg.substitute(currentRune)
			if processedRune != 0 {
				res.WriteRune(processedRune)
				// if cfg.squeezeFlag && processedRune != currentRune {
				// squeezeMap[currentRune] = struct{}{}
				// }
			}
		}
	}
	return res.String()
}

func (cfg *config) checkAndLoadExpression() {
	// wanna get rid of this, subst is initialized in loadConfig
	// but have to fix testing first
	if cfg.subst == nil {
		cfg.subst = make(map[rune]rune)
	}

	switch cfg.inputType {
	case regToReg:
		cfg.subst = loadSubstitutionMap(cfg.target, cfg.translation)
		cfg.substitute = cfg.regToReg
	case regToFunc:
		cfg.subst = loadSubstitutionMap(cfg.target, nil)

		funcs, err := loadSubstFuncs(cfg.translation)
		if err != nil {
			fmt.Println(err)
		}
		cfg.translate = funcs.translate
		cfg.translation = nil
		cfg.substitute = cfg.regToFunc
	case funcToReg:
		funcs, err := loadSubstFuncs(cfg.target)
		if err != nil {
			fmt.Println(err)
		}
		cfg.check = funcs.check

		cfg.translationSlice = []rune(cfg.translation)
		cfg.target = nil
		cfg.substitute = cfg.funcToReg
	case funcToFunc:
		funcs, err := loadSubstFuncs(cfg.target)
		if err != nil {
			fmt.Println(err)
		}
		cfg.check = funcs.check
		cfg.target = nil

		funcs, err = loadSubstFuncs(cfg.translation)
		if err != nil {
			fmt.Println(err)
		}
		cfg.translate = funcs.translate
		cfg.translation = nil
		cfg.substitute = cfg.funcToFunc
	}
}
