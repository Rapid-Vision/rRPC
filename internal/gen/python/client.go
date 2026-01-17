package pygen

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed client.py.tmpl
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
	tmpl, err := template.New("client.py.tmpl").Funcs(template.FuncMap{
		"className":      className,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"pythonType":     pythonType,
		"rpcMethodName":  rpcMethodName,
		"resultField":    resultField,
		"decodeExpr":     decodeExpr,
		"hasParameters":  hasParameters,
		"hasModelFields": hasModelFields,
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

func GeneratePythonInit(schema *parser.Schema) string {
	var b strings.Builder
	b.WriteString("from .client import RPCClient\n")
	b.WriteString("from .client import RPCError\n")
	b.WriteString("from .client import RPCErrorException\n")
	b.WriteString("from .client import CustomRPCError\n")
	b.WriteString("from .client import ValidationRPCError\n")
	b.WriteString("from .client import InputRPCError\n")
	b.WriteString("from .client import UnauthorizedRPCError\n")
	b.WriteString("from .client import ForbiddenRPCError\n")
	b.WriteString("from .client import NotImplementedRPCError\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("from .client import ")
		b.WriteString(className)
		b.WriteString("\n")
	}
	b.WriteString("\n__all__ = [\n")
	b.WriteString("    \"RPCClient\",\n")
	b.WriteString("    \"RPCError\",\n")
	b.WriteString("    \"RPCErrorException\",\n")
	b.WriteString("    \"CustomRPCError\",\n")
	b.WriteString("    \"ValidationRPCError\",\n")
	b.WriteString("    \"InputRPCError\",\n")
	b.WriteString("    \"UnauthorizedRPCError\",\n")
	b.WriteString("    \"ForbiddenRPCError\",\n")
	b.WriteString("    \"NotImplementedRPCError\",\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("    \"")
		b.WriteString(className)
		b.WriteString("\",\n")
	}
	b.WriteString("]\n")
	return b.String()
}

func className(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Model"
}

func fieldName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func jsonName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func rpcMethodName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func pythonType(t parser.TypeRef) string {
	base := pythonBaseType(t)
	if t.Optional {
		return "Optional[" + base + "]"
	}
	return base
}

func pythonBaseType(t parser.TypeRef) string {
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return "List[Any]"
		}
		return "List[" + pythonType(*t.Elem) + "]"
	case parser.TypeMap:
		valueType := "Any"
		if t.Value != nil {
			valueType = pythonType(*t.Value)
		}
		return "Dict[str, " + valueType + "]"
	default:
		switch t.Name {
		case "string":
			return "str"
		case "int":
			return "int"
		case "bool":
			return "bool"
		case "json":
			return "Any"
		case "raw":
			return "Any"
		default:
			return utils.NewIdentifierName(t.Name).PascalCase() + "Model"
		}
	}
}

func decodeExpr(t parser.TypeRef, value string) string {
	if t.Optional {
		return fmt.Sprintf("None if %s is None else %s", value, decodeExpr(stripOptional(t), value))
	}
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return value
		}
		itemExpr := decodeExpr(*t.Elem, "item")
		return fmt.Sprintf("[%s for item in %s]", itemExpr, value)
	case parser.TypeMap:
		if t.Value == nil {
			return value
		}
		valExpr := decodeExpr(*t.Value, "v")
		return fmt.Sprintf("{k: %s for k, v in %s.items()}", valExpr, value)
	default:
		switch t.Name {
		case "string", "int", "bool", "json", "raw":
			return value
		default:
			return utils.NewIdentifierName(t.Name).PascalCase() + "Model.from_dict(" + value + ")"
		}
	}
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

func hasModelFields(model parser.Model) bool {
	return len(model.Fields) > 0
}

func stripOptional(t parser.TypeRef) parser.TypeRef {
	t.Optional = false
	return t
}

func prefixPath(prefix string) string {
	p := strings.Trim(prefix, "/")
	if p == "" {
		return ""
	}
	return "/" + p
}
