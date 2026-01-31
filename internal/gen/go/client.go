package gogen

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"

	_ "embed"
	"go/format"
)

//go:embed client_models.go.tmpl
var clientModelsTemplate string

//go:embed client_errors.go.tmpl
var clientErrorsTemplate string

//go:embed client_client.go.tmpl
var clientClientTemplate string

//go:embed client_transport.go.tmpl
var clientTransportTemplate string

//go:embed client_rpcs.go.tmpl
var clientRPCsTemplate string

func GenerateClient(schema *parser.Schema, pkg string) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	return GenerateClientWithPrefix(schema, pkg, "rpc")
}

func GenerateClientWithPrefix(schema *parser.Schema, pkg, prefix string) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	data := templateData{
		Package: pkg,
		Models:  schema.Models,
		RPCs:    schema.RPCs,
	}
	funcMap := template.FuncMap{
		"modelTypeName": modelTypeName,
		"fieldName":     fieldName,
		"jsonName":      jsonName,
		"goType":        goType,
		"rpcParamsName": rpcParamsName,
		"rpcResultName": rpcResultName,
		"rpcMethodName": rpcMethodName,
		"rpcPath": func(name string) string {
			return rpcPath(prefix, name)
		},
		"resultField": resultField,
		"hasReturn":   hasReturn,
		"usesRawInModels": func(data templateData) bool {
			return usesRawInModels(data.Models)
		},
		"usesRawInRPCs": func(data templateData) bool {
			return usesRawInRPCs(data.RPCs)
		},
		"hasRPCs": func(data templateData) bool {
			return len(data.RPCs) > 0
		},
	}

	templates := map[string]string{
		"models.go":    clientModelsTemplate,
		"errors.go":    clientErrorsTemplate,
		"client.go":    clientClientTemplate,
		"transport.go": clientTransportTemplate,
		"rpcs.go":      clientRPCsTemplate,
	}

	files := make(map[string]string, len(templates))
	for name, tmplText := range templates {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplText)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		var buf bytes.Buffer
		buf.WriteString("// THIS CODE IS GENERATED\n\n")
		buf.WriteString("package ")
		buf.WriteString(pkg)
		buf.WriteString("\n\n")
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("execute template %s: %w", name, err)
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("formatting error %s: %w", name, err)
		}
		files[name] = string(formatted)
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

func usesRawInModels(models []parser.Model) bool {
	for _, model := range models {
		for _, field := range model.Fields {
			if hasRawType(field.Type) {
				return true
			}
		}
	}
	return false
}

func usesRawInRPCs(rpcs []parser.RPC) bool {
	for _, rpc := range rpcs {
		for _, param := range rpc.Parameters {
			if hasRawType(param.Type) {
				return true
			}
		}
		if rpc.HasReturn && hasRawType(rpc.Returns) {
			return true
		}
	}
	return false
}

func hasRawType(t parser.TypeRef) bool {
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return false
		}
		return hasRawType(*t.Elem)
	case parser.TypeMap:
		if t.Value == nil {
			return false
		}
		return hasRawType(*t.Value)
	default:
		return t.Name == "raw"
	}
}
