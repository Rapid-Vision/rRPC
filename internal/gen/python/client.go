package pygen

import (
	"fmt"
	"strings"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
)

type Generator struct {
	b strings.Builder
}

func GenerateClient(schema *parser.Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	gen := &Generator{}
	gen.writeHeader()
	gen.writeModels(schema.Models)
	gen.writeClient(schema.RPCs)
	return gen.b.String(), nil
}

func (g *Generator) writeHeader() {
	g.b.WriteString("from __future__ import annotations\n\n")
	g.b.WriteString("from dataclasses import dataclass\n")
	g.b.WriteString("from typing import Any, Dict, List, Optional\n")
	g.b.WriteString("import json\n")
	g.b.WriteString("import urllib.error\n")
	g.b.WriteString("import urllib.request\n\n")
}

func (g *Generator) writeModels(models []parser.Model) {
	for _, model := range models {
		name := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		g.b.WriteString("@dataclass\n")
		g.b.WriteString("class ")
		g.b.WriteString(name)
		g.b.WriteString(":\n")
		if len(model.Fields) == 0 {
			g.b.WriteString("    pass\n\n")
			continue
		}
		for _, field := range model.Fields {
			fieldName := utils.NewIdentifierName(field.Name).SnakeCase()
			g.b.WriteString("    ")
			g.b.WriteString(fieldName)
			g.b.WriteString(": ")
			g.b.WriteString(g.pythonType(field.Type))
			g.b.WriteString("\n")
		}
		g.b.WriteString("\n")
		g.b.WriteString("    @staticmethod\n")
		g.b.WriteString("    def from_dict(data: Dict[str, Any]) -> \"")
		g.b.WriteString(name)
		g.b.WriteString("\":\n")
		g.b.WriteString("        return ")
		g.b.WriteString(name)
		g.b.WriteString("(\n")
		for _, field := range model.Fields {
			fieldName := utils.NewIdentifierName(field.Name).SnakeCase()
			jsonKey := utils.NewIdentifierName(field.Name).SnakeCase()
			g.b.WriteString("            ")
			g.b.WriteString(fieldName)
			g.b.WriteString("=")
			g.b.WriteString(g.decodeExpr(field.Type, "data.get(\""+jsonKey+"\")"))
			g.b.WriteString(",\n")
		}
		g.b.WriteString("        )\n\n")
	}
}

func (g *Generator) writeClient(rpcs []parser.RPC) {
	g.b.WriteString("class RPCClient:\n")
	g.b.WriteString("    def __init__(self, base_url: str) -> None:\n")
	g.b.WriteString("        self.base_url = base_url.rstrip(\"/\")\n\n")

	g.b.WriteString("    def _request(self, path: str, payload: Optional[Dict[str, Any]]) -> Any:\n")
	g.b.WriteString("        url = f\"{self.base_url}/{path}\"\n")
	g.b.WriteString("        data = None\n")
	g.b.WriteString("        if payload is not None:\n")
	g.b.WriteString("            data = json.dumps(payload).encode(\"utf-8\")\n")
	g.b.WriteString("        req = urllib.request.Request(url, data=data, method='POST', headers={\"Content-Type\": \"application/json\"})\n")
	g.b.WriteString("        try:\n")
	g.b.WriteString("            with urllib.request.urlopen(req) as resp:\n")
	g.b.WriteString("                body = resp.read()\n")
	g.b.WriteString("        except urllib.error.HTTPError as err:\n")
	g.b.WriteString("            detail = err.read().decode(\"utf-8\")\n")
	g.b.WriteString("            raise RuntimeError(f\"rpc error: {detail}\") from err\n")
	g.b.WriteString("        if not body:\n")
	g.b.WriteString("            return None\n")
	g.b.WriteString("        return json.loads(body.decode(\"utf-8\"))\n\n")

	for _, rpc := range rpcs {
		methodName := utils.NewIdentifierName(rpc.Name).SnakeCase()
		returnType := g.pythonType(rpc.Returns)
		g.b.WriteString("    def ")
		g.b.WriteString(methodName)
		g.b.WriteString("(self")
		for _, param := range rpc.Parameters {
			paramName := utils.NewIdentifierName(param.Name).SnakeCase()
			g.b.WriteString(", ")
			g.b.WriteString(paramName)
			g.b.WriteString(": ")
			g.b.WriteString(g.pythonType(param.Type))
			if param.Type.Optional {
				g.b.WriteString(" = None")
			}
		}
		g.b.WriteString(") -> ")
		g.b.WriteString(returnType)
		g.b.WriteString(":\n")
		if len(rpc.Parameters) == 0 {
			g.b.WriteString("        payload = None\n")
		} else {
			g.b.WriteString("        payload = {\n")
			for _, param := range rpc.Parameters {
				paramName := utils.NewIdentifierName(param.Name).SnakeCase()
				jsonKey := utils.NewIdentifierName(param.Name).SnakeCase()
				g.b.WriteString("            \"")
				g.b.WriteString(jsonKey)
				g.b.WriteString("\": ")
				g.b.WriteString(paramName)
				g.b.WriteString(",\n")
			}
			g.b.WriteString("        }\n")
		}
		g.b.WriteString("        data = self._request(\"")
		g.b.WriteString(methodName)
		g.b.WriteString("\", payload)\n")
		resultKey := g.resultFieldName(rpc.Returns)
		g.b.WriteString("        value = data.get(\"")
		g.b.WriteString(resultKey)
		g.b.WriteString("\") if isinstance(data, dict) else data\n")
		g.b.WriteString("        return ")
		g.b.WriteString(g.decodeExpr(rpc.Returns, "value"))
		g.b.WriteString("\n\n")
	}
}

func (g *Generator) pythonType(t parser.TypeRef) string {
	base := g.pythonBaseType(t)
	if t.Optional {
		return "Optional[" + base + "]"
	}
	return base
}

func (g *Generator) pythonBaseType(t parser.TypeRef) string {
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return "List[Any]"
		}
		return "List[" + g.pythonType(*t.Elem) + "]"
	case parser.TypeMap:
		keyType := "str"
		valueType := "Any"
		if t.Key != nil {
			keyType = g.pythonType(*t.Key)
		}
		if t.Value != nil {
			valueType = g.pythonType(*t.Value)
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

func (g *Generator) decodeExpr(t parser.TypeRef, value string) string {
	if t.Optional {
		return fmt.Sprintf("None if %s is None else %s", value, g.decodeExpr(stripOptional(t), value))
	}
	switch t.Kind {
	case parser.TypeList:
		if t.Elem == nil {
			return value
		}
		itemExpr := g.decodeExpr(*t.Elem, "item")
		return fmt.Sprintf("[%s for item in %s]", itemExpr, value)
	case parser.TypeMap:
		if t.Key == nil || t.Value == nil {
			return value
		}
		keyExpr := g.decodeMapKeyExpr(*t.Key, "k")
		valExpr := g.decodeExpr(*t.Value, "v")
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

func (g *Generator) decodeMapKeyExpr(t parser.TypeRef, value string) string {
	switch t.Name {
	case "int":
		return "int(" + value + ")"
	case "bool":
		return value + ".lower() == \"true\""
	default:
		return value
	}
}

func (g *Generator) resultFieldName(t parser.TypeRef) string {
	if t.Kind == parser.TypeIdent {
		return utils.NewIdentifierName(t.Name).SnakeCase()
	}
	return "result"
}

func stripOptional(t parser.TypeRef) parser.TypeRef {
	t.Optional = false
	return t
}
