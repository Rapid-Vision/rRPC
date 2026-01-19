package formatter

import (
	"testing"

	"github.com/Rapid-Vision/rRPC/internal/parser"

	_ "embed"
)

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
	schema, err := parser.Parse(input)
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
	schema, err := parser.Parse(input)
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
	schema, err := parser.Parse(input)
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

func TestFormatSchemaInlineCommentAfterEmptyModel(t *testing.T) {
	input := `model User {
}

rpc FindUser() User # comment
`
	schema, err := parser.Parse(input)
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
	schema, err := parser.Parse(input)
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
	schema, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	schema2, err := parser.Parse(first)
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

//go:embed test_input.rrpc
var testInput string

//go:embed test_expected.rrpc
var testExpected string

func TestFormatFromTestFile(t *testing.T) {
	schema, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != testExpected {
		t.Fatalf("formatted output mismatch:\n%s", formatted)
	}
}
