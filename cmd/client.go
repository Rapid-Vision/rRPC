package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pygen "github.com/Rapid-Vision/rRPC/internal/gen/python"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/Rapid-Vision/rRPC/internal/utils"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate client code from a schema",
	RunE:  RunClientCmd,
}

var (
	clientLang  string
	clientPkg   string
	clientOut   string
	clientForce bool
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientLang, "lang", "python", "Output language")
	clientCmd.Flags().StringVarP(&clientPkg, "pkg", "p", "rpc_client", "Python package name for generated code")
	clientCmd.Flags().StringVarP(&clientOut, "output", "o", "", "Output base directory (default: .)")
	clientCmd.Flags().BoolVarP(&clientForce, "force", "f", false, "Overwrite output file if it exists")
}

func RunClientCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	if clientLang != "python" && clientLang != "py" {
		return fmt.Errorf("unsupported language %q for client", clientLang)
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
	outputDir := clientOut
	if outputDir == "" {
		outputDir = "."
	}
	baseDir := filepath.Join(outputDir, "generated", clientPkg)
	clientPath := filepath.Join(baseDir, "client.py")
	initPath := filepath.Join(baseDir, "__init__.py")
	if !clientForce {
		if _, statErr := os.Stat(clientPath); statErr == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", clientPath)
		}
		if _, statErr := os.Stat(initPath); statErr == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", initPath)
		}
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(clientPath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	if err := os.WriteFile(initPath, []byte(buildPythonInit(schema)), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

func buildPythonInit(schema *parser.Schema) string {
	var b strings.Builder
	b.WriteString("from .client import RPCClient\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("from .client import ")
		b.WriteString(className)
		b.WriteString("\n")
	}
	b.WriteString("\n__all__ = [\n")
	b.WriteString("    \"RPCClient\",\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("    \"")
		b.WriteString(className)
		b.WriteString("\",\n")
	}
	b.WriteString("]\n")
	return b.String()
}
