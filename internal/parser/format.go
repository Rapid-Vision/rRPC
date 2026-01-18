package parser

import (
	"fmt"
	"strings"
)

func FormatSchema(schema *Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	if err := ValidateSchema(schema); err != nil {
		return "", err
	}

	var b strings.Builder
	totalBlocks := len(schema.Models) + len(schema.RPCs)
	blockIndex := 0

	for _, model := range schema.Models {
		blockIndex++
		b.WriteString("model ")
		b.WriteString(model.Name)
		b.WriteString(" {\n")
		for _, field := range model.Fields {
			b.WriteString("    ")
			b.WriteString(field.Name)
			b.WriteString(": ")
			b.WriteString(formatType(field.Type))
			b.WriteString("\n")
		}
		b.WriteString("}\n")
		if blockIndex < totalBlocks {
			b.WriteString("\n")
		}
	}

	for _, rpc := range schema.RPCs {
		blockIndex++
		if len(rpc.Parameters) == 0 {
			b.WriteString("rpc ")
			b.WriteString(rpc.Name)
			b.WriteString("() ")
			b.WriteString(formatType(rpc.Returns))
			b.WriteString("\n")
		} else {
			b.WriteString("rpc ")
			b.WriteString(rpc.Name)
			b.WriteString("(\n")
			for _, param := range rpc.Parameters {
				b.WriteString("    ")
				b.WriteString(param.Name)
				b.WriteString(": ")
				b.WriteString(formatType(param.Type))
				b.WriteString(",\n")
			}
			b.WriteString(") ")
			b.WriteString(formatType(rpc.Returns))
			b.WriteString("\n")
		}
		if blockIndex < totalBlocks {
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}
