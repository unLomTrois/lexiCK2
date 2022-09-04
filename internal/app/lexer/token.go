package lexer

import (
	"fmt"
	"strconv"
)

type TokenType string

const (
	COMMENT TokenType = "COMMENT"
	WORD    TokenType = "WORD"
	NUMBER  TokenType = "NUMBER"
	NULL    TokenType = "NULL"
	EQUALS  TokenType = "EQUALS"
	START   TokenType = "START"
	END     TokenType = "END"
)

var Spec = map[string]TokenType{
	`^#.+`:                COMMENT,
	`^\w+:?(\w+)?(\.\d)?`: WORD,
	`^\d+`:                NUMBER,
	`^\s+`:                NULL,
	`^=`:                  EQUALS,
	`^{`:                  START,
	`^}`:                  END,
}

type Token struct {
	Type  TokenType `json:"type"`
	Value []byte    `json:"value"`
}

func (t Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}
