package parser

import (
	"fmt"
	"math"
	"strings"

	"github.com/Rapid-Vision/rRPC/internal/utils"
)

func FormatSchema(schema *Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	if err := ValidateSchema(schema); err != nil {
		return "", err
	}

	lineIsNode := utils.NewSet[int]()
	for _, model := range schema.Models {
		if model.Line > 0 {
			lineIsNode.Add(model.Line)
		}
		if model.EndLine > 0 {
			lineIsNode.Add(model.EndLine)
		}
		for _, field := range model.Fields {
			if field.Line > 0 {
				lineIsNode.Add(field.Line)
			}
		}
	}
	for _, rpc := range schema.RPCs {
		if rpc.Line > 0 {
			lineIsNode.Add(rpc.Line)
		}
		if rpc.ParamsEndLine > 0 {
			lineIsNode.Add(rpc.ParamsEndLine)
		}
		if rpc.HasReturn && rpc.Returns.Line > 0 {
			lineIsNode.Add(rpc.Returns.Line)
		}
		for _, param := range rpc.Parameters {
			if param.Line > 0 {
				lineIsNode.Add(param.Line)
			}
		}
	}

	var leadingComments []Comment
	trailingComments := make(map[int][]Comment)
	for _, comment := range schema.Comments {
		if lineIsNode.Has(comment.Line) {
			trailingComments[comment.Line] = append(trailingComments[comment.Line], comment)
			continue
		}
		leadingComments = append(leadingComments, comment)
	}

	var b strings.Builder
	totalBlocks := len(schema.Models) + len(schema.RPCs)
	blockIndex := 0
	leadingIdx := 0

	emitLeading := func(line int, indent string) {
		for leadingIdx < len(leadingComments) && leadingComments[leadingIdx].Line < line {
			b.WriteString(indent)
			b.WriteString(leadingComments[leadingIdx].Text)
			b.WriteString("\n")
			leadingIdx++
		}
	}

	appendTrailing := func(line int) {
		for _, comment := range trailingComments[line] {
			b.WriteString(" ")
			b.WriteString(comment.Text)
		}
	}

	for _, model := range schema.Models {
		emitLeading(model.Line, "")
		blockIndex++
		b.WriteString("model ")
		b.WriteString(model.Name)
		b.WriteString(" {")
		appendTrailing(model.Line)
		b.WriteString("\n")
		for _, field := range model.Fields {
			emitLeading(field.Line, "    ")
			b.WriteString("    ")
			b.WriteString(field.Name)
			b.WriteString(": ")
			b.WriteString(formatType(field.Type))
			appendTrailing(field.Line)
			b.WriteString("\n")
		}
		if model.EndLine > 0 {
			emitLeading(model.EndLine, "    ")
		}
		b.WriteString("}")
		appendTrailing(model.EndLine)
		b.WriteString("\n")
		if blockIndex < totalBlocks {
			b.WriteString("\n")
		}
	}

	for _, rpc := range schema.RPCs {
		emitLeading(rpc.Line, "")
		blockIndex++
		if len(rpc.Parameters) == 0 {
			b.WriteString("rpc ")
			b.WriteString(rpc.Name)
			b.WriteString("()")
			appendTrailing(rpc.Line)
			if rpc.HasReturn {
				b.WriteString(" ")
				b.WriteString(formatType(rpc.Returns))
				if rpc.Returns.Line != rpc.Line {
					appendTrailing(rpc.Returns.Line)
				}
			}
			b.WriteString("\n")
		} else {
			b.WriteString("rpc ")
			b.WriteString(rpc.Name)
			b.WriteString("(")
			appendTrailing(rpc.Line)
			b.WriteString("\n")
			for _, param := range rpc.Parameters {
				emitLeading(param.Line, "    ")
				b.WriteString("    ")
				b.WriteString(param.Name)
				b.WriteString(": ")
				b.WriteString(formatType(param.Type))
				b.WriteString(",")
				appendTrailing(param.Line)
				b.WriteString("\n")
			}
			if rpc.ParamsEndLine > 0 {
				emitLeading(rpc.ParamsEndLine, "    ")
			}
			if rpc.HasReturn && rpc.Returns.Line > 0 {
				emitLeading(rpc.Returns.Line, "")
			}
			b.WriteString(")")
			if rpc.HasReturn {
				b.WriteString(" ")
				b.WriteString(formatType(rpc.Returns))
				if rpc.ParamsEndLine > 0 {
					appendTrailing(rpc.ParamsEndLine)
				}
				if rpc.Returns.Line != rpc.ParamsEndLine {
					appendTrailing(rpc.Returns.Line)
				}
			} else if rpc.ParamsEndLine > 0 {
				appendTrailing(rpc.ParamsEndLine)
			}
			b.WriteString("\n")
		}
		if blockIndex < totalBlocks {
			b.WriteString("\n")
		}
	}

	emitLeading(math.MaxInt, "")

	return b.String(), nil
}
