package main

import (
	"fmt"
	"os"

	"github.com/Rapid-Vision/rRPC/internal/lexer"
)

func main() {
	// cmd.Execute()

	data, err := os.ReadFile("examples/example.rrpc")
	if err != nil {
		fmt.Print(err)
		return
	}

	lex := lexer.NewLexer(string(data))
	tokens := lex.Tokenize()

	for i, tok := range tokens {
		fmt.Println(i, tok)
	}
}
