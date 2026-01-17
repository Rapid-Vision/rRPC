package openapi

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed openapi.json.tmpl
var openApiTemplate string

type templateData struct {
	Title   string
	Version string
	Models  []parser.Model
	RPCs    []parser.RPC
	Prefix  string
}

func Generate(schema *parser.Schema, title, version string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	return GenerateWithPrefix(schema, title, version, "rpc")
}

func GenerateWithPrefix(schema *parser.Schema, title, version, prefix string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	if title == "" {
		title = "rRPC API"
	}
	if version == "" {
		version = "0.1.0"
	}
	tmpl, err := template.New("openapi.json.tmpl").Funcs(template.FuncMap{
		"modelSchemaName":  modelSchemaName,
		"paramsSchemaName": paramsSchemaName,
		"resultSchemaName": resultSchemaName,
		"errorSchemaName":  errorSchemaName,
		"rpcRoute": func(name string) string {
			return rpcRoute(prefix, name)
		},
		"rpcMethodName":    rpcMethodName,
		"jsonName":         jsonName,
		"schemaJSON":       schemaJSON,
		"requiredList":     requiredList,
		"toJSON":           toJSON,
		"hasParameters":    hasParameters,
		"resultField":      resultField,
		"add":              add,
	}).Parse(openApiTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := templateData{
		Title:   title,
		Version: version,
		Models:  schema.Models,
		RPCs:    schema.RPCs,
		Prefix:  prefix,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func modelSchemaName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Model"
}

func paramsSchemaName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Params"
}

func resultSchemaName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Result"
}

func errorSchemaName() string {
	return "RPCError"
}

func rpcRoute(prefix, name string) string {
	p := strings.Trim(prefix, "/")
	route := utils.NewIdentifierName(name).SnakeCase()
	if p == "" {
		return "/" + route
	}
	return "/" + p + "/" + route
}

func rpcMethodName(name string) string {
	return utils.NewIdentifierName(name).PascalCase()
}

func jsonName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func schemaJSON(t parser.TypeRef) string {
	schema := schemaForType(t)
	return toJSON(schema)
}

func schemaForType(t parser.TypeRef) map[string]any {
	var schema map[string]any
	switch t.Kind {
	case parser.TypeList:
		items := map[string]any{"type": "string"}
		if t.Elem != nil {
			items = schemaForType(*t.Elem)
		}
		schema = map[string]any{
			"type":  "array",
			"items": items,
		}
	case parser.TypeMap:
		additional := map[string]any{"type": "string"}
		if t.Value != nil {
			additional = schemaForType(*t.Value)
		}
		schema = map[string]any{
			"type":                 "object",
			"additionalProperties": additional,
		}
	default:
		switch t.Name {
		case "string":
			schema = map[string]any{"type": "string"}
		case "int":
			schema = map[string]any{"type": "integer", "format": "int32"}
	case "bool":
		schema = map[string]any{"type": "boolean"}
	case "json":
		schema = map[string]any{}
	case "raw":
		schema = map[string]any{}
	default:
		schema = map[string]any{
			"$ref": "#/components/schemas/" + modelSchemaName(t.Name),
		}
		}
	}

	if t.Optional {
		if _, ok := schema["$ref"]; ok {
			return map[string]any{
				"allOf":    []any{schema},
				"nullable": true,
			}
		}
		schema["nullable"] = true
	}
	return schema
}

func requiredList(fields []parser.Field) []string {
	required := make([]string, 0, len(fields))
	for _, field := range fields {
		if !field.Type.Optional {
			required = append(required, jsonName(field.Name))
		}
	}
	return required
}

func resultField(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).SnakeCase()
	}
	return "result"
}

func hasParameters(rpc parser.RPC) bool {
	return len(rpc.Parameters) > 0
}

func toJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func add(a, b int) int {
	return a + b
}
