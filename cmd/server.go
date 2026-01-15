package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	gogen "github.com/Rapid-Vision/rRPC/internal/gen/go"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Generate server code from a schema",
	RunE:  RunServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("lang", "go", "Output language")
	serverCmd.Flags().StringP("pkg", "p", "rpcserver", "Package name for generated code")
	serverCmd.Flags().StringP("output", "o", "", "Output base directory (default: .)")
	serverCmd.Flags().BoolP("force", "f", false, "Overwrite output file if it exists")
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return fmt.Errorf("read lang flag: %w", err)
	}
	if lang != "go" {
		return fmt.Errorf("unsupported language %q for server", lang)
	}
	pkg, err := cmd.Flags().GetString("pkg")
	if err != nil {
		return fmt.Errorf("read pkg flag: %w", err)
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("read output flag: %w", err)
	}
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return fmt.Errorf("read force flag: %w", err)
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
	code, err := gogen.Generate(schema, pkg)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	outputDir := output
	if outputDir == "" {
		outputDir = "."
	}
	outputPath := filepath.Join(outputDir, pkg, "server.go")
	if !force {
		if _, statErr := os.Stat(outputPath); statErr == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", outputPath)
		}
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outputPath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
