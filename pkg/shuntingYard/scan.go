package shuntingYard

import (
	"strings"
	"text/scanner"
)

// Scan splits a given input string by whitespaces.
func Scan(input string) ([]string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))

	var tok rune
	var result = make([]string, 0)
	for tok != scanner.EOF {
		tok = s.Scan()
		value := strings.TrimSpace(s.TokenText())
		if len(value) > 0 {
			result = append(result, s.TokenText())
		}
	}
	return result, nil
}
