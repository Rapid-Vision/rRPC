package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	gogen "github.com/Rapid-Vision/rRPC/internal/gen/go"
	pygen "github.com/Rapid-Vision/rRPC/internal/gen/python"
	"github.com/Rapid-Vision/rRPC/internal/parser"
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
	if err := os.WriteFile(initPath, []byte(pygen.GeneratePythonInit(schema)), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
