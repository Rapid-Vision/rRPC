package pygen

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed errors.py.tmpl
var errorsTemplate string

//go:embed models.py.tmpl
var modelsTemplate string

//go:embed client.py.tmpl
var clientTemplate string

type templateData struct {
	Models []parser.Model
	RPCs   []parser.RPC
	Prefix string
}

func GenerateClient(schema *parser.Schema) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	return GenerateClientWithPrefix(schema, "rpc")
}

func GenerateClientWithPrefix(schema *parser.Schema, prefix string) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	data := templateData{
		Models: schema.Models,
		RPCs:   schema.RPCs,
		Prefix: prefixPath(prefix),
	}
	funcMap := template.FuncMap{
		"className":      className,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"pythonType":     pythonType,
		"rpcMethodName":  rpcMethodName,
		"resultField":    resultField,
		"decodeExpr":     decodeExpr,
		"hasParameters":  hasParameters,
		"hasModelFields": hasModelFields,
		"hasReturn":      hasReturn,
		"hasModels": func(data templateData) bool {
			return len(data.Models) > 0
		},
	}

	templates := map[string]string{
		"errors.py": errorsTemplate,
		"models.py": modelsTemplate,
		"client.py": clientTemplate,
	}

	files := make(map[string]string, len(templates))
	for name, tmplText := range templates {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplText)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		var buf bytes.Buffer
		buf.WriteString("# THIS CODE IS GENERATED\n\n")
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("execute template %s: %w", name, err)
		}
		files[name] = buf.String()
	}

	ordered := make([]string, 0, len(files))
	for name := range files {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	for _, name := range ordered {
		if strings.TrimSpace(files[name]) == "" {
			delete(files, name)
		}
	}
	return files, nil
}

func GeneratePythonInit(schema *parser.Schema) string {
	var b strings.Builder
	b.WriteString("# THIS CODE IS GENERATED\n\n")

	b.WriteString("from .client import RPCClient\n")
	b.WriteString("from .errors import RPCError\n")
	b.WriteString("from .errors import RPCErrorException\n")
	b.WriteString("from .errors import CustomRPCError\n")
	b.WriteString("from .errors import ValidationRPCError\n")
	b.WriteString("from .errors import InputRPCError\n")
	b.WriteString("from .errors import UnauthorizedRPCError\n")
	b.WriteString("from .errors import ForbiddenRPCError\n")
	b.WriteString("from .errors import NotImplementedRPCError\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("from .models import ")
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

func hasReturn(rpc parser.RPC) bool {
	return rpc.HasReturn
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
