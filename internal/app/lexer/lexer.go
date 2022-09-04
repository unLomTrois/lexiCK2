package lexer

import (
	"regexp"
	"unicode"
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

func (l *Lexer) GetNextToken() ([]byte, error) {
	if !l.hasMoreTokens() {
		return nil, nil
	}

	_string := l.Text[l.cursor:]

	if unicode.IsLetter(rune(_string[0])) {
		reg, err := regexp.Compile(`^\w+`)
		if err != nil {
			return nil, err
		}
		match := reg.Find(_string)
		l.cursor += len(match)
		return match, nil
	}

	if unicode.In(rune(_string[0]), unicode.Sm) {
		l.cursor++
		return []byte{_string[0]}, nil
	}

	if unicode.IsSpace(rune(_string[0])) {
		l.cursor++
		return []byte{' '}, nil
	}

	return _string, nil
}
