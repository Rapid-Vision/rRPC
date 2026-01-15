package parser

import (
	"fmt"
	"strings"

	"github.com/Rapid-Vision/rRPC/internal/lexer"
)

type Schema struct {
	Models []Model
	RPCs   []RPC
}

func (s *Schema) Dump() string {
	var b strings.Builder
	totalChildren := len(s.Models) + len(s.RPCs)
	modelsLeft := len(s.Models)
	for _, model := range s.Models {
		modelsLeft--
		writeTreeLine(&b, 0, "Model: "+model.Name)
		fieldsLeft := len(model.Fields)
		for _, field := range model.Fields {
			fieldsLeft--
			writeTreeLine(&b, 1, "Field: "+field.Name)
			writeTreeLine(&b, 2, "Type: "+formatType(field.Type))
		}
		if len(model.Fields) == 0 {
			writeTreeLine(&b, 1, "Field: (none)")
		}
		totalChildren--
		if totalChildren == 0 {
			break
		}
	}

	rpcsLeft := len(s.RPCs)
	for _, rpc := range s.RPCs {
		rpcsLeft--
		writeTreeLine(&b, 0, "RPC: "+rpc.Name)
		writeTreeLine(&b, 1, "Params")
		if len(rpc.Parameters) == 0 {
			writeTreeLine(&b, 2, "Field: (none)")
		} else {
			paramsLeft := len(rpc.Parameters)
			for _, param := range rpc.Parameters {
				paramsLeft--
				writeTreeLine(&b, 2, "Field: "+param.Name)
				writeTreeLine(&b, 3, "Type: "+formatType(param.Type))
			}
		}
		writeTreeLine(&b, 1, "Returns")
		writeTreeLine(&b, 2, "Type: "+formatType(rpc.Returns))
	}
	return b.String()
}

type Model struct {
	Name   string
	Fields []Field
}

type RPC struct {
	Name       string
	Parameters []Field
	Returns    TypeRef
}

type Field struct {
	Name string
	Type TypeRef
}

type TypeRef struct {
	Kind     TypeKind
	Name     string
	Elem     *TypeRef
	Value    *TypeRef
	Optional bool
}

type TypeKind int

