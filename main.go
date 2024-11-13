package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// Token

type Token struct {
	Literal   string
	TokenType TokenType
	Value     any
}

type TokenType int

func (t TokenType) String() string {
	switch t {
	case TokenTypeIdentifier:
		return "Identifier"
	case TokenTypeString:
		return "String"
	case TokenTypeNumber:
		return "Number"
	default:
		return "Unknown"
	}
}

const (
	TokenTypeIdentifier = iota
	TokenTypeString
	TokenTypeNumber
)

func tokenize(input string) ([]Token, error) {
	start := 0
	cursor := 0

	tokens := []Token{}

	var peek = func() string {
		if cursor >= len(input) {
			return ""
		}
		return string(input[cursor])
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
		AllowedIdentifierStartChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		AllowedIdentifierChars      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	)

	for cursor < len(input) {
		currChar := peek()
		switch {
		case currChar == `'`:
			start = cursor
			fetch()
			for fetch() != `'` {
			}
			tokens = append(tokens, Token{Literal: input[start:cursor], TokenType: TokenTypeString})
		case strings.Contains(AllowedIdentifierStartChars, currChar):
			start = cursor
			for strings.Contains(AllowedIdentifierChars, fetch()) {
			}
			tokens = append(tokens, Token{Literal: input[start:cursor], TokenType: TokenTypeIdentifier})
		case slices.Contains([]string{" ", "(", ")"}, currChar):
			fetch()
		default:
			// fetch()
			panic(fmt.Sprintf("unexpected character %s", currChar))
		}

		// switch currChar {
		// case '(':
		// case strings.Contains("", currChar):
		// }
	}
	tokens = append(tokens, Token{Literal: "EOF", TokenType: TokenTypeIdentifier, Value: io.EOF})
	return tokens, nil
}

// Parse

// Eval

type Expression interface {
	Eval() error
}

type Println struct {
	Message string
}

var _ Expression = &Println{}

func (e *Println) Eval() error {
	fmt.Println(e.Message)
	return nil
}

func parse(tokens []Token) (Expression, error) {
	return nil, nil
}

func eval(e Expression) error {
	if e == nil {
		panic(fmt.Sprintf("unexpected type %T", e))
	}
	return e.Eval()
}

func main() {
	program := "(println 'Hello, city!' 'What is your name?')"
	tokens, err := tokenize(program)
	if err != nil {
		panic(err)
	}
	for i, token := range tokens {
		fmt.Printf("Token: % 16d | %- 16s  | %- 32s | %- 16v\n", i, token.TokenType, token.Literal, token.Value)
	}
	expression, err := parse(tokens)
	if err != nil {
		panic(err)
	}
	eval(expression)
}

/*
Syntax:
(println "Hello, city!")
(pln "Hello, city!")
(pf "Hello, %s!" "city")
(printf "Hello, %q!" "city")
*/
