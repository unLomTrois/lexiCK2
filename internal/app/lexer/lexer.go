package lexer

type Lexer struct {
	text string
}

func New(text string) *Lexer {
	return &Lexer{
		text: text,
	}
}
