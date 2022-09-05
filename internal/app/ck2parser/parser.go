package ck2parser

import (
	"ck2-parser/internal/app/lexer"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Parser struct {
	Filepath  string       `json:"filepath"`
	lexer     *lexer.Lexer `json:"-"`
	lookahead *lexer.Token `json:"-"`
	Data      []*Statement `json:"data"`
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
		Filepath:  file_path,
		lexer:     lexer,
		lookahead: nil,
		Data:      nil,
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

	p.Data = p.StatementList()

	// fmt.Println(statement[2].Expression.(*BinaryExpression).right)
	json, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}
	w, err := os.Create("tmp/meta.json")
	if err != nil {
		return err
	}
	w.Write(json)

	return nil
}

type Statement struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (p *Parser) StatementList(opt_stop_lookahead ...lexer.TokenType) []*Statement {
	list := make([]*Statement, 0)

	for {
		if p.lookahead == nil {
			break
		}
		if len(opt_stop_lookahead) > 0 && p.lookahead.Type == opt_stop_lookahead[0] {
			p._eat(lexer.END)
			break
		}

		newitem := p.Statement()
		list = append(list, newitem)
	}

	return list
}

func (p *Parser) Statement() *Statement {
	switch p.lookahead.Type {
	case lexer.START:
		return &Statement{
			Type: "BLOCK",
			Data: "{",
		}
	case lexer.COMMENT:
		return p.CommentStatement()
	default:
		return p.ExpressionStatement()
	}
}

func (p *Parser) CommentStatement() *Statement {
	return &Statement{
		Type: "CommentStatement",
		Data: p.CommentLiteral(),
	}
}

func (p *Parser) ExpressionStatement() *Statement {
	return &Statement{
		Type: "ExpressionStatement",
		Data: p.Expression(),
	}
}

type BinaryExpression struct {
	Type     NodeType    `json:"type"`
	Left     Literal     `json:"left"`
	Operator string      `json:"operator"`
	Right    interface{} `json:"right"`
}

func (p *Parser) Expression() *BinaryExpression {
	left := p.Literal()
	operator := p._eat(lexer.EQUALS)
	var right interface{}

	switch p.lookahead.Type {
	case lexer.WORD:
		right = p.Literal()
		return &BinaryExpression{
			Type:     Property,
			Left:     left,
			Operator: string(operator.Value),
			Right:    right,
		}
	case lexer.START:
		right = p.BlockStatement()
		return &BinaryExpression{
			Type:     Block,
			Left:     left,
			Operator: string(operator.Value),
			Right:    right,
		}
	default:
		return nil
	}
}

func (p *Parser) BlockStatement() []*Statement {
	p._eat(lexer.START)

	if p.lookahead.Type == lexer.END {
		return nil
	} else {
		return p.StatementList(lexer.END)
	}
}

type Literal string

func (p *Parser) Literal() Literal {
	switch p.lookahead.Type {
	case lexer.WORD:
		return p.WordLiteral()
	case lexer.COMMENT:
		return p.CommentLiteral()
	default:
		panic("Unexpected Literal: " + strconv.Quote(string(p.lookahead.Value)))
	}
}

func (p *Parser) WordLiteral() Literal {
	token := p._eat(lexer.WORD)
	return Literal(string(token.Value))
}

func (p *Parser) CommentLiteral() Literal {
	token := p._eat(lexer.COMMENT)
	return Literal(string(token.Value))
}
