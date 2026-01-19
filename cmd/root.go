package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.0.2"

var (
	rootVersionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "rrpc",
	Short: "rRPC is a code generation tool for creating an RPC API from a schema",
	RunE:  RunRootCmd,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&rootVersionFlag, "version", "v", false, "Print version")
}

func RunRootCmd(cmd *cobra.Command, args []string) error {
	if rootVersionFlag {
		fmt.Println(version)
	} else {
		return cmd.Help()
	}

	return nil
}
