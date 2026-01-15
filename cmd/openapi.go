package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rapid-Vision/rRPC/internal/gen/openapi"
	"github.com/Rapid-Vision/rRPC/internal/parser"
	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Generate OpenAPI spec from a schema",
	RunE:  RunOpenAPICmd,
}

var (
	openapiTitle   string
	openapiVersion string
	openapiOut     string
	openapiForce   bool
)

func init() {
	rootCmd.AddCommand(openapiCmd)
	openapiCmd.Flags().StringVar(&openapiTitle, "title", "rRPC API", "OpenAPI title")
	openapiCmd.Flags().StringVar(&openapiVersion, "version", "0.1.0", "OpenAPI version")
	openapiCmd.Flags().StringVarP(&openapiOut, "output", "o", "", "Output base directory (default: .)")
	openapiCmd.Flags().BoolVarP(&openapiForce, "force", "f", false, "Overwrite output file if it exists")
}

func RunOpenAPICmd(cmd *cobra.Command, args []string) error {
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
	spec, err := openapi.Generate(schema, openapiTitle, openapiVersion)
	if err != nil {
		return fmt.Errorf("generate openapi: %w", err)
	}
	outputDir := openapiOut
	if outputDir == "" {
		outputDir = "."
	}
	outputPath := filepath.Join(outputDir, "openapi.json")
	if !openapiForce {
		if _, statErr := os.Stat(outputPath); statErr == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", outputPath)
		}
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outputPath, []byte(spec), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
