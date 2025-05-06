package main

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "go-indexer",
		Short: "A file indexing and search tool using Redis",
	}
	rootCmd.AddCommand(NewIndexCmd())
	rootCmd.AddCommand(NewSearchCmd())
	return rootCmd
}