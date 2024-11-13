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
	case TokenTypeLParen:
		return "LParen"
	case TokenTypeRParen:
		return "RParen"
	case TokenTypeEOF:
		return "EOF"
	default:
		panic(fmt.Sprintf("unexpected token type %d", t))
	}
}

const (
	TokenTypeIdentifier = iota
	TokenTypeString
	TokenTypeNumber
	TokenTypeLParen
	TokenTypeRParen
	TokenTypeEOF
)

var (
	TokenEOF = Token{Literal: "EOF", TokenType: TokenTypeEOF, Value: io.EOF}
)

func tokenize(input string) ([]Token, error) {
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
			// fetch()
			return nil, (fmt.Errorf("unexpected character %q", currChar))
		}

		// switch currChar {
		// case '(':
		// case strings.Contains("", currChar):
		// }
	}
	tokens = append(tokens, TokenEOF)
	return tokens, nil
}

// Parse

func parse(tokens []Token) ([]Expression, error) {
	expressions := []Expression{}

	cursor := 0
	debug_amountOfPeeksOnTheSameChar := map[Token]int{}
	debug_maxAmountOfPeeksOnTheSameChar := 10

	var peek = func() Token {
		c := tokens[cursor]
		debug_amountOfPeeksOnTheSameChar[c]++
		if debug_amountOfPeeksOnTheSameChar[c] > debug_maxAmountOfPeeksOnTheSameChar {
			panic(fmt.Sprintf("peeked on the same token (%q) %d times", c, debug_amountOfPeeksOnTheSameChar[c]))
		}
		if cursor >= len(tokens) {
			return TokenEOF
		}
		return (tokens[cursor])
	}

	var fetch = func() Token {
		t := peek()
		cursor++
		return t
	}

token_loop:
	for cursor < len(tokens) {
		currToken := peek()
		switch currToken.TokenType {
		case TokenTypeEOF:
			break token_loop
		case TokenTypeLParen:
			fetch()
		case TokenTypeString:
			value, ok := currToken.Value.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type %T", currToken.Value)
			}
			expressions = append(expressions, &StringExpr{Value: value})
			fetch()
		case TokenTypeIdentifier:
			{
				switch currToken.Literal {
				case "println":
					fetch()
					expr := &Println{}
					args := []Token{}
					for peek().TokenType != TokenTypeRParen {
						// expr.Args = append(expr.Args, peek())
						args = append(args, peek())
						// fmt.Println("ADDED: ", peek())
						fetch()
					}
					parsed, err := parse(args)
					if err != nil {
						return nil, err
					}
					expr.Args = parsed
					expressions = append(expressions, expr)
					fetch()
				case "printf":
					fetch()
					expr := &Printf{}
					args := []Token{}
					for peek().TokenType != TokenTypeRParen {
						// expr.Args = append(expr.Args, peek())
						args = append(args, peek())
						// fmt.Println("ADDED: ", peek())
						fetch()
					}
					parsed, err := parse(args)
					if err != nil {
						return nil, err
					}
					expr.Args = parsed
					expressions = append(expressions, expr)
					fetch()
				default:
					return nil, fmt.Errorf("unexpected identifier %q (%+v)", currToken.Literal, currToken)
				}
			}
		default:
			return nil, fmt.Errorf("unexpected token %v", currToken)
		}
	}

	return expressions, nil
}

// Eval

type Expression interface {
	Eval() error
	// String() string // TODO:
}

type StringExpr struct {
	Value string
}

var _ Expression = &StringExpr{}

func (e *StringExpr) Eval() error {
	return nil
}

func (e *StringExpr) String() string {
	return e.Value
}

type Println struct {
	Args []Expression
}

var _ Expression = &Println{}

func (e *Println) Eval() error {
	fmt.Print("OUTPUT: ")
	argsAsStr := []string{}
	for _, arg := range e.Args {
		argsAsStr = append(argsAsStr, fmt.Sprintf("%s", arg))
	}
	fmt.Println(strings.Join(argsAsStr, " "))
	return nil
}

func (e *Println) String() string {
	argsAsStr := []string{}
	for _, arg := range e.Args {
		argsAsStr = append(argsAsStr, fmt.Sprintf("%q", arg))
	}
	return fmt.Sprintf("(println %s)", strings.Join(argsAsStr, " "))
}

type Printf struct {
	Args []Expression
}

var _ Expression = &Printf{}

func (e *Printf) Eval() error {
	fmt.Print("OUTPUT: ")
	args := []any{}
	if len(e.Args) > 1 {
		for _, arg := range e.Args[1:] {
			args = append(args, arg)
		}
	}
	fmt.Printf(fmt.Sprint(e.Args[0]), args...)
	return nil
}

func (e *Printf) String() string {
	argsAsStr := []string{}
	for _, arg := range e.Args {
		argsAsStr = append(argsAsStr, fmt.Sprintf("%q", arg))
	}
	return fmt.Sprintf("(printf %s)", strings.Join(argsAsStr, " "))
}

func eval(e []Expression) error {
	if e == nil {
		return fmt.Errorf("unexpected type %T", e)
	}
	for _, e := range e {
		fmt.Printf("EVAL: % 16T %- 16s\n", e, e)
		err := e.Eval()
		if err != nil {
			return err
		}
	}
	return nil
}

// Main

func main() {
	// program := "(println 'Hello, city!' 'What is your name?')"
	program := `
	(println 'Hello, city!' 'What is your name?')
	(println 'Hello, state!' 'What is your age?')
	(println 'Hi!')
	(printf 'Hi! %s!' 'Country')
	`
	tokens, err := tokenize(program)
	if err != nil {
		panic(err)
	}
	fmt.Println("=== Tokens ===")
	for i, token := range tokens {
		fmt.Printf("Token: % 4d | %- 16s  | %- 32s | %- 16v\n", i, token.TokenType, token.Literal, token.Value)
	}
	fmt.Println("\n=== Expressions ===")
	expressions, err := parse(tokens)
	if err != nil {
		panic(err)
	}
	for i, expression := range expressions {
		fmt.Printf("Expression: % 4d | %- 4T | %s\n", i, expression, expression)
		// fmt.Printf("Expression: % 4d | %- 4T \n", i, expression)
	}
	fmt.Println("\n=== Eval ===")
	err = eval(expressions)
	if err != nil {
		panic(err)
	}
}

/*
Syntax:
(println "Hello, city!")
(pln "Hello, city!")
(pf "Hello, %s!" "city")
(printf "Hello, %q!" "city")
*/

// TODO: CLI baseado no codigo do aws-vault
