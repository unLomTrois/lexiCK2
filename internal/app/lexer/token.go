package lexer

import "fmt"

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
	`^#.+`:    COMMENT,
	`^[A-z]+`: WORD,
	`^\d+`:    NUMBER,
	`^\s+`:    NULL,
	`^=`:      EQUALS,
	`^{`:      START,
	`^}`:      END,
}

type Token struct {
	Type  TokenType `json:"type"`
	Value []byte    `json:"value"`
}

func (t Token) String() string {
	return fmt.Sprintf("type: %v, value: %v", t.Type, string(t.Value))
}
