package formatter

import (
	"testing"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/pmezard/go-difflib/difflib"

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
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(formatted),
			B:        difflib.SplitLines(testExpected),
			FromFile: "formatted",
			ToFile:   "expected",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)

		t.Fatalf("formatted output mismatch:\n%s", text)
	}
}

func TestFormatIdempotanceFromTestFile(t *testing.T) {
	schema, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	formatted, err := FormatSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	formatted_schema, err := parser.Parse(formatted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	formatted_twice, err := FormatSchema(formatted_schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted != formatted_twice {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(formatted),
			B:        difflib.SplitLines(formatted_twice),
			FromFile: "formatted",
			ToFile:   "formatted_twice",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)

		t.Fatalf("formatted idempotency failed:\n%s", text)
	}
}
