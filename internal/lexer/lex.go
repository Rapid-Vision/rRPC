package lexer

const (
	TokenModel = iota
	TokenRpc
	TokenIdentifier
	TokenOptional
	TokenColon
	TokenRoundBracketL
	TokenRoundBracketR
	TokenSquareBracketL
	TokenSquareBracketR
	TokenCurlyBracketL
	TokenCurlyBracketR
)

type Token struct {
	Value string
}

type Lexer struct {
	text string
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		text: text,
	}
}

func (l *Lexer) Tokenize() (tokens []Token) {
	return tokens
}