const (
	TypeIdent TypeKind = iota
	TypeList
	TypeMap
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func Parse(text string) (*Schema, error) {
	tokens, err := lexer.NewLexer(text).Tokenize()
	if err != nil {
		return nil, err
	}
	p := NewParser(tokens)
	schema, err := p.parseSchema()
	if err != nil {
		return nil, err
	}
	if err := ValidateSchema(schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) parseSchema() (*Schema, error) {
	var schema Schema
	for !p.atEnd() {
		switch p.peek().Type {
		case lexer.TokenModel:
			model, err := p.parseModel()
			if err != nil {
				return nil, err
			}
			schema.Models = append(schema.Models, model)
		case lexer.TokenRpc:
			rpc, err := p.parseRPC()
			if err != nil {
				return nil, err
			}
			schema.RPCs = append(schema.RPCs, rpc)
		default:
			return nil, p.unexpected("model or rpc")
		}
	}
	return &schema, nil
}

func (p *Parser) parseModel() (Model, error) {
	if _, err := p.expect(lexer.TokenModel); err != nil {
		return Model{}, err
	}
	name, err := p.expect(lexer.TokenIdentifier)
	if err != nil {
		return Model{}, err
	}
	if _, err := p.expect(lexer.TokenLBrace); err != nil {
		return Model{}, err
	}

	var fields []Field
	for !p.atEnd() && p.peek().Type != lexer.TokenRBrace {
		if p.peek().Type != lexer.TokenIdentifier {
			return Model{}, p.unexpected("field name or }")
		}
		field, err := p.parseField()
		if err != nil {
			return Model{}, err
		}
		fields = append(fields, field)
	}
	if _, err := p.expect(lexer.TokenRBrace); err != nil {
		return Model{}, err
	}
	return Model{Name: name.Value, Fields: fields}, nil
}

func (p *Parser) parseRPC() (RPC, error) {
	if _, err := p.expect(lexer.TokenRpc); err != nil {
		return RPC{}, err
	}
	name, err := p.expect(lexer.TokenIdentifier)
	if err != nil {
		return RPC{}, err
	}
	if _, err := p.expect(lexer.TokenLParen); err != nil {
		return RPC{}, err
	}

	var params []Field
	for !p.atEnd() && p.peek().Type != lexer.TokenRParen {
		if p.peek().Type != lexer.TokenIdentifier {
			return RPC{}, p.unexpected("parameter name or )")
		}
		field, err := p.parseField()
		if err != nil {
			return RPC{}, err
		}
		params = append(params, field)
		if p.match(lexer.TokenComma) {
			if p.atEnd() {
				return RPC{}, p.unexpected("parameter name or )")
			}
			if p.peek().Type == lexer.TokenRParen {
				break
			}
			continue
		}
		if p.atEnd() || p.peek().Type != lexer.TokenRParen {
			return RPC{}, p.unexpected("comma or )")
		}
	}
	if _, err := p.expect(lexer.TokenRParen); err != nil {
		return RPC{}, err
	}

	retType, err := p.parseType()
	if err != nil {
		return RPC{}, err
	}

	return RPC{Name: name.Value, Parameters: params, Returns: retType}, nil
}

func (p *Parser) parseField() (Field, error) {
	name, err := p.expect(lexer.TokenIdentifier)
	if err != nil {
		return Field{}, err
	}
	if _, err := p.expect(lexer.TokenColon); err != nil {
		return Field{}, err
	}
	fieldType, err := p.parseType()
	if err != nil {
		return Field{}, err
	}
	return Field{Name: name.Value, Type: fieldType}, nil
}

func (p *Parser) parseType() (TypeRef, error) {
	name, err := p.expect(lexer.TokenIdentifier)
	if err != nil {
		return TypeRef{}, err
	}
	switch name.Value {
	case "list":
		if _, err := p.expect(lexer.TokenLBrack); err != nil {
			return TypeRef{}, err
		}
		elem, err := p.parseType()
		if err != nil {
			return TypeRef{}, err
		}
		if _, err := p.expect(lexer.TokenRBrack); err != nil {
			return TypeRef{}, err
		}
		typeRef := TypeRef{Kind: TypeList, Elem: &elem}
		if p.match(lexer.TokenOptional) {
			typeRef.Optional = true
		}
		if err := ValidateType(typeRef); err != nil {
			return TypeRef{}, err
		}
		return typeRef, nil
	case "map":
		if _, err := p.expect(lexer.TokenLBrack); err != nil {
			return TypeRef{}, err
		}
		value, err := p.parseType()
		if err != nil {
			return TypeRef{}, err
		}
		if _, err := p.expect(lexer.TokenRBrack); err != nil {
			return TypeRef{}, err
		}
		typeRef := TypeRef{Kind: TypeMap, Value: &value}
		if p.match(lexer.TokenOptional) {
			typeRef.Optional = true
		}
		if err := ValidateType(typeRef); err != nil {
			return TypeRef{}, err
		}
		return typeRef, nil
	default:
		typeRef := TypeRef{Kind: TypeIdent, Name: name.Value}
		if p.match(lexer.TokenOptional) {
			typeRef.Optional = true
		}
		if err := ValidateType(typeRef); err != nil {
			return TypeRef{}, err
		}
		return typeRef, nil
	}
}

func (p *Parser) match(tt lexer.TokenType) bool {
	if p.atEnd() || p.peek().Type != tt {
		return false
	}
	p.pos++
	return true
}

func (p *Parser) expect(tt lexer.TokenType) (lexer.Token, error) {
	if p.atEnd() {
		return lexer.Token{}, fmt.Errorf("unexpected end of input, expected %s", lexer.TokenTypeName(tt))
	}
	token := p.peek()
	if token.Type != tt {
		return lexer.Token{}, p.unexpected(lexer.TokenTypeName(tt))
	}
	p.pos++
	return token, nil
}

func (p *Parser) atEnd() bool {
	return p.pos >= len(p.tokens)
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) unexpected(expected string) error {
	if p.atEnd() {
		return fmt.Errorf("unexpected end of input, expected %s", expected)
	}
	token := p.peek()
	return fmt.Errorf("unexpected token %q at line %d, column %d, expected %s", token.Value, token.Line, token.Col, expected)
}

func formatType(t TypeRef) string {
	var b strings.Builder
	switch t.Kind {
	case TypeList:
		b.WriteString("list[")
		if t.Elem != nil {
			b.WriteString(formatType(*t.Elem))
		}
		b.WriteString("]")
	case TypeMap:
		b.WriteString("map[")
		if t.Value != nil {
			b.WriteString(formatType(*t.Value))
		}
		b.WriteString("]")
	default:
		b.WriteString(t.Name)
	}
	if t.Optional {
		b.WriteString("?")
	}
	return b.String()
}

func writeTreeLine(b *strings.Builder, depth int, text string) {
	for i := 0; i < depth; i++ {
		b.WriteString("  ")
	}
	b.WriteString(text)
	b.WriteString("\n")
}

func ValidateType(t TypeRef) error {
	switch t.Kind {
	case TypeList:
		if t.Elem == nil {
			return fmt.Errorf("list type missing element")
		}
		return ValidateType(*t.Elem)
	case TypeMap:
		if t.Value == nil {
			return fmt.Errorf("map type missing value")
		}
		return ValidateType(*t.Value)
	case TypeIdent:
		if t.Name == "" {
			return fmt.Errorf("identifier type is empty")
		}
	}
	return nil
}

func ValidateSchema(schema *Schema) error {
	if schema == nil {
		return fmt.Errorf("schema is nil")
	}
	models := make(map[string]struct{}, len(schema.Models))
	for _, model := range schema.Models {
		if model.Name == "" {
			return fmt.Errorf("model name is empty")
		}
		if _, exists := models[model.Name]; exists {
			return fmt.Errorf("duplicate model %q", model.Name)
		}
		models[model.Name] = struct{}{}
	}
	rpcs := make(map[string]struct{}, len(schema.RPCs))
	for _, rpc := range schema.RPCs {
		if rpc.Name == "" {
			return fmt.Errorf("rpc name is empty")
		}
		if _, exists := rpcs[rpc.Name]; exists {
			return fmt.Errorf("duplicate rpc %q", rpc.Name)
		}
		rpcs[rpc.Name] = struct{}{}
	}
	for _, model := range schema.Models {
		fields := make(map[string]struct{}, len(model.Fields))
		for _, field := range model.Fields {
			if field.Name == "" {
				return fmt.Errorf("model %q has empty field name", model.Name)
			}
			if _, exists := fields[field.Name]; exists {
				return fmt.Errorf("model %q has duplicate field %q", model.Name, field.Name)
			}
			fields[field.Name] = struct{}{}
			if err := validateTypeRef(field.Type, models); err != nil {
				return fmt.Errorf("model %q field %q: %w", model.Name, field.Name, err)
			}
		}
	}
	for _, rpc := range schema.RPCs {
		params := make(map[string]struct{}, len(rpc.Parameters))
		for _, param := range rpc.Parameters {
			if param.Name == "" {
				return fmt.Errorf("rpc %q has empty parameter name", rpc.Name)
			}
			if _, exists := params[param.Name]; exists {
				return fmt.Errorf("rpc %q has duplicate parameter %q", rpc.Name, param.Name)
			}
			params[param.Name] = struct{}{}
			if err := validateTypeRef(param.Type, models); err != nil {
				return fmt.Errorf("rpc %q parameter %q: %w", rpc.Name, param.Name, err)
			}
		}
		if err := validateTypeRef(rpc.Returns, models); err != nil {
			return fmt.Errorf("rpc %q returns: %w", rpc.Name, err)
		}
	}
	return nil
}

func validateTypeRef(t TypeRef, models map[string]struct{}) error {
	if err := ValidateType(t); err != nil {
		return err
	}
	switch t.Kind {
	case TypeList:
		if t.Elem == nil {
			return fmt.Errorf("list type missing element")
		}
		return validateTypeRef(*t.Elem, models)
	case TypeMap:
		if t.Value == nil {
			return fmt.Errorf("map type missing value")
		}
		return validateTypeRef(*t.Value, models)
	case TypeIdent:
		if isBuiltinType(t.Name) {
			return nil
		}
		if _, ok := models[t.Name]; !ok {
			return fmt.Errorf("unknown type %q", t.Name)
		}
	}
	return nil
}

func isBuiltinType(name string) bool {
	switch name {
	case "string", "int", "bool":
		return true
	default:
		return false
	}
}
