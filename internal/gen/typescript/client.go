package tsgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed client.ts.tmpl
var clientTemplate string

type templateData struct {
	Models []parser.Model
	RPCs   []parser.RPC
	Prefix string
}

func GenerateClient(schema *parser.Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	return GenerateClientWithPrefix(schema, "rpc")
}

func GenerateClientWithPrefix(schema *parser.Schema, prefix string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	tmpl, err := template.New("client.ts.tmpl").Funcs(template.FuncMap{
		"className":      className,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"tsType":         tsType,
		"rpcMethodName":  rpcMethodName,
		"rpcPath":        rpcPath,
		"rpcParamsName":  rpcParamsName,
		"rpcResultName":  rpcResultName,
		"resultField":    resultField,
		"hasParameters":  hasParameters,
		"hasModelFields": hasModelFields,
		"hasReturn":      hasReturn,
	}).Parse(clientTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := templateData{
		Models: schema.Models,
		RPCs:   schema.RPCs,
		Prefix: prefixPath(prefix),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func GenerateTypeScriptIndex(schema *parser.Schema) string {
	var b strings.Builder
	b.WriteString("// THIS CODE IS GENERATED\n\n")
	b.WriteString("export {\n")
	b.WriteString("\tRPCClient,\n")
	b.WriteString("\tRPCErrorException,\n")
	b.WriteString("\tCustomRPCError,\n")
	b.WriteString("\tValidationRPCError,\n")
	b.WriteString("\tInputRPCError,\n")
	b.WriteString("\tUnauthorizedRPCError,\n")
	b.WriteString("\tForbiddenRPCError,\n")
	b.WriteString("\tNotImplementedRPCError,\n")
	b.WriteString("} from \"./client\";\n\n")

	b.WriteString("export type {\n")
	b.WriteString("\tRPCClientOptions,\n")
	b.WriteString("\tRPCErrorType,\n")
	b.WriteString("\tRPCError,\n")
	for _, model := range schema.Models {
		b.WriteString("\t")
		b.WriteString(className(model.Name))
		b.WriteString(",\n")
	}
	for _, rpc := range schema.RPCs {
		if len(rpc.Parameters) > 0 {
			b.WriteString("\t")
			b.WriteString(rpcParamsName(rpc.Name))
			b.WriteString(",\n")
		}
		if rpc.HasReturn {
			b.WriteString("\t")
			b.WriteString(rpcResultName(rpc.Name))
			b.WriteString(",\n")
		}
	}
	b.WriteString("} from \"./client\";\n")
	return b.String()
}

func className(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Model"
}

func rpcParamsName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Params"
}

func rpcResultName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Result"
}

func fieldName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func jsonName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func rpcMethodName(name string) string {
	pascal := utils.NewIdentifierName(name).PascalCase()
	if pascal == "" {
		return ""
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func rpcPath(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func tsType(t parser.TypeRef) string {
	base := tsBaseType(t)
	if t.Optional {
		return base + " | null"
	}
	return base
}

func tsBaseType(t parser.TypeRef) string {
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return "Array<any>"
		}
		return "Array<" + tsType(*t.Elem) + ">"
	case parser.TypeMap:
		valueType := "any"
		if t.Value != nil {
			valueType = tsType(*t.Value)
		}
		return "Record<string, " + valueType + ">"
	default:
		switch t.Name {
		case "string":
			return "string"
		case "int":
			return "number"
		case "bool":
			return "boolean"
		case "json":
			return "any"
		case "raw":
			return "any"
		default:
			return utils.NewIdentifierName(t.Name).PascalCase() + "Model"
		}
	}
}

func resultField(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).SnakeCase()
	}
	return "result"
}

func prefixPath(prefix string) string {
	p := strings.Trim(prefix, "/")
	if p == "" {
		return ""
	}
	return "/" + p
}

func hasParameters(rpc parser.RPC) bool {
	return len(rpc.Parameters) > 0
}

func hasReturn(rpc parser.RPC) bool {
	return rpc.HasReturn
}

func hasModelFields(model parser.Model) bool {
	return len(model.Fields) > 0
}
