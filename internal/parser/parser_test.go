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

func TestParseRPCNoReturn(t *testing.T) {
	input := `rpc Ping()
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.RPCs) != 1 {
		t.Fatalf("expected 1 rpc, got %d", len(schema.RPCs))
	}
	if schema.RPCs[0].HasReturn {
		t.Fatalf("expected no return type")
	}
}

func TestParseRPCNoReturnFollowedByModel(t *testing.T) {
	input := `rpc Ping()
model User {}
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.RPCs) != 1 || len(schema.Models) != 1 {
		t.Fatalf("expected 1 rpc and 1 model, got %d rpcs and %d models", len(schema.RPCs), len(schema.Models))
	}
	if schema.RPCs[0].HasReturn {
		t.Fatalf("expected no return type")
	}
}

func TestParseOptionalTypesInListAndMap(t *testing.T) {
	input := `model Text {}
rpc ListOptional() list[string?]
rpc MapOptional() map[Text?]
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schema.RPCs[0].Returns.Kind != TypeList || schema.RPCs[0].Returns.Elem == nil || !schema.RPCs[0].Returns.Elem.Optional {
		t.Fatalf("expected optional list element")
	}
	if schema.RPCs[1].Returns.Kind != TypeMap || schema.RPCs[1].Returns.Value == nil || !schema.RPCs[1].Returns.Value.Optional {
		t.Fatalf("expected optional map value")
	}
}

func TestParseEmptyModelsWithComments(t *testing.T) {
	input := `# leading
model Empty {
    # inside
}

model User {
    name: string # field
}
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.Models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(schema.Models))
	}
	if len(schema.Models[0].Fields) != 0 {
		t.Fatalf("expected empty model fields")
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

func TestFormatSchemaPreservesComments(t *testing.T) {
	input := `# top comment
model User { # model comment
    id: int # id comment
}

# between models and rpcs
rpc GetUser(
    id: int, # param comment
) User # return comment
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != input {
		t.Fatalf("formatted output mismatch:\n%s", formatted)
	}
}

func TestFormatSchemaNoReturnRPCWithComments(t *testing.T) {
	input := `# head
rpc Ping() # trailing
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != input {
		t.Fatalf("formatted output mismatch:\n%s", formatted)
	}
}

func TestFormatSchemaZeroArgReturnInlineComment(t *testing.T) {
	input := `model User {
}

rpc Get() User # comment
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != input {
		t.Fatalf("formatted output mismatch:\n%s", formatted)
	}
}

func TestFormatSchemaCommentPlacementAroundClosers(t *testing.T) {
	input := `model User {
    name: string
} # end model

rpc Hello(
    name: string,
) # end params
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != input {
		t.Fatalf("formatted output mismatch:\n%s", formatted)
	}
}

func TestFormatSchemaIdempotent(t *testing.T) {
	input := `# top
model User {
    id: int # id
}

rpc Ping()
rpc Echo(
    msg: string,
) string
`
	schema, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	schema2, err := Parse(first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := FormatSchema(schema2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first != second {
		t.Fatalf("expected idempotent formatting")
	}
}
