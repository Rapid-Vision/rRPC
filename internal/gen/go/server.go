package gogen

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"

	_ "embed"
	"go/format"
)

//go:embed server_models.go.tmpl
var serverModelsTemplate string

//go:embed server_errors.go.tmpl
var serverErrorsTemplate string

//go:embed server_utils.go.tmpl
var serverUtilsTemplate string

//go:embed server_rpcs.go.tmpl
var serverRPCsTemplate string

type templateData struct {
	Package string
	Models  []parser.Model
	RPCs    []parser.RPC
}

func Generate(schema *parser.Schema, pkg string) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	return GenerateWithPrefix(schema, pkg, "rpc")
}

func GenerateWithPrefix(schema *parser.Schema, pkg, prefix string) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	data := templateData{
		Package: pkg,
		Models:  schema.Models,
		RPCs:    schema.RPCs,
	}
	funcMap := template.FuncMap{
		"modelTypeName":  modelTypeName,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"goType":         goType,
		"rpcParamsName":  rpcParamsName,
		"rpcResultName":  rpcResultName,
		"rpcHandlerName": rpcHandlerName,
		"rpcMethodName":  rpcMethodName,
		"rpcRoute": func(name string) string {
			return rpcRoute(prefix, name)
		},
		"resultField": resultField,
		"hasReturn":   hasReturn,
		"hasRPCs": func(data templateData) bool {
			return len(data.RPCs) > 0
		},
		"usesRawInModels": func(data templateData) bool {
			return parser.UsesRawInModels(*schema)
		},
		"usesRawInRPCs": func(data templateData) bool {
			return parser.UsesRawInRPCs(*schema)
		},
		"usesJSONDecoder": func(data templateData) bool {
			return usesJSONDecoder(data.RPCs)
		},
	}

	templates := map[string]string{
		"models.go": serverModelsTemplate,
		"errors.go": serverErrorsTemplate,
		"utils.go":  serverUtilsTemplate,
		"rpcs.go":   serverRPCsTemplate,
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

func usesJSONDecoder(rpcs []parser.RPC) bool {
	for _, rpc := range rpcs {
		if len(rpc.Parameters) > 0 {
			return true
		}
	}
	return false
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

func rpcRoute(prefix, name string) string {
	return "POST " + rpcPath(prefix, name)
}

func rpcPath(prefix, name string) string {
	p := strings.Trim(prefix, "/")
	route := utils.NewIdentifierName(name).SnakeCase()
	if p == "" {
		return "/" + route
	}
	return "/" + p + "/" + route
}

func resultField(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).PascalCase()
	}
	return "Result"
}

func hasReturn(rpc parser.RPC) bool {
	return rpc.HasReturn
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
	case "json":
		return "any"
	case "raw":
		return "json.RawMessage"
	default:
		return utils.NewIdentifierName(name).PascalCase() + "Model"
	}
}
