package utils

import "unicode"

type IdentifierName struct {
	original string
	parts    []string
}

func NewIdentifierName(name string) IdentifierName {
	return IdentifierName{
		original: name,
		parts:    splitIdentifier(name),
	}
}

func (n IdentifierName) PascalCase() string {
	if len(n.parts) == 0 {
		return ""
	}
	out := make([]rune, 0, len(n.original))
	for _, part := range n.parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		out = append(out, unicode.ToUpper(runes[0]))
		for i := 1; i < len(runes); i++ {
			out = append(out, unicode.ToLower(runes[i]))
		}
	}
	return string(out)
}

func (n IdentifierName) SnakeCase() string {
	if len(n.parts) == 0 {
		return ""
	}
	out := make([]rune, 0, len(n.original)+len(n.parts))
	for i, part := range n.parts {
		if part == "" {
			continue
		}
		if i > 0 {
			out = append(out, '_')
		}
		for _, r := range part {
			out = append(out, unicode.ToLower(r))
		}
	}
	return string(out)
}

func splitIdentifier(name string) []string {
	var parts []string
	var buf []rune
	runes := []rune(name)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			if len(buf) > 0 {
				parts = append(parts, string(buf))
				buf = buf[:0]
			}
			continue
		}

		if len(buf) > 0 {
			prev := buf[len(buf)-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower)) {
				parts = append(parts, string(buf))
				buf = buf[:0]
			}
		}
		buf = append(buf, r)
	}
	if len(buf) > 0 {
		parts = append(parts, string(buf))
	}
	return parts
}
