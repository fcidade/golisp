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
			fetch()
			expr := &FnCall{Token: currToken}
			args := []Token{}
			for peek().TokenType != TokenTypeRParen {
				args = append(args, peek())
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
			return nil, fmt.Errorf("unexpected token %v", currToken)
		}
	}

	return expressions, nil
}
