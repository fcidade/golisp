package golisp

import "fmt"

func Debug(program string) error {
	tokens, err := Tokenize(program)
	if err != nil {
		return (err)
	}
	fmt.Println("=== Tokens ===")
	for i, token := range tokens {
		fmt.Printf("Token: % 4d | %- 16s  | %- 32s | %- 16v\n", i, token.TokenType, token.Literal, token.Value)
	}
	fmt.Println("\n=== Expressions ===")
	expressions, err := Parse(tokens)
	if err != nil {
		return (err)
	}
	for i, expression := range expressions {
		fmt.Printf("Expression: % 4d | %- 4T | %s\n", i, expression, expression)
		// fmt.Printf("Expression: % 4d | %- 4T \n", i, expression)
	}
	fmt.Println("\n=== Eval ===")
	err = Eval(expressions)
	if err != nil {
		return (err)
	}
	return nil
}
