package formatter

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Rapid-Vision/rRPC/internal/parser"
)

type anchorKey struct {
	line int
	col  int
	kind string
}

type anchorInfo struct {
	col int
	key anchorKey
}

func FormatSchema(schema *parser.Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}
	if err := parser.ValidateSchema(schema); err != nil {
		return "", err
	}

	anchorsByLine := make(map[int][]anchorInfo)
	addAnchor := func(key anchorKey) {
		if key.line <= 0 {
			return
		}
		anchorsByLine[key.line] = append(anchorsByLine[key.line], anchorInfo{col: key.col, key: key})
	}
	for _, model := range schema.Models {
		addAnchor(modelAnchorKey(model))
		addAnchor(modelEndAnchorKey(model))
		for _, field := range model.Fields {
			addAnchor(fieldAnchorKey(field))
		}
	}
	for _, rpc := range schema.RPCs {
		addAnchor(rpcAnchorKey(rpc))
		if len(rpc.Parameters) > 0 {
			addAnchor(rpcParamsEndAnchorKey(rpc))
		}
		if rpc.HasReturn {
			addAnchor(rpcReturnAnchorKey(rpc))
		}
		for _, param := range rpc.Parameters {
			addAnchor(fieldAnchorKey(param))
		}
	}
	for line, anchors := range anchorsByLine {
		sort.Slice(anchors, func(i, j int) bool {
			return anchors[i].col < anchors[j].col
		})
		anchorsByLine[line] = anchors
	}

	var leadingComments []parser.Comment
	trailingComments := make(map[anchorKey][]parser.Comment)
	for _, comment := range schema.Comments {
		anchors := anchorsByLine[comment.Line]
		var target *anchorInfo
		for i := range anchors {
			if anchors[i].col > comment.Col {
				break
			}
			target = &anchors[i]
		}
		if target != nil {
			trailingComments[target.key] = append(trailingComments[target.key], comment)
		} else {
			leadingComments = append(leadingComments, comment)
		}
	}

	decls := schema.Decls
	if len(decls) == 0 {
		decls = make([]parser.Decl, 0, len(schema.Models)+len(schema.RPCs))
		for i := range schema.Models {
			decls = append(decls, parser.Decl{Kind: parser.DeclModel, Model: &schema.Models[i]})
		}
		for i := range schema.RPCs {
			decls = append(decls, parser.Decl{Kind: parser.DeclRPC, RPC: &schema.RPCs[i]})
		}
	}

	var b strings.Builder
	totalBlocks := len(decls)
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

	appendTrailing := func(key anchorKey) {
		for _, comment := range trailingComments[key] {
			b.WriteString(" ")
			b.WriteString(comment.Text)
		}
	}

	for _, decl := range decls {
		blockIndex++
		switch decl.Kind {
		case parser.DeclModel:
			if decl.Model == nil {
				continue
			}
			model := *decl.Model
			emitLeading(model.Line, "")
			b.WriteString("model ")
			b.WriteString(model.Name)
			b.WriteString(" {")
			appendTrailing(modelAnchorKey(model))
			b.WriteString("\n")
			for _, field := range model.Fields {
				emitLeading(field.Line, "    ")
				b.WriteString("    ")
				b.WriteString(field.Name)
				b.WriteString(": ")
				b.WriteString(parser.FormatType(field.Type))
				appendTrailing(fieldAnchorKey(field))
				b.WriteString("\n")
			}
			if model.EndLine > 0 {
				emitLeading(model.EndLine, "    ")
			}
			b.WriteString("}")
			appendTrailing(modelEndAnchorKey(model))
			b.WriteString("\n")
		case parser.DeclRPC:
			if decl.RPC == nil {
				continue
			}
			rpc := *decl.RPC
			emitLeading(rpc.Line, "")
			if len(rpc.Parameters) == 0 {
				b.WriteString("rpc ")
				b.WriteString(rpc.Name)
				b.WriteString("()")
				appendTrailing(rpcAnchorKey(rpc))
				if rpc.HasReturn {
					b.WriteString(" ")
					b.WriteString(parser.FormatType(rpc.Returns))
					appendTrailing(rpcReturnAnchorKey(rpc))
				}
				b.WriteString("\n")
			} else {
				b.WriteString("rpc ")
				b.WriteString(rpc.Name)
				b.WriteString("(")
				appendTrailing(rpcAnchorKey(rpc))
				b.WriteString("\n")
				for _, param := range rpc.Parameters {
					emitLeading(param.Line, "    ")
					b.WriteString("    ")
					b.WriteString(param.Name)
					b.WriteString(": ")
					b.WriteString(parser.FormatType(param.Type))
					b.WriteString(",")
					appendTrailing(fieldAnchorKey(param))
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
					b.WriteString(parser.FormatType(rpc.Returns))
					appendTrailing(rpcReturnAnchorKey(rpc))
				} else if rpc.ParamsEndLine > 0 {
					appendTrailing(rpcParamsEndAnchorKey(rpc))
				}
				b.WriteString("\n")
			}
		}
		if blockIndex < totalBlocks {
			b.WriteString("\n")
		}
	}

	emitLeading(math.MaxInt, "")

	return b.String(), nil
}

func modelAnchorKey(model parser.Model) anchorKey {
	return anchorKey{line: model.Line, col: model.Col, kind: "model"}
}

func modelEndAnchorKey(model parser.Model) anchorKey {
	if model.EndLine == model.Line && model.Line > 0 {
		return anchorKey{line: model.EndLine, col: modelSingleLineEndCol(model), kind: "model_end"}
	}
	return anchorKey{line: model.EndLine, col: 1, kind: "model_end"}
}

func modelSingleLineEndCol(model parser.Model) int {
	nameLen := len([]rune(model.Name))
	return model.Col + len([]rune("model ")) + nameLen + 2
}

func fieldAnchorKey(field parser.Field) anchorKey {
	return anchorKey{line: field.Line, col: field.Col, kind: "field"}
}

func rpcAnchorKey(rpc parser.RPC) anchorKey {
	return anchorKey{line: rpc.Line, col: rpc.Col, kind: "rpc"}
}

func rpcParamsEndAnchorKey(rpc parser.RPC) anchorKey {
	return anchorKey{line: rpc.ParamsEndLine, col: 1, kind: "rpc_params_end"}
}

func rpcReturnAnchorKey(rpc parser.RPC) anchorKey {
	return anchorKey{line: rpc.Returns.Line, col: rpc.Returns.Col, kind: "rpc_return"}
}
