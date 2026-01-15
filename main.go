package main

import (
	"fmt"
	"os"

	"github.com/Rapid-Vision/rRPC/internal/parser"
)

func main() {
	// cmd.Execute()

	data, err := os.ReadFile("examples/example.rrpc")
	if err != nil {
		fmt.Print(err)
		return
	}

	res, err := parser.Parse(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Dump())
}
