package ck2parser

import (
	"ck2-parser/internal/app/lexer"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

	statement := p.StatementList()
	fmt.Printf("statement: %v\n", statement)
	// lit := p.Literal()
	// fmt.Println(lit)
	// lit = p.Literal()
	// fmt.Println(lit)
	// lit = p.Literal()
	// fmt.Println(lit)

	return nil
}

type Statement struct {
	Type       string
	Expression interface{}
}

func (p *Parser) StatementList() []*Statement {
	list := make([]*Statement, 0)

	for {
		if p.lookahead == nil {
			break
		}
		list = append(list, p.Statement().(*Statement))
	}

	return list
}

func (p *Parser) Statement() interface{} {
	switch p.lookahead.Type {
	case lexer.START:
		return &Statement{
			Type:       "BLOCK",
			Expression: "{",
		}
	case lexer.COMMENT:
		return p.CommentStatement()
	default:
		return p.ExpressionStatement()
	}
}

func (p *Parser) CommentStatement() interface{} {
	return &Statement{
		Type:       "CommentStatement",
		Expression: p.CommentLiteral(),
	}
}

func (p *Parser) ExpressionStatement() interface{} {
	b := p.EquationExpression()
	return &Statement{
		Type:       "ExpressionStatement",
		Expression: b,
	}
}

type BinaryExpression struct {
	left     *Literal
	operator string
	right    interface{}
}

func (p *Parser) EquationExpression() *BinaryExpression {
	left := p.Literal()
	operator := p._eat(lexer.EQUALS)

	var right interface{}

	switch p.lookahead.Type {
	case lexer.WORD:
		right = p.Literal()
		return &BinaryExpression{
			left:     left,
			operator: string(operator.Value),
			right:    right,
		}
	case lexer.START:
		right = p._eat(lexer.START)
		return &BinaryExpression{
			left:     left,
			operator: string(operator.Value),
			// !!!todo:!!! add block statement
			right: right,
		}
	default:
		return nil
	}
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
	default:
		panic("Unexpected Literal: " + strconv.Quote(string(p.lookahead.Value)))
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
