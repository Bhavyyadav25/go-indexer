package tokenizer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kljensen/snowball"
	stopwords "github.com/zoomio/stopwords"
)

var validExts = map[string]bool{
	".txt":  true,
	".md":   true,
	".go":   true,
	".html": true,
	".log":  true,
}

func IsTextFile(path string) bool {
	return validExts[filepath.Ext(path)]
}

func ProcessFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	tokens := tokenize(content)
	return processTokens(tokens), nil
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	})
}

func processTokens(tokens []string) []string {
	var result []string
	for _, token := range tokens {
		token = strings.ToLower(token)
		if stopwords.IsStopWord(token) {
			continue
		}
		stemmed, err := snowball.Stem(token, "english", true)
		if err == nil && len(stemmed) > 2 {
			result = append(result, stemmed)
		}
	}
	return result
}
