package golisp

import (
	"fmt"
	"io"
)

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
