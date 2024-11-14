package golisp

import (
	"fmt"
	"strings"
)

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
