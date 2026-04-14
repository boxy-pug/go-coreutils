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
	subst            map[rune]rune
	deleteFlag       bool
	squeezeFlag      bool
	target           []rune
	translation      []rune
	targetType       expressionType
	translationType  expressionType
	output           io.Writer
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

type substFuncs struct {
	check      checkFunc
	translate  translateFunc
	substitute substitutionFunc
}

type inputType int

const (
	regularToRegular inputType = iota
	regularToFunction
	functionToRegular
	functionToFunction
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
		fmt.Printf("couldn't load config %v", err)
		os.Exit(1)
	}

	cfg.translateCmd()

	// fmt.Println(res)
}

func loadConfig() (config, error) {
	cfg := config{
		subst:  make(map[rune]rune),
		input:  os.Stdin,
		output: os.Stdout,
	}

	flag.BoolVar(&cfg.deleteFlag, "d", false, "delete chosen chars")
	flag.BoolVar(&cfg.squeezeFlag, "s", false, "squeeze chosen chars")

	flag.Parse()
	args := flag.Args()

	switch {
	case len(args) == 1 && (cfg.deleteFlag || cfg.squeezeFlag):
		cfg.target = []rune(args[0])
	case len(args) < 2:
		return cfg, fmt.Errorf("please provide chars to translate and chars to translate into: %v", args)
	case len(args) == 2:
		cfg.target = []rune(args[0])
		cfg.translation = []rune(args[1])
	default:
		return cfg, fmt.Errorf("please provide cmd <target> <translation>: %v", args)
	}

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

		// if cfg.squeezeFlag {
		// _, exists := squeezeMap[currentRune]
		// if exists {
		// continue
		// }
		// }

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
	if cfg.subst == nil {
		cfg.subst = make(map[rune]rune)
	}

	switch cfg.inputType {
	case regularToRegular:
		cfg.subst = loadSubstitutionMap(cfg.target, cfg.translation)
		cfg.substitute = cfg.regToReg
	case regularToFunction:
		cfg.subst = loadSubstitutionMap(cfg.target, nil)

		funcs, err := loadSubstFuncs(cfg.translation)
		if err != nil {
			fmt.Println(err)
		}
		cfg.translate = funcs.translate
		cfg.translation = nil
		cfg.substitute = cfg.regToFunc
	case functionToRegular:
		funcs, err := loadSubstFuncs(cfg.target)
		if err != nil {
			fmt.Println(err)
		}
		cfg.check = funcs.check

		cfg.translationSlice = []rune(cfg.translation)
		cfg.target = nil
		cfg.substitute = cfg.funcToReg
	case functionToFunction:
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
