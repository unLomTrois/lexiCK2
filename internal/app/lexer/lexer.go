package lexer

import (
	"regexp"
)

type Lexer struct {
	Text   []byte
	cursor int
}

func New(text []byte) *Lexer {
	return &Lexer{
		Text:   text,
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

	_string := l.Text[l.cursor:]

	for k, token_type := range Spec {
		// fmt.Println("try:", k)
		reg := regexp.MustCompile(k)
		token_value := l._match(reg, _string)
		if token_value == nil {
			continue
		}
		if token_type == NULL {
			return l.GetNextToken()
		}
		return &Token{
			Type:  token_type,
			Value: token_value,
		}, nil
	}

	panic("Unexpected token: " + string(_string[0]))
}
