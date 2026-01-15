package pygen

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed client.py.tmpl
var clientTemplate string

type templateData struct {
	Models []parser.Model
	RPCs   []parser.RPC
}

func GenerateClient(schema *parser.Schema) (string, error) {
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
		"decodeKeyExpr":  decodeKeyExpr,
		"hasParameters":  hasParameters,
		"hasModelFields": hasModelFields,
	}).Parse(clientTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := templateData{
		Models: schema.Models,
		RPCs:   schema.RPCs,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
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
		keyType := "str"
		valueType := "Any"
		if t.Key != nil {
			keyType = pythonType(*t.Key)
		}
		if t.Value != nil {
			valueType = pythonType(*t.Value)
		}
		return "Dict[" + keyType + ", " + valueType + "]"
	default:
		switch t.Name {
		case "string":
			return "str"
		case "int":
			return "int"
		case "bool":
			return "bool"
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
		if t.Key == nil || t.Value == nil {
			return value
		}
		keyExpr := decodeKeyExpr(*t.Key, "k")
		valExpr := decodeExpr(*t.Value, "v")
		return fmt.Sprintf("{%s: %s for k, v in %s.items()}", keyExpr, valExpr, value)
	default:
		switch t.Name {
		case "string", "int", "bool":
			return value
		default:
			return utils.NewIdentifierName(t.Name).PascalCase() + "Model.from_dict(" + value + ")"
		}
	}
}

func decodeKeyExpr(t parser.TypeRef, value string) string {
	switch t.Name {
	case "int":
		return "int(" + value + ")"
	default:
		return value
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
