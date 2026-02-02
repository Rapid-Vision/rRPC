package tsgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

//go:embed errors.ts.tmpl
var errorsTemplate string

//go:embed models.ts.tmpl
var modelsTemplate string

//go:embed client.ts.tmpl
var clientTemplate string

type templateData struct {
	Models []parser.Model
	RPCs   []parser.RPC
	Prefix string
	Zod    bool
}

func GenerateClient(schema *parser.Schema) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	return GenerateClientWithPrefixAndZod(schema, "rpc", false)
}

func GenerateClientWithPrefix(schema *parser.Schema, prefix string) (map[string]string, error) {
	return GenerateClientWithPrefixAndZod(schema, prefix, false)
}

func GenerateClientWithPrefixAndZod(schema *parser.Schema, prefix string, zod bool) (map[string]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	data := templateData{
		Models: schema.Models,
		RPCs:   schema.RPCs,
		Prefix: prefixPath(prefix),
		Zod:    zod,
	}

	funcMap := template.FuncMap{
		"className":      className,
		"fieldName":      fieldName,
		"jsonName":       jsonName,
		"tsType":         tsType,
		"zodType":        zodType,
		"rpcMethodName":  rpcMethodName,
		"rpcPath":        rpcPath,
		"rpcParamsName":  rpcParamsName,
		"rpcResultName":  rpcResultName,
		"resultField":    resultField,
		"hasParameters":  hasParameters,
		"hasModelFields": hasModelFields,
		"hasReturn":      hasReturn,
		"hasTypes": func(data templateData) bool {
			if len(data.Models) > 0 {
				return true
			}
			for _, rpc := range data.RPCs {
				if len(rpc.Parameters) > 0 || rpc.HasReturn {
					return true
				}
			}
			return false
		},
	}

	templates := map[string]string{
		"errors.ts": errorsTemplate,
		"models.ts": modelsTemplate,
		"client.ts": clientTemplate,
	}

	files := make(map[string]string, len(templates))
	for name, tmplText := range templates {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplText)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		var buf bytes.Buffer
		buf.WriteString("// THIS CODE IS GENERATED\n\n")
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

func GenerateTypeScriptIndex(schema *parser.Schema) string {
	return GenerateTypeScriptIndexWithZod(schema, false)
}

func GenerateTypeScriptIndexWithZod(schema *parser.Schema, zod bool) string {
	var b strings.Builder
	b.WriteString("// THIS CODE IS GENERATED\n\n")

	b.WriteString("export { RPCClient } from \"./client\";\n")
	b.WriteString("export {\n")
	b.WriteString("\tRPCErrorException,\n")
	b.WriteString("\tCustomRPCError,\n")
	b.WriteString("\tValidationRPCError,\n")
	b.WriteString("\tInputRPCError,\n")
	b.WriteString("\tUnauthorizedRPCError,\n")
	b.WriteString("\tForbiddenRPCError,\n")
	b.WriteString("\tNotImplementedRPCError,\n")
	b.WriteString("} from \"./errors\";\n")
	hasModelsExports := len(schema.Models) > 0
	hasZodExports := false
	hasTypesExports := false
	for _, rpc := range schema.RPCs {
		if len(rpc.Parameters) > 0 {
			hasZodExports = true
			hasTypesExports = true
		}
		if rpc.HasReturn {
			hasTypesExports = true
		}
	}
	if hasModelsExports {
		hasTypesExports = true
		hasZodExports = true
	}
	if zod {
		if hasZodExports {
			b.WriteString("export {\n")
			for _, model := range schema.Models {
				b.WriteString("\t")
				b.WriteString(className(model.Name))
				b.WriteString("Schema,\n")
			}
			for _, rpc := range schema.RPCs {
				if len(rpc.Parameters) > 0 {
					b.WriteString("\t")
					b.WriteString(rpcParamsName(rpc.Name))
					b.WriteString("Schema,\n")
				}
			}
			b.WriteString("} from \"./models\";\n")
		}
	}
	b.WriteString("\n")

	if hasTypesExports {
		b.WriteString("export type {\n")
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
		b.WriteString("} from \"./models\";\n")
	}
	b.WriteString("export type { FetchFn, FetchInit, FetchResponse, RPCClientOptions } from \"./client\";\n")
	b.WriteString("export type { RPCErrorType, RPCError } from \"./errors\";\n")
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

func zodType(t parser.TypeRef) string {
	base := zodBaseType(t)
	if t.Optional {
		return "z.union([" + base + ", z.null()])"
	}
	return base
}

func zodBaseType(t parser.TypeRef) string {
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return "z.array(z.any())"
		}
		return "z.array(" + zodType(*t.Elem) + ")"
	case parser.TypeMap:
		valueType := "z.any()"
		if t.Value != nil {
			valueType = zodType(*t.Value)
		}
		return "z.record(z.string(), " + valueType + ")"
	default:
		switch t.Name {
		case "string":
			return "z.string()"
		case "int":
			return "z.number().int()"
		case "bool":
			return "z.boolean()"
		case "json":
			return "z.any()"
		case "raw":
			return "z.any()"
		default:
			return "z.lazy(() => " + utils.NewIdentifierName(t.Name).PascalCase() + "ModelSchema)"
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
