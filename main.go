package main

import (
	"bufio"
	"fmt"
	"os"

	golisp "github.com/fcidade/golisp/gollisp"
)

// func main() {
// 	program := `
// 	(println 'Hello, city!' 'What is your name?')
// 	(println 'Hello, state!' 'What is your age?')
// 	(println 'Hi!')
// 	(printf 'Hi! %s!' 'Country')
// 	`
// 	err := golisp.Debug(program)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		program, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = golisp.Debug(program)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
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
