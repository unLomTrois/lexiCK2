package ck2parser

import (
	"bufio"
	"ck2-parser/internal/app/lexer"
	"encoding/json"
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

	// lookahead, err := lexer.GetNextToken()
	// if err != nil {
	// 	return nil, err
	// }

	return &Parser{
		filepath:  file_path,
		lexer:     lexer,
		lookahead: nil,
	}, nil
}

func (p *Parser) _eat(tokentype lexer.TokenType) (*lexer.Token, error) {
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
		return nil, err
	}
	return token, nil
}

func (p *Parser) Parse() error {
	token_queue := make([]*lexer.Token, 0)
	for {
		token, err := p.lexer.GetNextToken()
		if err != nil {
			return err
		}
		if token == nil {
			break
		}
		token_queue = append(token_queue, token)
	}

	new_file, err := os.Create("./tmp/meta.json")
	if err != nil {
		return err
	}
	defer new_file.Close()

	w := bufio.NewWriter(new_file)
	aJSON, _ := json.MarshalIndent(token_queue, "", "  ")
	w.Write(aJSON)
	w.Flush()

	// kek, err = p._eat(lexer.WORD)
	// fmt.Println(kek)
	// kek, err = p._eat(lexer.EQUALS)
	// fmt.Println(kek)
	// kek, err = p._eat(lexer.START)
	// fmt.Println(kek)
	// kek, err = p._eat(lexer.COMMENT)
	// fmt.Println(kek)

	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))
	// token, _ = p.lexer.GetNextToken()
	// fmt.Println(strconv.Quote(string(token)))

	// p.lexer.GetNextToken()
	// fmt.Println(lookahead)

	return nil
}
