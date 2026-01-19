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

	anchorsByLine := buildAnchorsByLine(schema)
	leadingComments, trailingComments := partitionComments(schema, anchorsByLine)
	decls := resolveDecls(schema)

	var b strings.Builder
	comments := newCommentEmitter(&b, leadingComments, trailingComments)
	totalBlocks := len(decls)

	for i, decl := range decls {
		switch decl.Kind {
		case parser.DeclModel:
			if decl.Model == nil {
				continue
			}
			writeModel(&b, comments, *decl.Model)
		case parser.DeclRPC:
			if decl.RPC == nil {
				continue
			}
			writeRPC(&b, comments, *decl.RPC)
		}
		if i+1 < totalBlocks {
			b.WriteString("\n")
		}
	}

	comments.EmitLeading(math.MaxInt, "")

	return b.String(), nil
}

func buildAnchorsByLine(schema *parser.Schema) map[int][]anchorInfo {
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
	return anchorsByLine
}

func partitionComments(schema *parser.Schema, anchorsByLine map[int][]anchorInfo) ([]parser.Comment, map[anchorKey][]parser.Comment) {
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
	return leadingComments, trailingComments
}

func resolveDecls(schema *parser.Schema) []parser.Decl {
	if len(schema.Decls) > 0 {
		return schema.Decls
	}
	decls := make([]parser.Decl, 0, len(schema.Models)+len(schema.RPCs))
	for i := range schema.Models {
		decls = append(decls, parser.Decl{Kind: parser.DeclModel, Model: &schema.Models[i]})
	}
	for i := range schema.RPCs {
		decls = append(decls, parser.Decl{Kind: parser.DeclRPC, RPC: &schema.RPCs[i]})
	}
	return decls
}

type commentEmitter struct {
	builder    *strings.Builder
	leading    []parser.Comment
	trailing   map[anchorKey][]parser.Comment
	leadingIdx int
}

func newCommentEmitter(builder *strings.Builder, leading []parser.Comment, trailing map[anchorKey][]parser.Comment) *commentEmitter {
	return &commentEmitter{
		builder:  builder,
		leading:  leading,
		trailing: trailing,
	}
}

func (e *commentEmitter) EmitLeading(line int, indent string) {
	for e.leadingIdx < len(e.leading) && e.leading[e.leadingIdx].Line < line {
		e.builder.WriteString(indent)
		e.builder.WriteString(e.leading[e.leadingIdx].Text)
		e.builder.WriteString("\n")
		e.leadingIdx++
	}
}

func (e *commentEmitter) AppendTrailing(key anchorKey) {
	for _, comment := range e.trailing[key] {
		e.builder.WriteString(" ")
		e.builder.WriteString(comment.Text)
	}
}

func writeModel(b *strings.Builder, comments *commentEmitter, model parser.Model) {
	comments.EmitLeading(model.Line, "")
	b.WriteString("model ")
	b.WriteString(model.Name)
	b.WriteString(" {")
	comments.AppendTrailing(modelAnchorKey(model))
	b.WriteString("\n")
	for _, field := range model.Fields {
		comments.EmitLeading(field.Line, "    ")
		b.WriteString("    ")
		b.WriteString(field.Name)
		b.WriteString(": ")
		b.WriteString(parser.FormatType(field.Type))
		comments.AppendTrailing(fieldAnchorKey(field))
		b.WriteString("\n")
	}
	if model.EndLine > 0 {
		comments.EmitLeading(model.EndLine, "    ")
	}
	b.WriteString("}")
	comments.AppendTrailing(modelEndAnchorKey(model))
	b.WriteString("\n")
}

func writeRPC(b *strings.Builder, comments *commentEmitter, rpc parser.RPC) {
	comments.EmitLeading(rpc.Line, "")
	if len(rpc.Parameters) == 0 {
		b.WriteString("rpc ")
		b.WriteString(rpc.Name)
		b.WriteString("()")
		comments.AppendTrailing(rpcAnchorKey(rpc))
		if rpc.HasReturn {
			b.WriteString(" ")
			b.WriteString(parser.FormatType(rpc.Returns))
			comments.AppendTrailing(rpcReturnAnchorKey(rpc))
		}
		b.WriteString("\n")
		return
	}

	b.WriteString("rpc ")
	b.WriteString(rpc.Name)
	b.WriteString("(")
	comments.AppendTrailing(rpcAnchorKey(rpc))
	b.WriteString("\n")
	for _, param := range rpc.Parameters {
		comments.EmitLeading(param.Line, "    ")
		b.WriteString("    ")
		b.WriteString(param.Name)
		b.WriteString(": ")
		b.WriteString(parser.FormatType(param.Type))
		b.WriteString(",")
		comments.AppendTrailing(fieldAnchorKey(param))
		b.WriteString("\n")
	}
	if rpc.ParamsEndLine > 0 {
		comments.EmitLeading(rpc.ParamsEndLine, "    ")
	}
	if rpc.HasReturn && rpc.Returns.Line > 0 {
		comments.EmitLeading(rpc.Returns.Line, "")
	}
	b.WriteString(")")
	if rpc.HasReturn {
		b.WriteString(" ")
		b.WriteString(parser.FormatType(rpc.Returns))
		comments.AppendTrailing(rpcReturnAnchorKey(rpc))
	} else if rpc.ParamsEndLine > 0 {
		comments.AppendTrailing(rpcParamsEndAnchorKey(rpc))
	}
	b.WriteString("\n")
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
	return anchorKey{line: rpc.ParamsEndLine, col: rpc.ParamsEndCol, kind: "rpc_params_end"}
}

func rpcReturnAnchorKey(rpc parser.RPC) anchorKey {
	return anchorKey{line: rpc.Returns.Line, col: rpc.Returns.Col, kind: "rpc_return"}
}
