package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format [schema]",
	Short: "Format a schema",
	RunE:  RunFormatCmd,
}

var (
	formatWrite bool
)

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().BoolVarP(&formatWrite, "write", "w", false, "Rewrite input file in-place")
}

func RunFormatCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("expected zero or one schema path argument")
	}

	var (
		data []byte
		err  error
	)
	if len(args) == 1 {
		data, err = os.ReadFile(args[0])
	} else {
		data, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	schema, err := parser.Parse(string(data))
	if err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}

	formatted, err := parser.FormatSchema(schema)
	if err != nil {
		return fmt.Errorf("format schema: %w", err)
	}

	if formatWrite && len(args) == 1 {
		if err := os.WriteFile(args[0], []byte(formatted), 0o644); err != nil {
			return fmt.Errorf("write schema: %w", err)
		}
		return nil
	}

	if _, err := io.WriteString(os.Stdout, formatted); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
