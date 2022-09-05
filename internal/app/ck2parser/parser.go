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
	Data      []*Node      `json:"data"`
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

func (p *Parser) Parse() error {
	p.lookahead, _ = p.lexer.GetNextToken()

	p.Data = p.NodeList()

	w, err := os.Create("tmp/meta.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	enc.Encode(p)

	return nil
}

func (p *Parser) NodeList(opt_stop_lookahead ...lexer.TokenType) []*Node {
	nodes := make([]*Node, 0)

	for {
		if p.lookahead == nil {
			break
		}
		if len(opt_stop_lookahead) > 0 && p.lookahead.Type == opt_stop_lookahead[0] {
			p._eat(lexer.END)
			break
		}

		new_node := p.Node()
		nodes = append(nodes, new_node)
	}

	return nodes
}

func (p *Parser) Node() *Node {
	switch p.lookahead.Type {
	case lexer.COMMENT:
		return p.CommentNode()
	default:
		return p.ExpressionNode()
	}
}

func (p *Parser) CommentNode() *Node {
	return &Node{
		Type: Comment,
		Data: p.CommentLiteral(),
	}
}

func (p *Parser) ExpressionNode() *Node {
	key := p.Literal()

	var _type NodeType
	var _operator *lexer.Token
	var _opvalue string
	switch p.lookahead.Type {
	case lexer.EQUALS:
		_operator = p._eat(lexer.EQUALS)
		if string(_operator.Value) == "==" {
			_type = Comparison
		} else {
			_type = Property
		}
		_opvalue = string(_operator.Value)
	case lexer.COMPARISON:
		_operator = p._eat(lexer.COMPARISON)
		_type = Comparison
		_opvalue = string(_operator.Value)
	}

	var value interface{}

	switch p.lookahead.Type {
	case lexer.WORD, lexer.NUMBER:
		value = p.Literal()
		return &Node{
			Type:     _type,
			Key:      key,
			Operator: _opvalue,
			Data:     value,
		}
	case lexer.START:
		value = p.BlockNode()
		return &Node{
			Type:     Block,
			Key:      key,
			Operator: _opvalue,
			Data:     value,
		}
	default:
		return nil
	}
}

func (p *Parser) BlockNode() []*Node {
	p._eat(lexer.START)

	if p.lookahead.Type == lexer.END {
		return nil
	} else {
		return p.NodeList(lexer.END)
	}
}

func (p *Parser) Literal() interface{} {
	switch p.lookahead.Type {
	case lexer.WORD:
		return p.WordLiteral()
	case lexer.NUMBER:
		return p.NumberLiteral()
	case lexer.COMMENT:
		return p.CommentLiteral()
	default:
		panic("[Parser] Unexpected Literal: " + strconv.Quote(string(p.lookahead.Value)) + ", with type of: " + string(p.lookahead.Type))
	}
}

func (p *Parser) WordLiteral() string {
	token := p._eat(lexer.WORD)
	return string(token.Value)
}

func (p *Parser) NumberLiteral() float32 {
	token := p._eat(lexer.NUMBER)
	number, err := strconv.ParseFloat(string(token.Value), 32)
	if err != nil {
		panic("[Parser] Unexpected NumberLiteral: " + strconv.Quote(string(token.Value)))
	}
	return float32(number)
}

func (p *Parser) CommentLiteral() string {
	token := p._eat(lexer.COMMENT)
	return string(token.Value)
}

func (p *Parser) _eat(tokentype lexer.TokenType) *lexer.Token {
	token := p.lookahead
	if token == nil {
		panic("[Parser] Unexpected end of input, expected: " + string(tokentype))
	}
	if token.Type != tokentype {
		panic("[Parser] Unexpected token: \"" + string(token.Value) + "\" with type of " + string(token.Type) + ", expected type: " + string(tokentype))
	}

	var err error
	p.lookahead, err = p.lexer.GetNextToken()
	if err != nil {
		panic(err)
	}
	return token
}
