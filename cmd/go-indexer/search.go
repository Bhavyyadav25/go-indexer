package main

import (
	"fmt"

	"github.com/Bhavyyadav25/go-indexer/internal/indexer"
	"github.com/spf13/cobra"
)

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search indexed documents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := args[0]
			client := indexer.NewRedisClient()
			defer client.Close()

			res, err := indexer.Search(client, query)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, doc := range res {
				fmt.Printf("Path: %s\nScore: %f\nSummary: %s\n\n",
					doc["path"], doc["score"], doc["summary"])
			}
		},
	}
	return cmd
}
