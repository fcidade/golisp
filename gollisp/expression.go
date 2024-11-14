package golisp

import (
	"fmt"
	"strings"
)

type Expression interface {
	Eval() error
	// String() string // TODO:
}

// --

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

// --

type FnCall struct {
	Token Token
	Args  []Expression
}

var _ Expression = &FnCall{}

// TODO: Be able to return an expression
type Fn func(args []Expression) error

var FnMap = map[string]Fn{
	"println": func(args []Expression) error {
		fmt.Print("OUTPUT: ")
		argsAsStr := []string{}
		for _, arg := range args {
			argsAsStr = append(argsAsStr, fmt.Sprintf("%s", arg))
		}
		fmt.Println(strings.Join(argsAsStr, " "))
		return nil
	},
	"printf": func(args []Expression) error {
		if len(args) == 0 {
			return fmt.Errorf("printf: missing format string")
		}
		fmt.Print("OUTPUT: ")
		argsAny := []any{}
		if len(args) > 1 {
			for _, arg := range args[1:] {
				argsAny = append(argsAny, arg)
			}
		}
		fmt.Printf(fmt.Sprint(args[0]), argsAny...)
		return nil
	},
	"+": func(args []Expression) error {
		// if len(args) != 2 {
		// 	return fmt.Errorf("expected 2 arguments, got %d", len(args))
		// }
		// switch arg1 := args[0].(type) {
		// 	case *StringExpr:
		// 		switch arg2 := args[1].(type) {
		// 			case *StringExpr:
		// 				args[0] = &StringExpr{Value: arg1.Value + arg2.Value}
		// 				args[1] = nil
		// }
		return nil
	},
}

func (e *FnCall) Eval() error {
	// fmt.Print("OUTPUT: ")
	// if e.Token.Literal == "println" {
	// 	return (&Println{Args: e.Args}).Eval()
	// }
	// if e.Token.Literal == "printf" {
	// 	return (&Printf{Args: e.Args}).Eval()
	// }
	fn, ok := FnMap[e.Token.Literal]
	if !ok {
		return fmt.Errorf("unknown function %q", e.Token.Literal)
	}
	return fn(e.Args)
}

func (e *FnCall) String() string {
	argsAsStr := []string{}
	for _, arg := range e.Args {
		argsAsStr = append(argsAsStr, fmt.Sprintf("%q", arg))
	}
	return fmt.Sprintf("(%s %s)", e.Token.Literal, strings.Join(argsAsStr, " "))
}
