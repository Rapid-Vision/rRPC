package lexer_test

import (
	"testing"

	"github.com/Rapid-Vision/rRPC/internal/lexer"
)

func TestTokenizeExampleSchema(t *testing.T) {
	input := `model User {
    id: int
    name: string
    surname: string?
}

rpc GetUser(
    userId: int,
) User

rpc ListUsers() list[User]

rpc GetUsernameMap() map[User]
`

	tokens, err := lexer.NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) == 0 {
		t.Fatalf("expected tokens, got none")
	}
	if tokens[0].Type != lexer.TokenModel {
		t.Fatalf("expected first token TokenModel, got %v", tokens[0].Type)
	}
	if tokens[0].Line != 1 || tokens[0].Col != 1 {
		t.Fatalf("expected first token position 1:1, got %d:%d", tokens[0].Line, tokens[0].Col)
	}
}

func TestTokenizeReportsErrorPosition(t *testing.T) {
	input := "model User { name: string$ }\n"
	_, err := lexer.NewLexer(input).Tokenize()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	lexErr, ok := err.(lexer.LexerError)
	if !ok {
		t.Fatalf("expected LexerError, got %T", err)
	}
	if lexErr.Line != 1 || lexErr.Col == 0 {
		t.Fatalf("expected line 1 with column, got %d:%d", lexErr.Line, lexErr.Col)
	}
	if lexErr.Ch != "$" {
		t.Fatalf("expected error for '$', got %q", lexErr.Ch)
	}
}

func TestTokenizeCommentsAndPositions(t *testing.T) {
	input := "# first\nmodel User {}\n# second\n"
	tokens, err := lexer.NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) < 2 {
		t.Fatalf("expected tokens, got %d", len(tokens))
	}
	if tokens[0].Type != lexer.TokenComment || tokens[0].Value != "# first" {
		t.Fatalf("expected first token comment, got %v %q", tokens[0].Type, tokens[0].Value)
	}
	if tokens[0].Line != 1 || tokens[0].Col != 1 {
		t.Fatalf("expected first comment at 1:1, got %d:%d", tokens[0].Line, tokens[0].Col)
	}
	if tokens[1].Type != lexer.TokenModel {
		t.Fatalf("expected model token after comment, got %v", tokens[1].Type)
	}
	foundSecond := false
	for _, token := range tokens {
		if token.Type == lexer.TokenComment && token.Value == "# second" {
			if token.Line != 3 || token.Col != 1 {
				t.Fatalf("expected second comment at 3:1, got %d:%d", token.Line, token.Col)
			}
			foundSecond = true
		}
	}
	if !foundSecond {
		t.Fatalf("expected to find second comment token")
	}
}

func TestTokenizeOptionalAndBrackets(t *testing.T) {
	input := "rpc ListUsers() list[User?]\n"
	tokens, err := lexer.NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var hasOptional, hasLBrack, hasRBrack bool
	for _, token := range tokens {
		switch token.Type {
		case lexer.TokenOptional:
			hasOptional = true
		case lexer.TokenLBrack:
			hasLBrack = true
		case lexer.TokenRBrack:
			hasRBrack = true
		}
	}
	if !hasOptional || !hasLBrack || !hasRBrack {
		t.Fatalf("expected ?, [ and ] tokens, got optional=%v lbrack=%v rbrack=%v", hasOptional, hasLBrack, hasRBrack)
	}
}

func TestTokenizeInvalidCharactersInsideComments(t *testing.T) {
	input := "# comment with $ % ^\nmodel User {}\n"
	tokens, err := lexer.NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) == 0 || tokens[0].Type != lexer.TokenComment {
		t.Fatalf("expected leading comment token, got %v", tokens)
	}
}
