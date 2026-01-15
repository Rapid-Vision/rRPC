package cmd

import (
	"fmt"
	"os"

	pygen "github.com/Rapid-Vision/rRPC/internal/gen/python"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate Python client code from a schema",
	RunE:  RunClientCmd,
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func RunClientCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	schemaPath := args[0]
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}
	schema, err := parser.Parse(string(data))
	if err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}
	code, err := pygen.GenerateClient(schema)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	_, err = fmt.Fprint(cmd.OutOrStdout(), code)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
