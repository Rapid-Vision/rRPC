package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootVersionFlag bool
)

var rootCmd = &cobra.Command{
	Use:     "rrpc",
	Short:   "rRPC is a code generation tool for creating an RPC API from a schema",
	Version: "0.0.8",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&rootVersionFlag, "version", "v", false, "Print version")

	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
