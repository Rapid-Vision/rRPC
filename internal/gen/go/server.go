package gogen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"

	_ "embed"
	"go/format"
)

//go:embed server.go.tmpl
var serverTemplate string

type templateData struct {
	Package string
	Models  []parser.Model
	RPCs    []parser.RPC
}

func Generate(schema *parser.Schema, pkg string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	tmpl, err := template.New("server.go.tmpl").Funcs(template.FuncMap{
		"modelTypeName":  modelTypeName,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"goType":         goType,
		"rpcParamsName":  rpcParamsName,
		"rpcResultName":  rpcResultName,
		"rpcHandlerName": rpcHandlerName,
		"rpcMethodName":  rpcMethodName,
		"rpcRoute":       rpcRoute,
		"resultField":    resultField,
	}).Parse(serverTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := templateData{
		Package: pkg,
		Models:  schema.Models,
		RPCs:    schema.RPCs,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("formatting error: %w", err)
	}
	return string(formatted), nil
}

func modelTypeName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Model"
}

func fieldName(name string) string {
	return utils.NewIdentifierName(name).PascalCase()
}

func jsonName(name string) string {
	return utils.NewIdentifierName(name).SnakeCase()
}

func rpcParamsName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Params"
}

func rpcResultName(name string) string {
	return utils.NewIdentifierName(name).PascalCase() + "Result"
}

func rpcHandlerName(name string) string {
	return "Create" + utils.NewIdentifierName(name).PascalCase() + "Handler"
}

func rpcMethodName(name string) string {
	return utils.NewIdentifierName(name).PascalCase()
}

func rpcRoute(name string) string {
	return "POST /" + utils.NewIdentifierName(name).SnakeCase()
}

func rpcPath(name string) string {
	return "/" + utils.NewIdentifierName(name).SnakeCase()
}

func resultField(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).PascalCase()
	}
	return "Result"
}

func goType(t parser.TypeRef) string {
	var base string
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			base = "[]any"
		} else {
			base = "[]" + goType(*t.Elem)
		}
	case parser.TypeMap:
		valueType := "any"
		if t.Value != nil {
			valueType = goType(*t.Value)
		}
		base = "map[string]" + valueType
	default:
		base = identType(t.Name)
	}

	if t.Optional {
		return "*" + base
	}
	return base
}

func identType(name string) string {
	switch name {
	case "string":
		return "string"
	case "int":
		return "int"
	case "bool":
		return "bool"
	default:
		return utils.NewIdentifierName(name).PascalCase() + "Model"
	}
}
