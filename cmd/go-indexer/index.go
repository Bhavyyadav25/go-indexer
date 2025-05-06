package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Bhavyyadav25/go-indexer/internal/indexer"
	"github.com/Bhavyyadav25/go-indexer/internal/tokenizer"
	"github.com/Bhavyyadav25/go-indexer/internal/worker"
	"github.com/spf13/cobra"
)

func NewIndexCmd() *cobra.Command {
	var workers int
	cmd := &cobra.Command{
		Use:   "index [directory]",
		Short: "Index files in a directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			redisClient := indexer.NewRedisClient()
			defer redisClient.Close()

			pool := worker.NewPool(workers, func(path string) error {
				return indexer.IndexFile(redisClient, path)
			})

			go func() {
				err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && tokenizer.IsTextFile(path) {
						pool.Jobs <- path
					}
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}
				close(pool.Jobs)
			}()

			pool.Wait()
			fmt.Println("Indexing completed")
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 4, "Number of worker goroutines")
	return cmd
}
