package lexer

import (
	"regexp"
	"strings"
)

type TokenType int

const (
	TokenIgn TokenType = iota
	TokenModel
	TokenRpc
	TokenIdentifier
	TokenOptional
	TokenColon
	TokenLParen
	TokenRParen
	TokenLBrack
	TokenRBrack
	TokenLBrace
	TokenRBrace
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	text string
}

type rule struct {
	Name  string
	Regex string
	Type  TokenType
}

var rules = [...]rule{
	{
		Name:  "ws",
		Regex: `(?P<ws>\s+)`,
		Type:  TokenIgn,
	},
	{
		Name:  "comment",
		Regex: `(?P<comment>#[^\n]*)`,
		Type:  TokenIgn,
	},
	{
		Name:  "model",
		Regex: `(?P<model>model)\b`,
		Type:  TokenModel,
	},
	{
		Name:  "rpc",
		Regex: `(?P<rpc>rpc)\b`,
		Type:  TokenRpc,
	},
	{
		Name:  "ident",
		Regex: `(?P<ident>[A-Za-z_][A-Za-z0-9_]*)`,
		Type:  TokenIdentifier,
	},
	{
		Name:  "colon",
		Regex: `(?P<colon>:)`,
		Type:  TokenColon,
	},
	{
		Name:  "lparen",
		Regex: `(?P<lparen>\()`,
		Type:  TokenLParen,
	},
	{
		Name:  "rparen",
		Regex: `(?P<rparen>\))`,
		Type:  TokenRParen,
	},
	{
		Name:  "lbrack",
		Regex: `(?P<lbrack>\[)`,
		Type:  TokenLBrack,
	},
	{
		Name:  "rbrack",
		Regex: `(?P<rbrack>\])`,
		Type:  TokenRBrack,
	},
	{
		Name:  "lbrace",
		Regex: `(?P<lbrace>\{)`,
		Type:  TokenLBrace,
	},
	{
		Name:  "rbrace",
		Regex: `(?P<rbrace>\})`,
		Type:  TokenRBrace,
	},
	{
		Name:  "optional",
		Regex: `(?P<optional>\?)`,
		Type:  TokenOptional,
	},
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		text: text,
	}
}

func (l *Lexer) Tokenize() (tokens []Token) {
	patterns := make([]string, 0, len(rules))
	nameToType := make(map[string]TokenType, len(rules))
	for _, rule := range rules {
		patterns = append(patterns, rule.Regex)
		nameToType[rule.Name] = rule.Type
	}
	re := regexp.MustCompile(`^(?:` + strings.Join(patterns, "|") + `)`)
	names := re.SubexpNames()

	text := l.text
	for i := 0; i < len(text); {
		loc := re.FindStringSubmatchIndex(text[i:])
		if loc == nil || loc[0] != 0 {
			i++
			continue
		}

		var tokenType TokenType
		var value string
		for group := 1; group < len(loc)/2; group++ {
			start := loc[group*2]
			end := loc[group*2+1]
			if start == -1 || end == -1 {
				continue
			}
			value = text[i+start : i+end]
			name := names[group]
			matchedType, ok := nameToType[name]
			if !ok {
				value = ""
				break
			}
			tokenType = matchedType
			break
		}

		if value != "" && tokenType != TokenIgn {
			tokens = append(tokens, Token{Type: tokenType, Value: value})
		}
		i += loc[1]
	}
	return tokens
}
