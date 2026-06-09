package main

import (
	"fmt"
	"regexp"
	"strings"
)

func loadSubstitutionMap(target, translation []rune) map[rune]rune {
	res := make(map[rune]rune)

	if len(translation) == 0 {
		for _, r := range target {
			res[r] = 0
		}
		return res
	}

	for i, r := range target {
		if i < len(translation) {
			res[r] = translation[i]
		} else {
			res[r] = translation[len(translation)-1]
		}
	}
	return res
}

func checkExpressionAndExpandRange(ru []rune) (expressionType, []rune) {
	et := checkExpression(ru)

	if et != Range {
		return et, ru
	}

	s := string(ru)
	var res []rune
	idx := strings.Index(s, "-")

	startTarget := ru[idx-1]
	endTarget := ru[idx+1]

	// if startTarget is bigger than endTarget treat as normal subst
	if startTarget > endTarget {
		return et, ru
	}

	var startRest []rune
	var endRest []rune

	// Extract the parts before and after the range
	if idx > 1 {
		startRest = ru[:idx-1]
	}
	if idx < len(s)-2 {
		endRest = ru[idx+2:]
	}

	r := startTarget

	res = append(res, startRest...)

	for r <= endTarget {
		res = append(res, r)
		r++
	}

	res = append(res, endRest...)

	return et, res
}

func checkExpression(ru []rune) expressionType {
	s := string(ru)
	re := regexp.MustCompile(classPattern)

	switch {
	case re.MatchString(s):
		return Function
	case isValidRangeSubstitution(ru):
		return Range
	default:
		return Regular
	}
}

func isValidRangeSubstitution(ru []rune) bool {
	s := string(ru)
	idx := strings.Index(s, "-")
	return idx != -1 && len(s) >= 3 && idx > 0 && idx < len(s)-1
}

func loadSubstFuncs(ru []rune) (substFuncs, error) {
	var sf substFuncs
	s := string(ru)

	re := regexp.MustCompile(classPattern)
	matches := re.FindStringSubmatch(s)

	className := ""

	if len(matches) > 1 {
		className = matches[1]

		if funcs, exists := specifierFuncMap[className]; exists {
			return funcs, nil
		}
	}
	return sf, fmt.Errorf("class specifier not found: %q", className)
}

func determineInputType(targT, transT expressionType) inputType {
	switch {
	case targT == Function && transT != Function:
		return functionToRegular
	case targT != Function && transT == Function:
		return regularToFunction
	case targT == Function && transT == Function:
		return functionToFunction
	default:
		return regularToRegular
	}
}

func (cfg *config) regToReg(r rune) rune {
	val, exists := cfg.subst[r]
	if exists {
		return val
	}
	return r
}

func (cfg *config) regToFunc(r rune) rune {
	_, exists := cfg.subst[r]
	if exists {
		processedRune := cfg.translate(r)
		cfg.subst[r] = processedRune
		return processedRune
	}
	return r
}

func (cfg *config) funcToReg(r rune) rune {
	if cfg.check(r) {
		if cfg.deleteFlag {
			return 0
		}
		if len(cfg.translationSlice) == 0 {
			return 0
		}
		currentReplacementRune := cfg.translationSlice[0]
		cfg.subst[r] = currentReplacementRune
		if len(cfg.translationSlice) > 1 {
			cfg.translationSlice = cfg.translationSlice[1:]
		}
		return currentReplacementRune
	}
	return r
}

func (cfg *config) funcToFunc(r rune) rune {
	if cfg.check(r) {
		return cfg.translate(r)
	}
	return r
}
