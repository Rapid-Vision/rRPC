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
	for i, model := range s.Models {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("model ")
		b.WriteString(model.Name)
		b.WriteString(" {\n")
		for _, field := range model.Fields {
			b.WriteString("  ")
			b.WriteString(field.Name)
			b.WriteString(": ")
			b.WriteString(formatType(field.Type))
			b.WriteString("\n")
		}
		b.WriteString("}\n")
	}

	for i, rpc := range s.RPCs {
		if len(s.Models) > 0 || i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("rpc ")
		b.WriteString(rpc.Name)
		b.WriteString("(")
		for idx, param := range rpc.Parameters {
			if idx > 0 {
				b.WriteString(", ")
			}
			b.WriteString(param.Name)
			b.WriteString(": ")
			b.WriteString(formatType(param.Type))
		}
		b.WriteString(") ")
		b.WriteString(formatType(rpc.Returns))
		b.WriteString("\n")
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
	Key      *TypeRef
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
	tokens := lexer.NewLexer(text).Tokenize()
	p := NewParser(tokens)
	return p.parseSchema()
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
		key, err := p.parseType()
		if err != nil {
			return TypeRef{}, err
		}
		if _, err := p.expect(lexer.TokenComma); err != nil {
			return TypeRef{}, err
		}
		value, err := p.parseType()
		if err != nil {
			return TypeRef{}, err
		}
		if _, err := p.expect(lexer.TokenRBrack); err != nil {
			return TypeRef{}, err
		}
		typeRef := TypeRef{Kind: TypeMap, Key: &key, Value: &value}
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
	return fmt.Errorf("unexpected token %q, expected %s", p.peek().Value, expected)
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
		if t.Key != nil {
			b.WriteString(formatType(*t.Key))
		}
		b.WriteString(", ")
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

func ValidateType(t TypeRef) error {
	switch t.Kind {
	case TypeList:
		if t.Elem == nil {
			return fmt.Errorf("list type missing element")
		}
		return ValidateType(*t.Elem)
	case TypeMap:
		if t.Key == nil || t.Value == nil {
			return fmt.Errorf("map type missing key or value")
		}
		if !isValidMapKey(*t.Key) {
			return fmt.Errorf("map key type must be a non-optional string or int")
		}
		if err := ValidateType(*t.Value); err != nil {
			return err
		}
	case TypeIdent:
		if t.Name == "" {
			return fmt.Errorf("identifier type is empty")
		}
	}
	return nil
}

func isValidMapKey(t TypeRef) bool {
	if t.Kind != TypeIdent {
		return false
	}
	if t.Optional {
		return false
	}
	switch t.Name {
	case "string", "int":
		return true
	default:
		return false
	}
}
