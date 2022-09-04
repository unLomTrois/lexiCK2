package ck2parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Parser struct {
	filepath string
	file     *os.File
}

func New(file *os.File) (*Parser, error) {
	file_path, err := filepath.Abs(file.Name())
	if err != nil {
		return nil, err
	}

	return &Parser{
		filepath: file_path,
		file:     file,
	}, nil
}

func (p *Parser) Parse() error {
	b, err := io.ReadAll(p.file)
	if err != nil {
		return nil
	}
	fmt.Println(b)
	fmt.Println(strconv.Quote(string(b)))

	return nil
}
