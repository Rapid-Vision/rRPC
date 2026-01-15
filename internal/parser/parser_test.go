package parser

import "testing"

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

func TestParseMapKeyValidation(t *testing.T) {
	input := `rpc GetMap() map[bool, string]`
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseListAndMapTypes(t *testing.T) {
	input := `rpc GetUsers() list[User]
rpc GetUsersByName() map[string, list[User]]
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
