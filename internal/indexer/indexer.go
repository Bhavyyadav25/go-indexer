package indexer

import (
	"path/filepath"
)

func IsTextFile(path string) bool {
	ext := filepath.Ext(path)
	_, ok := map[string]struct{}{
		".txt":  {},
		".md":   {},
		".go":   {},
		".html": {},
		".log":  {},
	}[ext]
	return ok
}
