package gogen

import (
	"fmt"
	"strings"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

type Generator struct {
	pkg string
	b   strings.Builder
}

func Generate(schema *parser.Schema, pkg string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	gen := &Generator{pkg: pkg}
	gen.writeHeader()
	gen.writeModels(schema.Models)
	gen.writeRPCTypes(schema.RPCs)
	gen.writeHandlerInterface(schema.RPCs)
	gen.writeHTTPHandlers(schema.RPCs)
	return gen.b.String(), nil
}

func (g *Generator) writeHeader() {
	g.b.WriteString("package ")
	g.b.WriteString(g.pkg)
	g.b.WriteString("\n\n")
	g.b.WriteString("import (\n")
	g.b.WriteString("\t\"encoding/json\"\n")
	g.b.WriteString("\t\"io\"\n")
	g.b.WriteString("\t\"net/http\"\n")
	g.b.WriteString(")\n\n")
}

func (g *Generator) writeModels(models []parser.Model) {
	for i, model := range models {
		if i > 0 {
			g.b.WriteString("\n")
		}
		typeName := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		g.b.WriteString("type ")
		g.b.WriteString(typeName)
		g.b.WriteString(" struct {\n")
		for _, field := range model.Fields {
			fieldName := utils.NewIdentifierName(field.Name).PascalCase()
			g.b.WriteString("\t")
			g.b.WriteString(fieldName)
			g.b.WriteString(" ")
			g.b.WriteString(g.goType(field.Type))
			g.b.WriteString(" `json:\"")
			g.b.WriteString(utils.NewIdentifierName(field.Name).SnakeCase())
			g.b.WriteString("\"`\n")
		}
		g.b.WriteString("}\n")
	}
}

func (g *Generator) writeRPCTypes(rpcs []parser.RPC) {
	for _, rpc := range rpcs {
		g.b.WriteString("\n")
		rpcName := utils.NewIdentifierName(rpc.Name).PascalCase()
		paramsName := rpcName + "Params"
		resultName := rpcName + "Result"

		g.b.WriteString("type ")
		g.b.WriteString(paramsName)
		g.b.WriteString(" struct {\n")
		for _, param := range rpc.Parameters {
			fieldName := utils.NewIdentifierName(param.Name).PascalCase()
			g.b.WriteString("\t")
			g.b.WriteString(fieldName)
			g.b.WriteString(" ")
			g.b.WriteString(g.goType(param.Type))
			g.b.WriteString(" `json:\"")
			g.b.WriteString(utils.NewIdentifierName(param.Name).SnakeCase())
			g.b.WriteString("\"`\n")
		}
		g.b.WriteString("}\n\n")

		g.b.WriteString("type ")
		g.b.WriteString(resultName)
		g.b.WriteString(" struct {\n")
		resFieldName := g.resultFieldName(rpc.Returns)
		g.b.WriteString("\t")
		g.b.WriteString(resFieldName)
		g.b.WriteString(" ")
		g.b.WriteString(g.goType(rpc.Returns))
		g.b.WriteString(" `json:\"")
		g.b.WriteString(utils.NewIdentifierName(resFieldName).SnakeCase())
		g.b.WriteString("\"`\n")
		g.b.WriteString("}\n")
	}
}

func (g *Generator) writeHandlerInterface(rpcs []parser.RPC) {
	g.b.WriteString("\n")
	g.b.WriteString("type RPCHandler interface {\n")
	for _, rpc := range rpcs {
		rpcName := utils.NewIdentifierName(rpc.Name).PascalCase()
		paramsName := rpcName + "Params"
		resultName := rpcName + "Result"
		g.b.WriteString("\t")
		g.b.WriteString(rpcName)
		g.b.WriteString("(")
		g.b.WriteString(paramsName)
		g.b.WriteString(") (")
		g.b.WriteString(resultName)
		g.b.WriteString(", error)\n")
	}
	g.b.WriteString("}\n")
}

func (g *Generator) writeHTTPHandlers(rpcs []parser.RPC) {
	g.b.WriteString("\n")
	g.b.WriteString("func CreateHTTPHandler(rpc RPCHandler) http.Handler {\n")
	g.b.WriteString("\tmux := http.NewServeMux()\n")
	for _, rpc := range rpcs {
		route := "POST /" + utils.NewIdentifierName(rpc.Name).SnakeCase()
		handler := "Create" + utils.NewIdentifierName(rpc.Name).PascalCase() + "Handler"
		g.b.WriteString("\tmux.Handle(\"")
		g.b.WriteString(route)
		g.b.WriteString("\", ")
		g.b.WriteString(handler)
		g.b.WriteString("(rpc))\n")
	}
	g.b.WriteString("\treturn mux\n")
	g.b.WriteString("}\n")

	for _, rpc := range rpcs {
		g.writeRPCHandler(rpc)
	}

	g.b.WriteString("\n")
	g.b.WriteString("type rpcError struct {\n")
	g.b.WriteString("\tMessage string `json:\"message\"`\n")
	g.b.WriteString("}\n")

	g.b.WriteString("\n")
	g.b.WriteString("func writeJSON(w http.ResponseWriter, status int, payload any) {\n")
	g.b.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	g.b.WriteString("\tw.WriteHeader(status)\n")
	g.b.WriteString("\t_ = json.NewEncoder(w).Encode(payload)\n")
	g.b.WriteString("}\n")

	g.b.WriteString("\n")
	g.b.WriteString("func writeError(w http.ResponseWriter, status int, err error) {\n")
	g.b.WriteString("\tmsg := \"error\"\n")
	g.b.WriteString("\tif err != nil {\n")
	g.b.WriteString("\t\tmsg = err.Error()\n")
	g.b.WriteString("\t}\n")
	g.b.WriteString("\twriteJSON(w, status, rpcError{Message: msg})\n")
	g.b.WriteString("}\n")
}

func (g *Generator) writeRPCHandler(rpc parser.RPC) {
	rpcName := utils.NewIdentifierName(rpc.Name).PascalCase()
	paramsName := rpcName + "Params"
	handlerName := "Create" + rpcName + "Handler"

	g.b.WriteString("\n")
	g.b.WriteString("func ")
	g.b.WriteString(handlerName)
	g.b.WriteString("(rpc RPCHandler) http.Handler {\n")
	g.b.WriteString("\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n")
	g.b.WriteString("\t\tvar params ")
	g.b.WriteString(paramsName)
	g.b.WriteString("\n")

	if len(rpc.Parameters) > 0 {
		g.b.WriteString("\t\tdecoder := json.NewDecoder(r.Body)\n")
		g.b.WriteString("\t\tdecoder.DisallowUnknownFields()\n")
		g.b.WriteString("\t\tif err := decoder.Decode(&params); err != nil && err != io.EOF {\n")
		g.b.WriteString("\t\t\twriteError(w, http.StatusBadRequest, err)\n")
		g.b.WriteString("\t\t\treturn\n")
		g.b.WriteString("\t\t}\n")
	}

	g.b.WriteString("\t\tres, err := rpc.")
	g.b.WriteString(rpcName)
	g.b.WriteString("(params)\n")
	g.b.WriteString("\t\tif err != nil {\n")
	g.b.WriteString("\t\t\twriteError(w, http.StatusInternalServerError, err)\n")
	g.b.WriteString("\t\t\treturn\n")
	g.b.WriteString("\t\t}\n")
	g.b.WriteString("\t\twriteJSON(w, http.StatusOK, res)\n")
	g.b.WriteString("\t})\n")
	g.b.WriteString("}\n")

}

func (g *Generator) identType(name string) string {
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

func (g *Generator) resultFieldName(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).PascalCase()
	}
	return "Result"
}

func (g *Generator) goType(t parser.TypeRef) string {
	var base string
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			base = "[]any"
		} else {
			base = "[]" + g.goType(*t.Elem)
		}
	case parser.TypeMap:
		keyType := "string"
		valueType := "any"
		if t.Key != nil {
			keyType = g.goType(*t.Key)
		}
		if t.Value != nil {
			valueType = g.goType(*t.Value)
		}
		base = "map[" + keyType + "]" + valueType
	default:
		base = g.identType(t.Name)
	}

	if t.Optional {
		return "*" + base
	}
	return base
}
