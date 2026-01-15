package gogen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"

	_ "embed"
	"go/format"
)

//go:embed client.go.tmpl
var clientTemplate string

func GenerateClient(schema *parser.Schema, pkg string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	tmpl, err := template.New("client.go.tmpl").Funcs(template.FuncMap{
		"modelTypeName": modelTypeName,
		"fieldName":     fieldName,
		"jsonName":      jsonName,
		"goType":        goType,
		"rpcParamsName": rpcParamsName,
		"rpcResultName": rpcResultName,
		"rpcMethodName": rpcMethodName,
		"rpcPath":       rpcPath,
		"resultField":   resultField,
	}).Parse(clientTemplate)
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
