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

var (
	serverLang   string
	serverPkg    string
	serverOut    string
	serverForce  bool
	serverPrefix string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverLang, "lang", "go", "Output language")
	serverCmd.Flags().StringVarP(&serverPkg, "pkg", "p", "rpcserver", "Package name for generated code")
	serverCmd.Flags().StringVarP(&serverOut, "output", "o", ".", "Output base directory")
	serverCmd.Flags().BoolVarP(&serverForce, "force", "f", false, "Overwrite output file if it exists")
	serverCmd.Flags().StringVar(&serverPrefix, "prefix", "rpc", "URL path prefix (empty for none)")
}

func RunServerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected schema path argument")
	}
	if serverLang != "go" {
		return fmt.Errorf("unsupported language %q for server", serverLang)
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
	files, err := gogen.GenerateWithPrefix(schema, serverPkg, serverPrefix)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	outputDir := serverOut
	if outputDir == "" {
		outputDir = "."
	}
	filePaths := make([]string, 0, len(files))
	for name := range files {
		filePaths = append(filePaths, filepath.Join(outputDir, serverPkg, name))
	}
	if !serverForce {
		for _, path := range filePaths {
			if _, statErr := os.Stat(path); statErr == nil {
				return fmt.Errorf("output file exists: %s (use --force to overwrite)", path)
			}
		}
	}
	outputBase := filepath.Join(outputDir, serverPkg)
	if err := os.MkdirAll(outputBase, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	for name, contents := range files {
		outPath := filepath.Join(outputBase, name)
		if err := os.WriteFile(outPath, []byte(contents), 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}
	return nil
}
