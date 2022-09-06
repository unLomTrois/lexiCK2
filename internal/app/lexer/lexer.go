package lexer

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strconv"
)

type Lexer struct {
	Text    []byte
	cursor  int
	_string []byte
}

func NormalizeText(text []byte) []byte {

	text = bytes.TrimSpace(text)
	// text = bytes.ReplaceAll(text, []byte("    "), []byte("\t"))
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	text = bytes.ReplaceAll(text, []byte("\t\n"), []byte("\n"))
	text = bytes.ReplaceAll(text, []byte("\t"), []byte(""))
	text = bytes.ReplaceAll(text, []byte(" = "), []byte("="))
	text = bytes.ReplaceAll(text, []byte("= {"), []byte("={"))

	// replace \n\n\n.. with \n\n
	reg := regexp.MustCompile(`\n{3,}`)
	text = reg.ReplaceAll(text, []byte("\n\n"))

	return text
}

func New(text []byte) *Lexer {
	normalized := NormalizeText(text)
	// fmt.Println(strconv.Quote(string(normalized)))

	new_file, _ := os.Create("./tmp/normalized.txt")
	defer new_file.Close()

	w := bufio.NewWriter(new_file)
	w.Write(normalized)
	w.Flush()

	return &Lexer{
		Text:   normalized,
		cursor: 0,
	}
}

func (l *Lexer) hasMoreTokens() bool {
	return l.cursor < len(l.Text)
}

func (l *Lexer) isEOF() bool {
	return l.cursor == len(l.Text)
}

func (l *Lexer) _match(reg *regexp.Regexp, text []byte) []byte {
	if match := reg.Find(text); match != nil {
		l.cursor += len(match)
		return match
	}
	return nil
}

func (l *Lexer) GetNextToken() (*Token, error) {
	if !l.hasMoreTokens() {
		return nil, nil
	}

	l._string = l.Text[l.cursor:]

	for k, token_type := range Spec {
		// todo: implement less greedy matching
		// fmt.Println("try:", k, "on: ", string(l._string[0:10]))
		reg := regexp.MustCompile(k)
		token_value := l._match(reg, l._string)
		if token_value == nil {
			// fmt.Println("continue")
			continue
		}
		if token_type == NULL {
			// fmt.Println("null")
			return l.GetNextToken()
		}
		// fmt.Println("return")
		return &Token{
			Type:  token_type,
			Value: token_value,
		}, nil
	}

	panic("[Lexer] Unexpected token: " + strconv.Quote(string(l._string[0])))
}
