package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	gogen "github.com/Rapid-Vision/rRPC/internal/gen/go"
	pygen "github.com/Rapid-Vision/rRPC/internal/gen/python"
	tsgen "github.com/Rapid-Vision/rRPC/internal/gen/typescript"
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
	clientZod    bool
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientLang, "lang", "python", "Output language")
	clientCmd.Flags().StringVarP(&clientPkg, "pkg", "p", "rpcclient", "Output package/directory name for generated code")
	clientCmd.Flags().StringVarP(&clientOut, "output", "o", ".", "Output base directory")
	clientCmd.Flags().BoolVarP(&clientForce, "force", "f", false, "Overwrite output file if it exists")
	clientCmd.Flags().StringVar(&clientPrefix, "prefix", "rpc", "URL path prefix (empty for none)")
	clientCmd.Flags().BoolVar(&clientZod, "ts-zod", false, "Generate TypeScript client with zod input validation")
}

func RunClientCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	if clientLang != "python" && clientLang != "py" && clientLang != "go" && clientLang != "ts" && clientLang != "typescript" {
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
		files, err := gogen.GenerateClientWithPrefix(schema, clientPkg, clientPrefix)
		if err != nil {
			return fmt.Errorf("generate code: %w", err)
		}
		filePaths := make([]string, 0, len(files))
		for name := range files {
			filePaths = append(filePaths, filepath.Join(baseDir, name))
		}
		if !clientForce {
			for _, path := range filePaths {
				if _, statErr := os.Stat(path); statErr == nil {
					return fmt.Errorf("output file exists: %s (use --force to overwrite)", path)
				}
			}
		}
		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
		for name, contents := range files {
			outPath := filepath.Join(baseDir, name)
			if err := os.WriteFile(outPath, []byte(contents), 0o644); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
		}
		return nil
	}
	if clientLang == "ts" || clientLang == "typescript" {
		code, err := tsgen.GenerateClientWithPrefixAndZod(schema, clientPrefix, clientZod)
		if err != nil {
			return fmt.Errorf("generate code: %w", err)
		}
		clientPath := filepath.Join(baseDir, "client.ts")
		indexPath := filepath.Join(baseDir, "index.ts")
		if !clientForce {
			if _, statErr := os.Stat(clientPath); statErr == nil {
				return fmt.Errorf("output file exists: %s (use --force to overwrite)", clientPath)
			}
			if _, statErr := os.Stat(indexPath); statErr == nil {
				return fmt.Errorf("output file exists: %s (use --force to overwrite)", indexPath)
			}
		}
		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
		if err := os.WriteFile(clientPath, []byte(code), 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		if err := os.WriteFile(indexPath, []byte(tsgen.GenerateTypeScriptIndexWithZod(schema, clientZod)), 0o644); err != nil {
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
