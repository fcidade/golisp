package golisp

import "fmt"

func Eval(e []Expression) error {
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
