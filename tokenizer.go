package main

import(
	"strings"
	"unicode"
)

func analyze(text string) []string{
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

func tokenize(text string) []string{
	return strings.FieldsFunx(text,func(r rune) bool{
		return !unicode.IsLetter(r) && !IsNumber(r)
	})
}
