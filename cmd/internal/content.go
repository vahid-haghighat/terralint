package internal

import (
	"github.com/vahid-haghighat/terralint/parser"
	"os"
)

func getFormattedContent(filePath string) ([]byte, error) {
	root, err := parser.ParseTerraformFile(filePath)
	root = root
	if err != nil {
		return nil, err
	}

	// For now, just return the original file content
	// In the future, we can implement formatting based on the parsed AST
	return os.ReadFile(filePath)
}
