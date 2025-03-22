package debug

import (
	"context"
	"fmt"
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	sitterhcl "github.com/smacker/go-tree-sitter/hcl"
)

func ParseFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Initialize tree-sitter parser
	parser := sitter.NewParser()
	parser.SetLanguage(sitterhcl.GetLanguage())

	// Parse the input
	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		fmt.Printf("Failed to parse: %v\n", err)
		os.Exit(1)
	}

	// Start with the root node
	rootNode := tree.RootNode()
	fmt.Printf("Root node type: %s\n", rootNode.Type())

	// Walk the tree and print node types
	walkTree(rootNode, content, 0)
}

func walkTree(node *sitter.Node, content []byte, depth int) {
	indent := strings.Repeat("  ", depth)
	nodeText := string(content[node.StartByte():node.EndByte()])
	// Truncate long text
	if len(nodeText) > 50 {
		nodeText = nodeText[:47] + "..."
	}
	// Replace newlines with \n for display
	nodeText = strings.ReplaceAll(nodeText, "\n", "\\n")

	fmt.Printf("%sNode Type: %s, Line: %d, Text: %q\n",
		indent, node.Type(), node.StartPoint().Row+1, nodeText)

	// Print field names for children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		fieldName := node.FieldNameForChild(i)
		childText := string(content[child.StartByte():child.EndByte()])
		if len(childText) > 30 {
			childText = childText[:27] + "..."
		}
		childText = strings.ReplaceAll(childText, "\n", "\\n")
		fmt.Printf("%s  Child %d: Type: %s, Field: %s, Text: %q\n",
			indent, i, child.Type(), fieldName, childText)
	}

	// Recursively walk named children
	for i := 0; i < int(node.NamedChildCount()); i++ {
		walkTree(node.NamedChild(i), content, depth+1)
	}
}
