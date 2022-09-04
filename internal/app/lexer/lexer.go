package lexer

type Lexer struct {
	Text []byte
}

func New(text []byte) *Lexer {
	return &Lexer{
		Text: text,
	}
}
