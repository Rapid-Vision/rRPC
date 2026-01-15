package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Rapid-Vision/rRPC/internal/lexer"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Dump lexer tokens or parser AST for a schema",
	RunE:  RunDebugCmd,
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().String("stage", "tokens", "Debug stage: tokens or ast")
}

func RunDebugCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	stage, err := cmd.Flags().GetString("stage")
	if err != nil {
		return fmt.Errorf("read stage flag: %w", err)
	}
	schemaPath := args[0]
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	switch stage {
	case "tokens", "tok", "lex", "lexer":
		tokens, err := lexer.NewLexer(string(data)).Tokenize()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)
		for _, tok := range tokens {
			_, err := fmt.Fprintf(w, "(%d:%d)\t%s\t%q\n", tok.Line, tok.Col, lexer.TokenTypeName(tok.Type), tok.Value)
			if err != nil {
				return fmt.Errorf("write output: %w", err)
			}
		}
		w.Flush()
	case "ast", "parser":
		schema, err := parser.Parse(string(data))
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(cmd.OutOrStdout(), schema.Dump())
		if err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	default:
		return fmt.Errorf("unsupported stage %q (use tokens or ast)", stage)
	}
	return nil
}
