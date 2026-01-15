package parser

import (
	"strings"
	"testing"
)

func TestParseSchema(t *testing.T) {
	input := `model GreetingMessage {
    message: string
}

rpc HelloWorld(
    name: string,
    surname: string?,
) GreetingMessage
`

	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(schema.Models))
	}
	if len(schema.RPCs) != 1 {
		t.Fatalf("expected 1 rpc, got %d", len(schema.RPCs))
	}
	if schema.RPCs[0].Name != "HelloWorld" {
		t.Fatalf("unexpected rpc name %q", schema.RPCs[0].Name)
	}
	if len(schema.RPCs[0].Parameters) != 2 {
		t.Fatalf("expected 2 params, got %d", len(schema.RPCs[0].Parameters))
	}
}

func TestParseListAndMapTypes(t *testing.T) {
	input := `model User {}
rpc GetUsers() list[User]
rpc GetUsersByName() map[list[User]]
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.RPCs) != 2 {
		t.Fatalf("expected 2 rpcs, got %d", len(schema.RPCs))
	}
	if schema.RPCs[1].Returns.Kind != TypeMap {
		t.Fatalf("expected map return type, got %v", schema.RPCs[1].Returns.Kind)
	}
}

func TestParseSchemaValidationErrors(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name: "duplicate models",
			input: `model User {
    id: int
}
model User {
    name: string
}
`,
			wantErr: `duplicate model "User"`,
		},
		{
			name: "duplicate rpcs",
			input: `rpc GetUser() User
rpc GetUser() User
`,
			wantErr: `duplicate rpc "GetUser"`,
		},
		{
			name: "duplicate model fields",
			input: `model User {
    id: int
    id: int
}
`,
			wantErr: `model "User" has duplicate field "id"`,
		},
		{
			name: "duplicate rpc params",
			input: `rpc GetUser(
    id: int,
    id: int,
) User
`,
			wantErr: `rpc "GetUser" has duplicate parameter "id"`,
		},
		{
			name: "unknown model field type",
			input: `model User {
    profile: Profile
}
`,
			wantErr: `unknown type "Profile"`,
		},
		{
			name: "unknown rpc return type",
			input: `rpc GetUser() User
`,
			wantErr: `unknown type "User"`,
		},
		{
			name: "unknown rpc param type",
			input: `rpc GetUser(
    id: UserId,
) User
`,
			wantErr: `unknown type "UserId"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}
