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
	filepath string
	file     *os.File
	lexer    *lexer.Lexer
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

	return &Parser{
		filepath: file_path,
		file:     file,
		lexer:    lexer.New(b),
	}, nil
}

func (p *Parser) Parse() error {
	fmt.Println(p.lexer.Text)
	fmt.Println(strconv.Quote(string(p.lexer.Text)))

	return nil
}
