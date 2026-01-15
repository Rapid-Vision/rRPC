package lexer

import "testing"

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

	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) == 0 {
		t.Fatalf("expected tokens, got none")
	}
	if tokens[0].Type != TokenModel {
		t.Fatalf("expected first token TokenModel, got %v", tokens[0].Type)
	}
	if tokens[0].Line != 1 || tokens[0].Col != 1 {
		t.Fatalf("expected first token position 1:1, got %d:%d", tokens[0].Line, tokens[0].Col)
	}
}

func TestTokenizeReportsErrorPosition(t *testing.T) {
	input := "model User { name: string$ }\n"
	_, err := NewLexer(input).Tokenize()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	lexErr, ok := err.(LexerError)
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
