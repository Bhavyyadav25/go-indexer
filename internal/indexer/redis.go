package indexer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bhavyyadav25/go-indexer/internal/tokenizer"
	redis "github.com/redis/go-redis"
)

var (
	ctx       = context.Background()
	indexName = "docIndex"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func CreateIndex(client *redis.Client) error {
	schema := &redis.IndexSchema{
		Name: "content",
		Type: redis.TextField,
		Options: &redis.TextFieldOptions{
			Weight:   5.0,
			Sortable: false,
		},
	}
	pathField := &redis.IndexSchema{
		Name: "path",
		Type: redis.TextField,
		Options: &redis.TextFieldOptions{
			Sortable: true,
		},
	}
	sizeField := &redis.IndexSchema{
		Name: "size",
		Type: redis.NumericField,
		Options: &redis.NumericFieldOptions{
			Sortable: true,
		},
	}
	tsField := &redis.IndexSchema{
		Name: "timestamp",
		Type: redis.NumericField,
		Options: &redis.NumericFieldOptions{
			Sortable: true,
		},
	}

	res, err := client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil || res == nil {
		_, err := client.Do(ctx, "FT.CREATE", indexName,
			"ON", "HASH",
			"PREFIX", 1, "doc:",
			"SCHEMA",
			pathField.Name, pathField.Type, "SORTABLE",
			sizeField.Name, sizeField.Type, "SORTABLE",
			tsField.Name, tsField.Type, "SORTABLE",
			schema.Name, "TEXT", schema.Options,
		).Result()
		return err
	}
	return nil
}

func IndexFile(client *redis.Client, path string) error {
	content, err := tokenizer.ProcessFile(path)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	doc := map[string]interface{}{
		"content":   strings.Join(content, " "),
		"path":      path,
		"size":      info.Size(),
		"timestamp": info.ModTime().Unix(),
	}

	key := fmt.Sprintf("doc:%s", filepath.Base(path))
	return client.HSet(ctx, key, doc).Err()
}

func Search(client *redis.Client, query string) ([]map[string]string, error) {
	res, err := client.Do(ctx, "FT.SEARCH", indexName, query,
		"SUMMARIZE", "FIELDS", 1, "content",
		"HIGHLIGHT", "FIELDS", 1, "content",
	).Result()
	if err != nil {
		return nil, err
	}

	return parseSearchResult(res), nil
}

func parseSearchResult(res interface{}) []map[string]string {
	results := []map[string]string{}
	data, ok := res.([]interface{})
	if !ok || len(data) < 1 {
		return results
	}

	for i := 1; i < len(data); i += 2 {
		doc := make(map[string]string)
		fields, ok := data[i+1].([]interface{})
		if !ok {
			continue
		}

		for j := 0; j < len(fields); j += 2 {
			key, kok := fields[j].(string)
			val, vok := fields[j+1].(string)
			if kok && vok {
				doc[key] = val
			}
		}
		results = append(results, doc)
	}
	return results
}
