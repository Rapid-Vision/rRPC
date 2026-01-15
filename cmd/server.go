package cmd

import (
	"fmt"
	"os"

	gogen "github.com/Rapid-Vision/rRPC/internal/gen/go"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Generate a go server from schema",
	RunE:  RunServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
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
	code, err := gogen.Generate(schema, "rpc")
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	_, err = fmt.Fprint(cmd.OutOrStdout(), code)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
