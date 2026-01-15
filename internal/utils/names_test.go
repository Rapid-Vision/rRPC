package utils

import "testing"

func TestIdentifierNamePascalCase(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"user_name", "UserName"},
		{"userName", "UserName"},
		{"UserName", "UserName"},
		{"userID", "UserId"},
	}

	for _, tt := range tests {
		got := NewIdentifierName(tt.in).PascalCase()
		if got != tt.want {
			t.Fatalf("PascalCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIdentifierNameSnakeCase(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"user_name", "user_name"},
		{"userName", "user_name"},
		{"UserName", "user_name"},
		{"UserID", "user_id"},
		{"userID", "user_id"},
	}

	for _, tt := range tests {
		got := NewIdentifierName(tt.in).SnakeCase()
		if got != tt.want {
			t.Fatalf("SnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
