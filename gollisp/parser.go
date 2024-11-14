package golisp

import "fmt"

func Parse(tokens []Token) ([]Expression, error) {
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
				// TODO: create a fn mapper
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
					parsed, err := Parse(args)
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
					parsed, err := Parse(args)
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
