package ck2parser

import (
	"ck2-parser/internal/app/lexer"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Parser struct {
	filepath  string
	lexer     *lexer.Lexer
	lookahead *lexer.Token
}

func New(file *os.File) (*Parser, error) {
	file_path, err := filepath.Abs(file.Name())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lexer := lexer.New(b)

	return &Parser{
		filepath:  file_path,
		lexer:     lexer,
		lookahead: nil,
	}, nil
}

func (p *Parser) _eat(tokentype lexer.TokenType) *lexer.Token {
	token := p.lookahead
	if token == nil {
		panic("Unexpected end of input, expected: " + string(tokentype))
	}
	if token.Type != tokentype {
		panic("Unexpected token: \"" + string(token.Value) + "\" with type of " + string(token.Type) + ", expected type: " + string(tokentype))
	}

	var err error
	p.lookahead, err = p.lexer.GetNextToken()
	if err != nil {
		panic(err)
	}
	return token
}

func (p *Parser) Parse() error {
	p.lookahead, _ = p.lexer.GetNextToken()

	lit := p.Literal()
	fmt.Println(lit)
	lit = p.Literal()
	fmt.Println(lit)
	lit = p.Literal()
	fmt.Println(lit)

	return nil
}

type Literal struct {
	Type  string
	Value string
}

func (p *Parser) Literal() *Literal {
	switch p.lookahead.Type {
	case lexer.WORD:
		return p.WordLiteral()
	case lexer.COMMENT:
		return p.CommentLiteral()
	case lexer.EQUALS:
		fmt.Println(p.lookahead)
		return nil
	default:
		panic("Literal: unexpected literal production")
	}
}

func (p *Parser) WordLiteral() *Literal {
	token := p._eat(lexer.WORD)
	return &Literal{
		Type:  "WordLiteral",
		Value: string(token.Value),
	}
}

func (p *Parser) CommentLiteral() *Literal {
	token := p._eat(lexer.COMMENT)
	return &Literal{
		Type:  "CommentLiteral",
		Value: string(token.Value),
	}
}
