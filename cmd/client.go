package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gogen "github.com/Rapid-Vision/rRPC/internal/gen/go"
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
	clientLang   string
	clientPkg    string
	clientOut    string
	clientForce  bool
	clientPrefix string
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientLang, "lang", "python", "Output language")
	clientCmd.Flags().StringVarP(&clientPkg, "pkg", "p", "rpc_client", "Python package name for generated code")
	clientCmd.Flags().StringVarP(&clientOut, "output", "o", ".", "Output base directory")
	clientCmd.Flags().BoolVarP(&clientForce, "force", "f", false, "Overwrite output file if it exists")
	clientCmd.Flags().StringVar(&clientPrefix, "prefix", "rpc", "URL path prefix (empty for none)")
}

func RunClientCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	if clientLang != "python" && clientLang != "py" && clientLang != "go" {
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
	outputDir := clientOut
	if outputDir == "" {
		outputDir = "."
	}
	baseDir := filepath.Join(outputDir, clientPkg)
	if clientLang == "go" {
		code, err := gogen.GenerateClientWithPrefix(schema, clientPkg, clientPrefix)
		if err != nil {
			return fmt.Errorf("generate code: %w", err)
		}
		clientPath := filepath.Join(baseDir, "client.go")
		if !clientForce {
			if _, statErr := os.Stat(clientPath); statErr == nil {
				return fmt.Errorf("output file exists: %s (use --force to overwrite)", clientPath)
			}
		}
		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
		if err := os.WriteFile(clientPath, []byte(code), 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}

	code, err := pygen.GenerateClientWithPrefix(schema, clientPrefix)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
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
	b.WriteString("from .client import RPCError\n")
	b.WriteString("from .client import RPCErrorException\n")
	b.WriteString("from .client import CustomRPCError\n")
	b.WriteString("from .client import ValidationRPCError\n")
	b.WriteString("from .client import InputRPCError\n")
	b.WriteString("from .client import UnauthorizedRPCError\n")
	b.WriteString("from .client import ForbiddenRPCError\n")
	b.WriteString("from .client import NotImplementedRPCError\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("from .client import ")
		b.WriteString(className)
		b.WriteString("\n")
	}
	b.WriteString("\n__all__ = [\n")
	b.WriteString("    \"RPCClient\",\n")
	b.WriteString("    \"RPCError\",\n")
	b.WriteString("    \"RPCErrorException\",\n")
	b.WriteString("    \"CustomRPCError\",\n")
	b.WriteString("    \"ValidationRPCError\",\n")
	b.WriteString("    \"InputRPCError\",\n")
	b.WriteString("    \"UnauthorizedRPCError\",\n")
	b.WriteString("    \"ForbiddenRPCError\",\n")
	b.WriteString("    \"NotImplementedRPCError\",\n")
	for _, model := range schema.Models {
		className := utils.NewIdentifierName(model.Name).PascalCase() + "Model"
		b.WriteString("    \"")
		b.WriteString(className)
		b.WriteString("\",\n")
	}
	b.WriteString("]\n")
	return b.String()
}
