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

func TestTokenizeCommentsAndPositions(t *testing.T) {
	input := "# first\nmodel User {}\n# second\n"
	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) < 2 {
		t.Fatalf("expected tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TokenComment || tokens[0].Value != "# first" {
		t.Fatalf("expected first token comment, got %v %q", tokens[0].Type, tokens[0].Value)
	}
	if tokens[0].Line != 1 || tokens[0].Col != 1 {
		t.Fatalf("expected first comment at 1:1, got %d:%d", tokens[0].Line, tokens[0].Col)
	}
	if tokens[1].Type != TokenModel {
		t.Fatalf("expected model token after comment, got %v", tokens[1].Type)
	}
	foundSecond := false
	for _, token := range tokens {
		if token.Type == TokenComment && token.Value == "# second" {
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
	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var hasOptional, hasLBrack, hasRBrack bool
	for _, token := range tokens {
		switch token.Type {
		case TokenOptional:
			hasOptional = true
		case TokenLBrack:
			hasLBrack = true
		case TokenRBrack:
			hasRBrack = true
		}
	}
	if !hasOptional || !hasLBrack || !hasRBrack {
		t.Fatalf("expected ?, [ and ] tokens, got optional=%v lbrack=%v rbrack=%v", hasOptional, hasLBrack, hasRBrack)
	}
}

func TestTokenizeInvalidCharactersInsideComments(t *testing.T) {
	input := "# comment with $ % ^\nmodel User {}\n"
	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) == 0 || tokens[0].Type != TokenComment {
		t.Fatalf("expected leading comment token, got %v", tokens)
	}
}
