package main

import "unicode"

func ToDigit(r rune) rune {
	if unicode.IsDigit(r) {
		return r
	}
	if unicode.IsLetter(r) {
		lowerCaseR := unicode.ToLower(r)
		if lowerCaseR >= 'a' && lowerCaseR <= 'i' {
			return rune('0' + (lowerCaseR - 'a'))
		} else {
			return '9'
		}
	}
	return '9'
}

func ToPunct(r rune) rune {
	if unicode.IsPunct(r) {
		return r
	}
	return '.'
}

func ToSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return r
	}
	return ' '
}

// TODO: how to convert to print?
func ToPrint(r rune) rune {
	if unicode.IsPrint(r) {
		return r
	}
	return 'x'
}

func ToLetter(r rune) rune {
	if unicode.IsLetter(r) {
		return r
	}
	return rune('a' + (r % 26))
}
