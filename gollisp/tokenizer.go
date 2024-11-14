package golisp

import (
	"fmt"
	"slices"
	"strings"
)

func Tokenize(input string) ([]Token, error) {
	start := 0
	cursor := 0

	tokens := []Token{}

	debug_amountOfPeeksOnTheSameChar := map[string]int{}
	debug_maxAmountOfPeeksOnTheSameChar := 10

	var peek = func() string {
		c := string(input[cursor])
		debug_amountOfPeeksOnTheSameChar[c]++
		if debug_amountOfPeeksOnTheSameChar[c] > debug_maxAmountOfPeeksOnTheSameChar {
			panic(fmt.Sprintf("peeked on the same char (%q) %d times", c, debug_amountOfPeeksOnTheSameChar[c]))
		}
		if cursor >= len(input) {
			return ""
		}
		return c
	}

	var fetch = func() string {
		if cursor >= len(input) {
			return ""
		}
		o := string(input[cursor])
		cursor++
		return o
	}

	const (
		AllowedIdentifierStartChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_+-*/"
		AllowedIdentifierChars      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+-*/"
	)

	for cursor < len(input) {
		currChar := peek()
		switch {
		case currChar == `'`:
			start = cursor
			fetch()
			for fetch() != `'` {
			}
			tokens = append(tokens, Token{
				Literal:   input[start:cursor],
				Value:     input[start+1 : cursor-1],
				TokenType: TokenTypeString,
			})
		case strings.Contains(AllowedIdentifierStartChars, currChar):
			start = cursor
			for strings.Contains(AllowedIdentifierChars, fetch()) {
			}
			tokens = append(tokens, Token{Literal: input[start : cursor-1], TokenType: TokenTypeIdentifier})
		case currChar == `(`:
			tokens = append(tokens, Token{Literal: currChar, TokenType: TokenTypeLParen})
			fetch()
		case currChar == `)`:
			tokens = append(tokens, Token{Literal: currChar, TokenType: TokenTypeRParen})
			fetch()
		case slices.Contains([]string{" ", "\n", "\t"}, currChar):
			fetch()
		default:
			return nil, (fmt.Errorf("unexpected character %q", currChar))
		}
	}
	tokens = append(tokens, TokenEOF)
	return tokens, nil
}
