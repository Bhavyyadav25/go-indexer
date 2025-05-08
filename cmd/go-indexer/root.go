package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-indexer",
	Short: "A file indexing and search tool using Redis",
}

func NewRootCmd() {
	rootCmd.AddCommand(NewIndexCmd())
	rootCmd.AddCommand(NewSearchCmd())
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
