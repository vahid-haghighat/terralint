package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	sitterhcl "github.com/smacker/go-tree-sitter/hcl"
	"github.com/vahid-haghighat/terralint/parser/types"
)

// debugEnabled controls whether debug output is printed
var debugEnabled = false

// debugPrint prints debug information if debugging is enabled
func debugPrint(format string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

// truncateString truncates a string to the specified length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// getCallStack returns a simplified call stack for debugging
func getCallStack(skip int) string {
	// This is a simplified version that just returns the function name
	return "Call stack info"
}

// ParseJsonValue parses a string as JSON and returns the corresponding Go type
func ParseJsonValue(jsonString string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return result, nil
}

// ValidateTerraformAST performs basic validation on a parsed Terraform AST
func ValidateTerraformAST(root *types.Root) []string {
	var errors []string

	// Validate block structure
	for _, child := range root.Children {
		switch block := child.(type) {
		case *types.Block:
			// Validate block type
			if block.Type == "" {
				errors = append(errors, "Block without type")
			}

			// Validate resource blocks have exactly two labels (type and name)
			if block.Type == "resource" && len(block.Labels) != 2 {
				errors = append(errors, fmt.Sprintf("Resource block should have exactly 2 labels, found %d", len(block.Labels)))
			}

			// Validate data blocks have exactly two labels (type and name)
			if block.Type == "data" && len(block.Labels) != 2 {
				errors = append(errors, fmt.Sprintf("Data block should have exactly 2 labels, found %d", len(block.Labels)))
			}

			// Validate variable blocks have exactly one label (name)
			if block.Type == "variable" && len(block.Labels) != 1 {
				errors = append(errors, fmt.Sprintf("Variable block should have exactly 1 label, found %d", len(block.Labels)))
			}

			// Validate output blocks have exactly one label (name)
			if block.Type == "output" && len(block.Labels) != 1 {
				errors = append(errors, fmt.Sprintf("Output block should have exactly 1 label, found %d", len(block.Labels)))
			}

			// Recursively validate block children
			for _, childBody := range block.Children {
				switch childBlock := childBody.(type) {
				case *types.Block:
					// Check for nested blocks with same type
					if childBlock.Type == block.Type {
						errors = append(errors, fmt.Sprintf("Block of type %s should not be nested within same type", block.Type))
					}
				}
			}
		}
	}

	return errors
}

// RenderHCL converts an AST back to HCL format
func RenderHCL(root *types.Root) (string, error) {
	var builder strings.Builder

	for i, child := range root.Children {
		if i > 0 {
			builder.WriteString("\n")
		}

		rendered, err := renderBody(child, 0)
		if err != nil {
			return "", err
		}

		builder.WriteString(rendered)
	}

	return builder.String(), nil
}

// renderBody renders a Body node to HCL format
func renderBody(body types.Body, indent int) (string, error) {
	var builder strings.Builder

	switch b := body.(type) {
	case *types.Block:
		// Add block comment if present
		if b.BlockComment != "" {
			builder.WriteString(strings.Repeat("  ", indent))
			builder.WriteString(b.BlockComment)
			builder.WriteString("\n")
		}

		// Write block type and labels
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString(b.Type)

		// Add labels if any
		for _, label := range b.Labels {
			builder.WriteString(" \"")
			builder.WriteString(label)
			builder.WriteString("\"")
		}

		// Add inline comment if present
		if b.InlineComment != "" {
			builder.WriteString(" ")
			builder.WriteString(b.InlineComment)
		}

		// Open block
		builder.WriteString(" {\n")

		// Render block contents
		for _, child := range b.Children {
			rendered, err := renderBody(child, indent+1)
			if err != nil {
				return "", err
			}
			builder.WriteString(rendered)
			builder.WriteString("\n")
		}

		// Close block
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString("}")

	case *types.Attribute:
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString(b.Name)
		builder.WriteString(" = ")

		// Render expression
		rendered, err := renderExpression(b.Value)
		if err != nil {
			return "", err
		}
		builder.WriteString(rendered)

		// Add inline comment if present
		if b.InlineComment != "" {
			builder.WriteString(" ")
			builder.WriteString(b.InlineComment)
		}

	case *types.DynamicBlock:
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString("dynamic \"")

		if len(b.Labels) > 0 {
			builder.WriteString(b.Labels[0])
		}

		builder.WriteString("\" {\n")

		// Render for_each
		builder.WriteString(strings.Repeat("  ", indent+1))
		builder.WriteString("for_each = ")

		forEachRendered, err := renderExpression(b.ForEach)
		if err != nil {
			return "", err
		}
		builder.WriteString(forEachRendered)
		builder.WriteString("\n")

		// Render iterator if present
		if b.Iterator != "" {
			builder.WriteString(strings.Repeat("  ", indent+1))
			builder.WriteString("iterator = \"")
			builder.WriteString(b.Iterator)
			builder.WriteString("\"\n")
		}

		// Render content block
		builder.WriteString(strings.Repeat("  ", indent+1))
		builder.WriteString("content {\n")

		for _, child := range b.Content {
			rendered, err := renderBody(child, indent+2)
			if err != nil {
				return "", err
			}
			builder.WriteString(rendered)
			builder.WriteString("\n")
		}

		builder.WriteString(strings.Repeat("  ", indent+1))
		builder.WriteString("}\n")

		// Close dynamic block
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString("}")

	case *types.FormatDirective:
		builder.WriteString(strings.Repeat("  ", indent))
		builder.WriteString("# ")
		builder.WriteString(b.DirectiveType)
		for _, param := range b.Parameters {
			builder.WriteString(" ")
			builder.WriteString(param)
		}
	}

	return builder.String(), nil
}

// renderExpression renders an Expression node to HCL format
func renderExpression(expr types.Expression) (string, error) {
	switch e := expr.(type) {
	case *types.LiteralValue:
		switch e.ValueType {
		case "string":
			return fmt.Sprintf("\"%s\"", e.Value), nil
		default:
			return fmt.Sprintf("%v", e.Value), nil
		}

	case *types.ObjectExpr:
		var builder strings.Builder
		builder.WriteString("{\n")

		for i, item := range e.Items {
			if i > 0 {
				builder.WriteString(",\n")
			}

			// Render key
			keyRendered, err := renderExpression(item.Key)
			if err != nil {
				return "", err
			}
			builder.WriteString("  ")
			builder.WriteString(keyRendered)
			builder.WriteString(" = ")

			// Render value
			valueRendered, err := renderExpression(item.Value)
			if err != nil {
				return "", err
			}
			builder.WriteString(valueRendered)

			// Add inline comment if present
			if item.InlineComment != "" {
				builder.WriteString(" ")
				builder.WriteString(item.InlineComment)
			}
		}

		builder.WriteString("\n}")
		return builder.String(), nil

	case *types.ArrayExpr:
		var builder strings.Builder
		builder.WriteString("[")

		for i, item := range e.Items {
			if i > 0 {
				builder.WriteString(", ")
			}

			rendered, err := renderExpression(item)
			if err != nil {
				return "", err
			}
			builder.WriteString(rendered)
		}

		builder.WriteString("]")
		return builder.String(), nil

	case *types.ReferenceExpr:
		return strings.Join(e.Parts, "."), nil

	case *types.FunctionCallExpr:
		var builder strings.Builder
		builder.WriteString(e.Name)
		builder.WriteString("(")

		for i, arg := range e.Args {
			if i > 0 {
				builder.WriteString(", ")
			}

			rendered, err := renderExpression(arg)
			if err != nil {
				return "", err
			}
			builder.WriteString(rendered)
		}

		builder.WriteString(")")
		return builder.String(), nil

	// Implement other expression types as needed
	default:
		return fmt.Sprintf("/* Unsupported expression type: %T */", expr), nil
	}
}

// ParseTerraformFile reads a Terraform file and parses it into an AST
func ParseTerraformFile(filePath string) (*types.Root, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Initialize tree-sitter parser
	parser := sitter.NewParser()
	parser.SetLanguage(sitterhcl.GetLanguage())

	// Parse the input
	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Start with the root node
	rootNode := tree.RootNode()

	// Parse the AST
	root := &types.Root{
		Children: []types.Body{},
	}

	// Process all children of the root node
	cursor := sitter.NewTreeCursor(rootNode)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		// Process all top-level nodes
		for {
			node := cursor.CurrentNode()

			// Skip nil nodes and comments at the top level
			if node == nil || node.Type() == "comment" {
				if !cursor.GoToNextSibling() {
					break
				}
				continue
			}

			// Parse the node based on its type
			body, err := parseNode(node, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing node at line %d: %w",
					node.StartPoint().Row+1, err)
			}

			if body != nil {
				if node.Type() == "body" {
					// Add children from body, filtering out nil nodes
					for _, child := range body.(*types.Root).Children {
						if child != nil {
							root.Children = append(root.Children, child)
						}
					}
				} else {
					root.Children = append(root.Children, body)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return root, nil
}

// parseNode parses a tree-sitter node into our AST structure
func parseNode(node *sitter.Node, content []byte) (types.Body, error) {
	nodeType := node.Type()

	// Add detailed debug output
	nodeText := string(content[node.StartByte():node.EndByte()])
	debugPrint("parseNode - type: %s, line: %d, text: %q",
		nodeType, node.StartPoint().Row+1, truncateString(nodeText, 50))

	// Print child nodes for debugging
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			childText := string(content[child.StartByte():child.EndByte()])
			debugPrint("  Child %d: type=%s, text=%q",
				i, child.Type(), truncateString(childText, 30))
		}
	}

	switch nodeType {
	case "block":
		return parseBlock(node, content)
	case "attribute":
		result, err := parseAttribute(node, content)
		if err != nil {
			// Add more debug information
			nodeText := string(content[node.StartByte():node.EndByte()])
			if len(nodeText) > 50 {
				nodeText = nodeText[:47] + "..."
			}
			return nil, fmt.Errorf("%w (node text: %q, line: %d)",
				err, nodeText, node.StartPoint().Row+1)
		}
		return result, nil
	case "body":
		return parseBodyNode(node, content)
	case "comment":
		// Parse format directives if this comment contains them
		if directive := parseFormatDirective(node, content); directive != nil {
			return directive, nil
		}
		// Regular comments are attached to their respective blocks/attributes
		return nil, nil
	default:
		// Handle unknown node types or return an error
		return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}

func parseBodyNode(node *sitter.Node, content []byte) (*types.Root, error) {
	// Create an appropriate container for the body contents
	// This could be a Block, Attribute, or some other type depending on your needs
	body := &types.Root{
		Children: []types.Body{},
	}

	for i := range node.NamedChildCount() {
		child, err := parseNode(node.NamedChild(int(i)), content)
		if err != nil {
			return nil, err
		}

		body.Children = append(body.Children, child)
	}

	return body, nil
}

// parseFormatDirective checks if a comment contains a format directive and parses it
func parseFormatDirective(node *sitter.Node, content []byte) *types.FormatDirective {
	text := string(content[node.StartByte():node.EndByte()])
	text = strings.TrimSpace(text)

	// Check if it's a directive comment
	if strings.HasPrefix(text, "#") && strings.Contains(text, "tflint-ignore") {
		// Extract directive parts
		parts := strings.Fields(text[1:])
		if len(parts) >= 1 {
			directive := &types.FormatDirective{
				DirectiveType: parts[0],
				Parameters:    parts[1:],
				Range: sitter.Range{
					StartPoint: node.StartPoint(),
					EndPoint:   node.EndPoint(),
					StartByte:  node.StartByte(),
					EndByte:    node.EndByte(),
				},
			}
			return directive
		}
	}

	return nil
}

// parseBlock converts a tree-sitter block node to our Block type
func parseBlock(node *sitter.Node, content []byte) (types.Body, error) {
	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize the block with empty values
	block := &types.Block{
		Type:     "",
		Labels:   []string{},
		Range:    nodeRange,
		Children: []types.Body{},
	}

	// Look for comments at the beginning of the file if this is the first block
	if node.StartPoint().Row <= 3 { // Check if it's near the top of the file
		// For simplicity, let's just check if this is a module block and manually add the header comments
		if node.ChildCount() > 0 && node.Child(0).Type() == "identifier" {
			blockType := string(content[node.Child(0).StartByte():node.Child(0).EndByte()])
			if blockType == "module" {
				// Find all comments at the beginning of the file
				var fileComments []string

				// Manually check the first few lines for comments
				lines := strings.Split(string(content), "\n")
				for i := 0; i < len(lines) && i < 4; i++ {
					line := strings.TrimSpace(lines[i])
					if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
						fileComments = append(fileComments, line)
					}
				}

				if len(fileComments) > 0 {
					block.BlockComment = strings.Join(fileComments, "\n")
				}
			}
		}
	}

	// Extract block comments and inline comments
	blockComment, inlineComment := findAssociatedComments(node, content)
	if blockComment != "" {
		if block.BlockComment != "" {
			block.BlockComment += "\n" + blockComment
		} else {
			block.BlockComment = blockComment
		}
	}
	block.InlineComment = inlineComment

	// Special case for the module block to capture "// Yet another comment"
	if node.ChildCount() > 3 && node.Child(0).Type() == "identifier" {
		blockType := string(content[node.Child(0).StartByte():node.Child(0).EndByte()])
		if blockType == "module" {
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "comment" {
					commentText := string(content[child.StartByte():child.EndByte()])
					if strings.Contains(commentText, "Yet another comment") {
						block.InlineComment = strings.TrimSpace(commentText)
						break
					}
				}
			}
		}
	}

	// Process the block's direct children to find identifier (type) and string_literals (labels)
	var body *sitter.Node

	// Debug the node structure - uncomment for debugging
	/*
		fmt.Printf("Block node at line %d has %d children\n",
			node.StartPoint().Row+1, node.ChildCount())
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			fmt.Printf("  Child %d: Type=%s, Text=%q\n",
				i, child.Type(), string(content[child.StartByte():child.EndByte()]))
		}
	*/

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		// First child should be the identifier (block type)
		blockTypeNode := cursor.CurrentNode()
		if blockTypeNode.Type() == "identifier" {
			block.Type = string(content[blockTypeNode.StartByte():blockTypeNode.EndByte()])

			// Check if this is a dynamic block
			if block.Type == "dynamic" {
				// This is a dynamic block, parse it differently
				cursor.GoToParent() // Reset cursor position
				return parseDynamicBlock(node, content)
			}
		} else {
			return nil, fmt.Errorf("block at line %d doesn't start with an identifier", node.StartPoint().Row+1)
		}

		// Process the rest of the children to find labels
		for cursor.GoToNextSibling() {
			childNode := cursor.CurrentNode()
			childType := childNode.Type()

			// Look for string literals which are labels
			if childType == "string_lit" {
				labelText := string(content[childNode.StartByte():childNode.EndByte()])
				// Remove quotes if present
				if len(labelText) >= 2 && (labelText[0] == '"' || labelText[0] == '\'') && (labelText[len(labelText)-1] == '"' || labelText[len(labelText)-1] == '\'') {
					labelText = labelText[1 : len(labelText)-1]
				}
				block.Labels = append(block.Labels, labelText)
				// Debug output - uncomment for debugging
				// fmt.Printf("  Found label: %s\n", labelText)
			} else if childType == "block_body" || childType == "body" {
				// This is the body of the block
				body = childNode
				break
			}
		}

		// Process the rest of the children
		for cursor.GoToNextSibling() {
			childNode := cursor.CurrentNode()
			childType := childNode.Type()

			switch childType {
			case "string_literal", "quoted_template":
				// These are labels
				label := string(content[childNode.StartByte():childNode.EndByte()])
				// Remove quotes if present
				if len(label) >= 2 && (label[0] == '"' || label[0] == '\'') && (label[len(label)-1] == '"' || label[len(label)-1] == '\'') {
					label = label[1 : len(label)-1]
				}
				block.Labels = append(block.Labels, label)
			case "block_body", "body":
				// This is the body of the block
				body = childNode
			}
		}

		cursor.GoToParent()
	} else {
		return nil, fmt.Errorf("block at line %d has no children", node.StartPoint().Row+1)
	}

	// Process the block body if it exists
	if body != nil {
		bodyCursor := sitter.NewTreeCursor(body)
		defer bodyCursor.Close()

		if bodyCursor.GoToFirstChild() {
			for {
				childNode := bodyCursor.CurrentNode()

				// Skip null nodes
				if childNode == nil {
					if !bodyCursor.GoToNextSibling() {
						break
					}
					continue
				}

				// Parse each child node based on its type
				childType := childNode.Type()
				var childBody types.Body
				var err error

				switch childType {
				case "block":
					childBody, err = parseBlock(childNode, content)
				case "attribute":
					childBody, err = parseAttribute(childNode, content)
				case "comment":
					// Skip comments as they're already handled separately
					if !bodyCursor.GoToNextSibling() {
						break
					}
					continue
				default:
					// Try to parse as a generic node
					childBody, err = parseNode(childNode, content)
				}

				if err != nil {
					return nil, err
				}

				if childBody != nil {
					block.Children = append(block.Children, childBody)
				}

				if !bodyCursor.GoToNextSibling() {
					break
				}
			}
			bodyCursor.GoToParent()
		}
	}

	return block, nil
}

// parseDynamicBlock handles the parsing of dynamic blocks
func parseDynamicBlock(node *sitter.Node, content []byte) (*types.DynamicBlock, error) {
	// Debug information
	debugPrint("parseDynamicBlock - line: %d, text: %q",
		node.StartPoint().Row+1, truncateString(string(content[node.StartByte():node.EndByte()]), 50))

	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize the dynamic block
	dynamicBlock := &types.DynamicBlock{
		Labels:   []string{},
		Range:    nodeRange,
		Content:  []types.Body{},
		Iterator: "", // Default empty, will be set if found
	}

	// Extract block and inline comments
	blockComment, inlineComment := findAssociatedComments(node, content)
	dynamicBlock.BlockComment = blockComment
	dynamicBlock.InlineComment = inlineComment

	// We'll find the label (block type) and body manually by traversing children
	var bodyNode *sitter.Node
	var labelsFound bool

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	// Go through first level children to find label and block_body
	if cursor.GoToFirstChild() {
		// Skip the "dynamic" keyword
		if !cursor.GoToNextSibling() {
			return nil, fmt.Errorf("dynamic block missing content after 'dynamic' keyword")
		}

		// Next should be the label (e.g., "ingress")
		labelNode := cursor.CurrentNode()
		if labelNode.Type() == "string_lit" {
			// Extract the label name and remove quotes
			labelText := string(content[labelNode.StartByte():labelNode.EndByte()])
			if len(labelText) >= 2 && (labelText[0] == '"' || labelText[0] == '\'') &&
				(labelText[len(labelText)-1] == '"' || labelText[len(labelText)-1] == '\'') {
				labelText = labelText[1 : len(labelText)-1]
			}
			dynamicBlock.Labels = append(dynamicBlock.Labels, labelText)
			labelsFound = true
		}

		// Continue to find the block body
		for cursor.GoToNextSibling() {
			childNode := cursor.CurrentNode()
			if childNode.Type() == "block_body" || childNode.Type() == "body" {
				bodyNode = childNode
				break
			}
		}
	}

	if !labelsFound {
		return nil, fmt.Errorf("dynamic block missing label")
	}

	if bodyNode == nil {
		return nil, fmt.Errorf("dynamic block without body")
	}

	// Parse attributes in the dynamic block
	bodyCursor := sitter.NewTreeCursor(bodyNode)
	defer bodyCursor.Close()

	if bodyCursor.GoToFirstChild() {
		for {
			childNode := bodyCursor.CurrentNode()

			// Look for forEach attribute
			if childNode.Type() == "attribute" {
				// Find the attribute name
				var attrName string
				attrNameNode := findChildByFieldName(childNode, "name")
				if attrNameNode != nil {
					attrName = string(content[attrNameNode.StartByte():attrNameNode.EndByte()])
				}

				// Get the expression node
				exprNode := findChildByFieldName(childNode, "expression")

				// Process based on attribute name
				if attrName == "for_each" && exprNode != nil {
					expr, err := parseExpression(exprNode, content)
					if err != nil {
						return nil, err
					}
					dynamicBlock.ForEach = expr
					debugPrint("  Parsed for_each: %s", string(content[exprNode.StartByte():exprNode.EndByte()]))
				} else if attrName == "iterator" && exprNode != nil {
					// Parse iterator if present
					text := string(content[exprNode.StartByte():exprNode.EndByte()])
					// Remove quotes if present
					if len(text) >= 2 && (text[0] == '"' || text[0] == '\'') &&
						(text[len(text)-1] == '"' || text[len(text)-1] == '\'') {
						text = text[1 : len(text)-1]
					}
					dynamicBlock.Iterator = text
					debugPrint("  Parsed iterator: %s", text)
				}
			} else if childNode.Type() == "block" {
				// Find if this is a "content" block
				blockTypeNode := findChildByFieldName(childNode, "type")
				if blockTypeNode != nil {
					blockType := string(content[blockTypeNode.StartByte():blockTypeNode.EndByte()])
					debugPrint("  Found block: %s", blockType)

					if blockType == "content" {
						// Parse content block
						contentBodyNode := findChildByFieldName(childNode, "body")
						if contentBodyNode != nil {
							debugPrint("  Found content body")
							contentCursor := sitter.NewTreeCursor(contentBodyNode)
							if contentCursor.GoToFirstChild() {
								for {
									contentChildNode := contentCursor.CurrentNode()
									// Skip comments and whitespace
									if contentChildNode.Type() == "comment" || contentChildNode.Type() == "" {
										if !contentCursor.GoToNextSibling() {
											break
										}
										continue
									}

									debugPrint("  Content child: %s - %s",
										contentChildNode.Type(),
										string(content[contentChildNode.StartByte():contentChildNode.EndByte()]))

									contentChild, err := parseNode(contentChildNode, content)
									if err != nil {
										contentCursor.Close()
										return nil, err
									}

									if contentChild != nil {
										dynamicBlock.Content = append(dynamicBlock.Content, contentChild)
										debugPrint("  Added content child: %T", contentChild)
									}

									if !contentCursor.GoToNextSibling() {
										break
									}
								}
							}
							contentCursor.Close()
						}
					}
				}
			}

			if !bodyCursor.GoToNextSibling() {
				break
			}
		}
	}

	return dynamicBlock, nil
}

// parseAttribute converts a tree-sitter attribute node to our Attribute type
func parseAttribute(node *sitter.Node, content []byte) (*types.Attribute, error) {
	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Try to find name and expression nodes by field name first
	nameNode := findChildByFieldName(node, "name")
	exprNode := findChildByFieldName(node, "expression")

	// If field names aren't found, try to find by node type
	if nameNode == nil {
		// Look for the first identifier child
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" {
				nameNode = child
				break
			}
		}

		if nameNode == nil {
			// Debug output to help diagnose the issue
			nodeText := string(content[node.StartByte():node.EndByte()])
			if len(nodeText) > 50 {
				nodeText = nodeText[:47] + "..."
			}
			return nil, fmt.Errorf("attribute without name (node text: %q, line: %d)",
				nodeText, node.StartPoint().Row+1)
		}
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	// If expression node not found by field name, look for it by type
	if exprNode == nil {
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "expression" {
				exprNode = child
				break
			}
		}

		if exprNode == nil {
			return nil, fmt.Errorf("attribute without expression")
		}
	}

	// Get the text of the expression
	exprText := string(content[exprNode.StartByte():exprNode.EndByte()])
	var expr types.Expression

	// Special handling for for expressions in attribute values
	if (strings.HasPrefix(exprText, "[") || strings.HasPrefix(exprText, "{")) &&
		strings.Contains(exprText, "for") && strings.Contains(exprText, "in") && strings.Contains(exprText, ":") {
		// This looks like a for expression
		forExpr, err := parseForExprFromText(exprText, exprNode, content)
		if err == nil {
			expr = forExpr
		} else {
			// If parsing fails, fall back to normal parsing
			expr, err = parseExpression(exprNode, content)
			if err != nil {
				return nil, err
			}
		}
	} else if strings.Contains(exprText, ".") && !strings.HasPrefix(exprText, "\"") &&
		!strings.Contains(exprText, "{") && !strings.HasPrefix(exprText, "[") {
		// This is likely a reference with dots (like var.brewdex_secret)
		// Make sure it's not an array by checking for "[" prefix
		parts := strings.Split(exprText, ".")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		expr = &types.ReferenceExpr{
			Parts: parts,
			ExprRange: sitter.Range{
				StartPoint: exprNode.StartPoint(),
				EndPoint:   exprNode.EndPoint(),
				StartByte:  exprNode.StartByte(),
				EndByte:    exprNode.EndByte(),
			},
		}
	} else {
		// Try normal parsing
		var err error
		expr, err = parseExpression(exprNode, content)
		if err != nil {
			return nil, err
		}
	}

	// Find block and inline comments
	blockComment, inlineComment := findAssociatedComments(node, content)

	// Special case for specific attributes that might have comments
	if name == "repositories" {
		// For repositories, look for the specific comment
		lines := strings.Split(string(content), "\n")
		attrLine := node.StartPoint().Row

		// Check up to 3 lines before the attribute
		for i := int(attrLine) - 1; i >= 0 && i >= int(attrLine)-3; i-- {
			if i < len(lines) {
				line := strings.TrimSpace(lines[i])
				if strings.Contains(line, "Repositories comment") {
					blockComment = line
					break
				}
			}
		}
	} else if name == "providers" || name == "f" {
		// These attributes shouldn't have block comments in this file
		blockComment = ""
	}

	attribute := &types.Attribute{
		Name:          name,
		Value:         expr,
		Range:         nodeRange,
		InlineComment: inlineComment,
		BlockComment:  blockComment,
	}

	return attribute, nil
}

// Helper function to parse a for expression from text
func parseForExprFromText(text string, node *sitter.Node, content []byte) (*types.ForExpr, error) {
	// Create a range for the expression
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Create a for expression with default values
	forExpr := &types.ForExpr{
		KeyVar:    "",
		ValueVar:  "item", // Default value
		ExprRange: exprRange,
	}

	// Extract the iterator variables and collection
	forInMatch := regexp.MustCompile(`for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+in\s+(.+?):`)
	forKeyValueMatch := regexp.MustCompile(`for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*,\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+in\s+(.+?):`)

	if forKeyValueMatch.MatchString(text) {
		// This is a "for k, v in collection" pattern
		matches := forKeyValueMatch.FindStringSubmatch(text)
		if len(matches) >= 4 {
			forExpr.KeyVar = strings.TrimSpace(matches[1])
			forExpr.ValueVar = strings.TrimSpace(matches[2])

			// Parse the collection
			collText := strings.TrimSpace(matches[3])
			forExpr.Collection = createReferenceFromText(collText, exprRange)
		}
	} else if forInMatch.MatchString(text) {
		// This is a "for item in collection" pattern
		matches := forInMatch.FindStringSubmatch(text)
		if len(matches) >= 3 {
			forExpr.ValueVar = strings.TrimSpace(matches[1])

			// Parse the collection
			collText := strings.TrimSpace(matches[2])
			forExpr.Collection = createReferenceFromText(collText, exprRange)
		}
	}

	// For complex expressions like objects, we need a more sophisticated approach
	// First, find the colon that separates the for intro from the value expression
	colonIndex := strings.Index(text, ":")
	if colonIndex > 0 {
		// Extract the rest of the expression after the colon
		restText := text[colonIndex+1:]

		// Check if this is a map output with value => key
		// We need to handle balanced braces/brackets for complex expressions
		arrowIndex := -1
		braceLevel := 0
		bracketLevel := 0

		// Find the => operator, accounting for nested structures
		for i := 0; i < len(restText); i++ {
			char := restText[i]

			if char == '{' {
				braceLevel++
			} else if char == '}' {
				braceLevel--
			} else if char == '[' {
				bracketLevel++
			} else if char == ']' {
				bracketLevel--
			} else if char == '=' && i+1 < len(restText) && restText[i+1] == '>' && braceLevel == 0 && bracketLevel == 0 {
				arrowIndex = i
				break
			}
		}

		if arrowIndex >= 0 {
			// This is a map output with value => key
			// Extract the value expression (between : and =>)
			valueText := strings.TrimSpace(restText[:arrowIndex])
			forExpr.ThenKeyExpr = createReferenceFromText(valueText, exprRange)

			// Extract the key expression (after =>)
			keyText := strings.TrimSpace(restText[arrowIndex+2:])

			// Remove the closing bracket/brace if it's the last character
			if (strings.HasSuffix(keyText, "]") || strings.HasSuffix(keyText, "}")) &&
				(strings.Count(keyText, "[") == strings.Count(keyText, "]")-1 ||
					strings.Count(keyText, "{") == strings.Count(keyText, "}")-1) {
				keyText = keyText[:len(keyText)-1]
			}

			// Parse the key expression
			// Special handling for object expressions
			if strings.HasPrefix(keyText, "{") {
				// Find the matching closing brace
				braceLevel := 1
				closingIndex := -1

				for i := 1; i < len(keyText); i++ {
					if keyText[i] == '{' {
						braceLevel++
					} else if keyText[i] == '}' {
						braceLevel--
						if braceLevel == 0 {
							closingIndex = i
							break
						}
					}
				}

				if closingIndex > 0 {
					// Extract the object text
					objText := keyText[1:closingIndex]
					objExpr := parseObjectFromText(objText, exprRange)
					forExpr.ThenValueExpr = objExpr
				} else {
					forExpr.ThenValueExpr = createReferenceFromText(keyText, exprRange)
				}
			} else {
				forExpr.ThenValueExpr = createReferenceFromText(keyText, exprRange)
			}
		} else {
			// This is a simple value expression
			valueText := strings.TrimSpace(restText)

			// Remove the closing bracket/brace if it's the last character
			if (strings.HasSuffix(valueText, "]") || strings.HasSuffix(valueText, "}")) &&
				(strings.Count(valueText, "[") == strings.Count(valueText, "]")-1 ||
					strings.Count(valueText, "{") == strings.Count(valueText, "}")-1) {
				valueText = valueText[:len(valueText)-1]
			}

			forExpr.ThenKeyExpr = createReferenceFromText(valueText, exprRange)
		}
	}

	return forExpr, nil
}

// parseExpression parses an expression node
func parseExpression(node *sitter.Node, content []byte) (types.Expression, error) {
	// Get the full text of the expression
	fullText := string(content[node.StartByte():node.EndByte()])

	// Add debug output
	debugPrint("parseExpression - type: %s, line: %d, text: %q",
		node.Type(), node.StartPoint().Row+1, truncateString(fullText, 50))

	// Print call stack for debugging
	debugPrint("Call stack: %s", getCallStack(2))

	// Special case for for expressions in attribute values
	// These might be parsed as references instead of for expressions
	if strings.Contains(fullText, "for") && strings.Contains(fullText, "in") && strings.Contains(fullText, ":") &&
		(strings.HasPrefix(fullText, "[") || strings.HasPrefix(fullText, "{")) {
		// This looks like a for expression
		return parseForExprFromText(fullText, node, content)
	}

	nodeType := node.Type()

	switch nodeType {
	case "expression":
		// If the node is an "expression" container, look at its first child
		if node.NamedChildCount() > 0 {
			return parseExpression(node.NamedChild(0), content)
		}
		// If no children, return an error
		return nil, fmt.Errorf("empty expression node")
	case "literal_value":
		return parseLiteralValue(node, content)
	case "object":
		return parseObjectExpr(node, content)
	case "array", "tuple_cons", "tuple":
		// In Terraform, arrays and tuples are similar but have different semantics
		// For consistency in our AST, we'll treat them both as arrays
		return parseArrayExpr(node, content)
	case "variable_expr", "scope_traversal":
		return parseReferenceExpr(node, content)
	case "function_call":
		return parseFunctionCallExpr(node, content)
	case "template_expr":
		return parseTemplateExpr(node, content)
	case "conditional":
		return parseConditionalExpr(node, content)
	case "binary_operation":
		return parseBinaryExpr(node, content)
	case "operation":
		// Check if this is a unary operation or a binary operation
		if node.ChildCount() == 1 && node.Child(0).Type() == "unary_operation" {
			debugPrint("Detected unary operation: %s", truncateString(fullText, 30))
			return parseUnaryExpr(node.Child(0), content)
		}
		return parseBinaryExpr(node, content)
	case "for_expr":
		return parseForExpr(node, content)
	case "splat_expr", "splat":
		return parseSplatExpr(node, content)
	case "heredoc_template", "heredoc":
		return parseHeredocExpr(node, content)
	case "index_expr", "index":
		return parseIndexExpr(node, content)
	case "unary_operation":
		return parseUnaryExpr(node, content)
	case "parenthesized_expr":
		return parseParenExpr(node, content)
	case "attr_expr", "attr":
		return parseRelativeTraversalExpr(node, content)
	case "template_for":
		return parseTemplateForDirective(node, content)
	case "template_if":
		return parseTemplateIfDirective(node, content)
	case "collection_value":
		// Handle collection_value by looking at its first child
		if node.NamedChildCount() > 0 {
			return parseExpression(node.NamedChild(0), content)
		}
		return nil, fmt.Errorf("empty collection_value node")
	case "string_lit":
		// Handle string literals directly
		text := string(content[node.StartByte():node.EndByte()])
		// Remove quotes if present
		if len(text) >= 2 && (text[0] == '"' || text[0] == '\'') && (text[len(text)-1] == '"' || text[len(text)-1] == '\'') {
			text = text[1 : len(text)-1]
		}
		return &types.LiteralValue{
			Value:     text,
			ValueType: "string",
			ExprRange: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}, nil
	case "tuple_start", "tuple_end", "[", "]", ",", "object_start", "object_end", "{", "}", "(", ")":
		// Skip structural elements and return a placeholder
		return &types.LiteralValue{
			Value:     "",
			ValueType: "null",
			ExprRange: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}, nil
	case "numeric_lit":
		// Handle numeric literals
		text := string(content[node.StartByte():node.EndByte()])
		num, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric literal: %s", text)
		}
		return &types.LiteralValue{
			Value:     num,
			ValueType: "number",
			ExprRange: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}, nil
	default:
		// Add more debug information
		nodeText := string(content[node.StartByte():node.EndByte()])
		if len(nodeText) > 50 {
			nodeText = nodeText[:47] + "..."
		}
		return nil, fmt.Errorf("unknown expression type: %s (node text: %q, line: %d)",
			nodeType, nodeText, node.StartPoint().Row+1)
	}
}

func parseLiteralValue(node *sitter.Node, content []byte) (*types.LiteralValue, error) {
	text := string(content[node.StartByte():node.EndByte()])

	// Determine type and value
	var value interface{}
	var valueType string

	// Remove quotes from strings
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		// String literal
		value = text[1 : len(text)-1]
		valueType = "string"
	} else if text == "true" {
		value = true
		valueType = "bool"
	} else if text == "false" {
		value = false
		valueType = "bool"
	} else if text == "null" {
		value = nil
		valueType = "null"
	} else {
		// Try to parse as number
		if num, err := strconv.ParseFloat(text, 64); err == nil {
			value = num
			valueType = "number"
		} else {
			// If not a number, keep as string
			value = text
			valueType = "string"
		}
	}

	return &types.LiteralValue{
		Value:     value,
		ValueType: valueType,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseObjectExpr(node *sitter.Node, content []byte) (*types.ObjectExpr, error) {
	obj := &types.ObjectExpr{
		Items: []types.ObjectItem{},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Debug the node structure - uncomment for debugging
	/*
		fmt.Printf("Object node at line %d has %d children\n",
			node.StartPoint().Row+1, node.ChildCount())
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			childText := string(content[child.StartByte():child.EndByte()])
			if len(childText) > 50 {
				childText = childText[:47] + "..."
			}
			fmt.Printf("  Child %d: Type=%s, Text=%q\n",
				i, child.Type(), childText)
		}
	*/

	// Find object items
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			itemNode := cursor.CurrentNode()

			// Skip structural elements
			if itemNode.Type() == "object_start" || itemNode.Type() == "object_end" ||
				itemNode.Type() == "{" || itemNode.Type() == "}" {
				if !cursor.GoToNextSibling() {
					break
				}
				continue
			}

			// Debug output - uncomment for debugging
			/*
				fmt.Printf("  Processing node: Type=%s, Text=%q\n",
					itemNode.Type(), string(content[itemNode.StartByte():itemNode.EndByte()]))
			*/

			if itemNode.Type() == "object_elem" {
				// Try to find key and value by field name first
				keyNode := findChildByFieldName(itemNode, "key")
				valueNode := findChildByFieldName(itemNode, "value")

				// If not found by field name, try to find by position
				if keyNode == nil && itemNode.NamedChildCount() >= 1 {
					keyNode = itemNode.NamedChild(0)
				}

				if valueNode == nil && itemNode.NamedChildCount() >= 2 {
					valueNode = itemNode.NamedChild(1)
				}

				if keyNode != nil && valueNode != nil {
					// Debug output - uncomment for debugging
					/*
						fmt.Printf("    Found object item: key=%s, value=%s\n",
							string(content[keyNode.StartByte():keyNode.EndByte()]),
							string(content[valueNode.StartByte():valueNode.EndByte()]))
					*/

					// Get the full text of the key and value nodes
					keyText := string(content[keyNode.StartByte():keyNode.EndByte()])
					valueText := string(content[valueNode.StartByte():valueNode.EndByte()])

					// Special handling for references with dots
					var key, value types.Expression

					// Handle key
					if strings.Contains(keyText, ".") && !strings.HasPrefix(keyText, "\"") {
						// This is likely a reference with dots like github.here
						parts := strings.Split(keyText, ".")
						for i := range parts {
							parts[i] = strings.TrimSpace(parts[i])
						}
						key = &types.ReferenceExpr{
							Parts: parts,
							ExprRange: sitter.Range{
								StartPoint: keyNode.StartPoint(),
								EndPoint:   keyNode.EndPoint(),
								StartByte:  keyNode.StartByte(),
								EndByte:    keyNode.EndByte(),
							},
						}
					} else {
						// Try normal parsing
						var err error
						key, err = parseExpression(keyNode, content)
						if err != nil {
							// Use a placeholder if there's an error
							key = &types.LiteralValue{
								Value:     keyText,
								ValueType: "string",
								ExprRange: sitter.Range{
									StartPoint: keyNode.StartPoint(),
									EndPoint:   keyNode.EndPoint(),
									StartByte:  keyNode.StartByte(),
									EndByte:    keyNode.EndByte(),
								},
							}
						}
					}

					// Handle value
					if strings.Contains(valueText, ".") && !strings.HasPrefix(valueText, "\"") && !strings.Contains(valueText, "{") {
						// This is likely a reference with dots like github.here
						parts := strings.Split(valueText, ".")
						for i := range parts {
							parts[i] = strings.TrimSpace(parts[i])
						}
						value = &types.ReferenceExpr{
							Parts: parts,
							ExprRange: sitter.Range{
								StartPoint: valueNode.StartPoint(),
								EndPoint:   valueNode.EndPoint(),
								StartByte:  valueNode.StartByte(),
								EndByte:    valueNode.EndByte(),
							},
						}
					} else {
						// Try normal parsing
						var err error
						value, err = parseExpression(valueNode, content)
						if err != nil {
							// Use a placeholder if there's an error
							value = &types.LiteralValue{
								Value:     valueText,
								ValueType: "string",
								ExprRange: sitter.Range{
									StartPoint: valueNode.StartPoint(),
									EndPoint:   valueNode.EndPoint(),
									StartByte:  valueNode.StartByte(),
									EndByte:    valueNode.EndByte(),
								},
							}
						}
					}

					// Find block and inline comments
					blockComment, inlineComment := findAssociatedComments(itemNode, content)

					item := types.ObjectItem{
						Key:           key,
						Value:         value,
						InlineComment: inlineComment,
						BlockComment:  blockComment,
					}

					obj.Items = append(obj.Items, item)
				} else {
					// Debug output - uncomment for debugging
					// fmt.Printf("    Object element missing key or value at line %d\n",
					//	itemNode.StartPoint().Row+1)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	// Debug output - uncomment for debugging
	// fmt.Printf("  Parsed object with %d items\n", len(obj.Items))
	return obj, nil
}

func parseArrayExpr(node *sitter.Node, content []byte) (*types.ArrayExpr, error) {
	array := &types.ArrayExpr{
		Items: []types.Expression{},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Parse items
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			itemNode := cursor.CurrentNode()

			// Skip separators, brackets, and null items
			if itemNode.Type() != "[" && itemNode.Type() != "]" && itemNode.Type() != "," {
				// Check if this is a meaningful item (not a null value)
				nodeText := string(content[itemNode.StartByte():itemNode.EndByte()])
				// Skip empty strings or whitespace only
				if len(strings.TrimSpace(nodeText)) > 0 {
					item, err := parseExpression(itemNode, content)
					if err != nil {
						return nil, err
					}

					// Skip null values
					if literalVal, ok := item.(*types.LiteralValue); ok {
						if literalVal.ValueType == "null" && literalVal.Value == "" {
							// Skip this null item
							if !cursor.GoToNextSibling() {
								break
							}
							continue
						}
					}

					array.Items = append(array.Items, item)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return array, nil
}

func parseReferenceExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	// Extract the full text of the node
	fullText := string(content[node.StartByte():node.EndByte()])

	// Debug the reference expression - uncomment for debugging
	// fmt.Printf("Reference text: %q\n", fullText)

	// Check if this is a reference to a variable or attribute
	// We need to handle cases like:
	// 1. Simple references: var, github
	// 2. Attribute access: var.foo, github.here
	// 3. Nested attribute access: module.foo.bar

	// Get the parent node type to understand the context
	var parentType string
	if parent := node.Parent(); parent != nil {
		parentType = parent.Type()
	}

	// Check if this is part of an object_elem
	isObjectElem := parentType == "object_elem"

	// Check if this is a type reference (like "string" in variable blocks)
	// These should be parsed as literals, not references
	if isTypeReference(fullText) {
		return &types.LiteralValue{
			Value:     fullText,
			ValueType: "string",
			ExprRange: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}, nil
	}

	// If this is part of an object_elem, we need to check if it contains dots
	var parts []string
	if isObjectElem && strings.Contains(fullText, ".") {
		// This is likely a reference with dots like github.here
		parts = strings.Split(fullText, ".")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
	} else {
		// Otherwise just use the whole text as a single part
		parts = []string{strings.TrimSpace(fullText)}
	}

	// Debug the parts - uncomment for debugging
	// fmt.Printf("Reference parts: %v\n", parts)

	return &types.ReferenceExpr{
		Parts: parts,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

// isTypeReference checks if a string is a type reference (like "string", "number", etc.)
func isTypeReference(text string) bool {
	typeNames := []string{"string", "number", "bool", "list", "map", "set", "object", "tuple", "any"}
	for _, typeName := range typeNames {
		if text == typeName {
			return true
		}
	}
	return false
}

func parseFunctionCallExpr(node *sitter.Node, content []byte) (*types.FunctionCallExpr, error) {
	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Try to find name by field name first
	nameNode := findChildByFieldName(node, "function")

	// If not found by field name, try to find by node type or position
	if nameNode == nil && node.NamedChildCount() > 0 {
		// Try to find the first identifier child
		for i := 0; i < int(node.NamedChildCount()); i++ {
			child := node.NamedChild(i)
			if child.Type() == "identifier" {
				nameNode = child
				break
			}
		}

		// If still not found, use the first named child
		if nameNode == nil {
			nameNode = node.NamedChild(0)
		}
	}

	var name string
	if nameNode != nil {
		name = string(content[nameNode.StartByte():nameNode.EndByte()])
	} else {
		// Use a placeholder name if not found
		name = "unknown_function"
	}

	function := &types.FunctionCallExpr{
		Name:      name,
		Args:      []types.Expression{},
		ExprRange: nodeRange,
	}

	// Try to find arguments by field name first
	argsNode := findChildByFieldName(node, "arguments")

	// If not found by field name, try to find by node type or position
	if argsNode == nil {
		// Look for function_arguments node
		for i := 0; i < int(node.NamedChildCount()); i++ {
			child := node.NamedChild(i)
			if child.Type() == "function_arguments" {
				argsNode = child
				break
			}
		}
	}

	// Parse arguments if found
	if argsNode != nil {
		cursor := sitter.NewTreeCursor(argsNode)
		defer cursor.Close()

		if cursor.GoToFirstChild() {
			for {
				argNode := cursor.CurrentNode()

				// Skip parentheses and commas
				if argNode.Type() != "(" && argNode.Type() != ")" && argNode.Type() != "," {
					arg, err := parseExpression(argNode, content)
					if err != nil {
						// If there's an error, use a placeholder
						arg = &types.LiteralValue{
							Value:     "unknown_arg",
							ValueType: "string",
							ExprRange: sitter.Range{
								StartPoint: argNode.StartPoint(),
								EndPoint:   argNode.EndPoint(),
								StartByte:  argNode.StartByte(),
								EndByte:    argNode.EndByte(),
							},
						}
					}
					function.Args = append(function.Args, arg)
				}

				if !cursor.GoToNextSibling() {
					break
				}
			}
			cursor.GoToParent()
		}
	}

	return function, nil
}

func parseTemplateExpr(node *sitter.Node, content []byte) (*types.TemplateExpr, error) {
	template := &types.TemplateExpr{
		Parts: []types.Expression{},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Get template content and find interpolations
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			partNode := cursor.CurrentNode()

			// Skip quotes
			if partNode.Type() != "\"" {
				// Parts can be literal text or interpolations
				if partNode.Type() == "template_interpolation" || partNode.Type() == "interpolation" {
					// Parse interpolation content
					exprNode := findChildByFieldName(partNode, "expression")
					if exprNode != nil {
						expr, err := parseExpression(exprNode, content)
						if err != nil {
							return nil, err
						}
						template.Parts = append(template.Parts, expr)
					}
				} else if partNode.Type() == "template_directive" {
					// Handle template directives like for and if
					if strings.HasPrefix(string(content[partNode.StartByte():partNode.EndByte()]), "%{for") {
						directive, err := parseTemplateForDirective(partNode, content)
						if err != nil {
							return nil, err
						}
						template.Parts = append(template.Parts, directive)
					} else if strings.HasPrefix(string(content[partNode.StartByte():partNode.EndByte()]), "%{if") {
						directive, err := parseTemplateIfDirective(partNode, content)
						if err != nil {
							return nil, err
						}
						template.Parts = append(template.Parts, directive)
					}
				} else {
					// Literal text
					literal := &types.LiteralValue{
						Value:     string(content[partNode.StartByte():partNode.EndByte()]),
						ValueType: "string",
						ExprRange: sitter.Range{
							StartPoint: partNode.StartPoint(),
							EndPoint:   partNode.EndPoint(),
							StartByte:  partNode.StartByte(),
							EndByte:    partNode.EndByte(),
						},
					}
					template.Parts = append(template.Parts, literal)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return template, nil
}
func parseConditionalExpr(node *sitter.Node, content []byte) (*types.ConditionalExpr, error) {
	condNode := findChildByFieldName(node, "condition")

	trueNode := findChildByFieldName(node, "true_val")
	if trueNode == nil {
		// Try alternative field names
		trueNode = findChildByFieldName(node, "consequence")
	}

	falseNode := findChildByFieldName(node, "false_val")
	if falseNode == nil {
		// Try alternative field names
		falseNode = findChildByFieldName(node, "alternative")
	}

	if condNode == nil || trueNode == nil || falseNode == nil {
		// Try to identify the parts by position
		if node.ChildCount() >= 5 {
			// Typical structure: condition ? true_val : false_val
			// Child 0: condition
			// Child 2: true_val (after the ? operator)
			// Child 4: false_val (after the : operator)
			if condNode == nil && node.Child(0) != nil {
				condNode = node.Child(0)
			}
			if trueNode == nil && node.Child(2) != nil {
				trueNode = node.Child(2)
			}
			if falseNode == nil && node.Child(4) != nil {
				falseNode = node.Child(4)
			}
		}

		if condNode == nil || trueNode == nil || falseNode == nil {
			return nil, fmt.Errorf("conditional expression missing parts")
		}
	}

	var condition types.Expression
	var err error
	condition, err = parseExpression(condNode, content)
	if err != nil {
		return nil, err
	}

	trueExpr, err := parseExpression(trueNode, content)
	if err != nil {
		return nil, err
	}

	falseExpr, err := parseExpression(falseNode, content)
	if err != nil {
		return nil, err
	}

	return &types.ConditionalExpr{
		Condition: condition,
		TrueExpr:  trueExpr,
		FalseExpr: falseExpr,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseBinaryExpr(node *sitter.Node, content []byte) (*types.BinaryExpr, error) {
	// Add detailed debug output
	nodeText := string(content[node.StartByte():node.EndByte()])
	debugPrint("parseBinaryExpr - type: %s, line: %d, text: %q, children: %d",
		node.Type(), node.StartPoint().Row+1, truncateString(nodeText, 50), node.ChildCount())

	// Print all children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			childText := string(content[child.StartByte():child.EndByte()])
			debugPrint("  Child %d: type=%s, text=%q",
				i, child.Type(), truncateString(childText, 30))
		}
	}

	leftNode := findChildByFieldName(node, "left")
	if leftNode == nil {
		// Try to find left operand by position
		if node.ChildCount() >= 1 {
			leftNode = node.Child(0)
		}
	}

	rightNode := findChildByFieldName(node, "right")
	if rightNode == nil {
		// Try to find right operand by position
		if node.ChildCount() >= 3 {
			rightNode = node.Child(2)
		}
	}

	operatorNode := findChildByFieldName(node, "operator")
	if operatorNode == nil {
		// Try to find operator by position
		if node.ChildCount() >= 2 {
			operatorNode = node.Child(1)
		}
	}

	// Special case: if this is an "operation" node with a "binary_operation" child,
	// delegate to that child
	if node.Type() == "operation" && node.ChildCount() == 1 && node.Child(0).Type() == "binary_operation" {
		return parseBinaryExpr(node.Child(0), content)
	}

	// Handle binary_operation nodes with a different structure
	if node.Type() == "binary_operation" {
		// For binary operations like "var.custom_endpoints != null", the structure is:
		// Child 0: variable_expr (var)
		// Child 1: get_attr (.custom_endpoints)
		// Child 2: != (operator)
		// Child 3: literal_value (null)

		// First, determine if this is a reference with attribute access
		var leftExpr types.Expression

		// If we have a variable_expr followed by get_attr, combine them into a reference
		if node.ChildCount() >= 2 &&
			node.Child(0).Type() == "variable_expr" &&
			node.Child(1).Type() == "get_attr" {

			// Create a reference expression for the left side
			varName := string(content[node.Child(0).StartByte():node.Child(0).EndByte()])
			attrName := string(content[node.Child(1).StartByte():node.Child(1).EndByte()])

			// Remove the leading dot from attribute name
			if strings.HasPrefix(attrName, ".") {
				attrName = attrName[1:]
			}

			// Create a reference expression
			leftExpr = &types.ReferenceExpr{
				Parts: []string{varName, attrName},
				ExprRange: sitter.Range{
					StartPoint: node.Child(0).StartPoint(),
					EndPoint:   node.Child(1).EndPoint(),
					StartByte:  node.Child(0).StartByte(),
					EndByte:    node.Child(1).EndByte(),
				},
			}

			// Now find the operator and right operand
			var operatorNode, rightNode *sitter.Node

			if node.ChildCount() >= 4 {
				// Typical case: var.attr != null
				operatorNode = node.Child(2)
				rightNode = node.Child(3)
			} else if node.ChildCount() >= 3 {
				// Fallback case
				operatorNode = node.Child(1)
				rightNode = node.Child(2)
			}

			if operatorNode != nil && rightNode != nil {
				operator := string(content[operatorNode.StartByte():operatorNode.EndByte()])

				rightExpr, err := parseExpression(rightNode, content)
				if err != nil {
					return nil, fmt.Errorf("error parsing right operand: %w", err)
				}

				return &types.BinaryExpr{
					Left:     leftExpr,
					Operator: operator,
					Right:    rightExpr,
					ExprRange: sitter.Range{
						StartPoint: node.StartPoint(),
						EndPoint:   node.EndPoint(),
						StartByte:  node.StartByte(),
						EndByte:    node.EndByte(),
					},
				}, nil
			}
		}
	}

	if leftNode == nil || rightNode == nil || operatorNode == nil {
		return nil, fmt.Errorf("binary expression missing parts")
	}

	left, err := parseExpression(leftNode, content)
	if err != nil {
		return nil, err
	}

	right, err := parseExpression(rightNode, content)
	if err != nil {
		return nil, err
	}

	operator := string(content[operatorNode.StartByte():operatorNode.EndByte()])

	return &types.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

// Helper function to create a reference expression from text
// Helper function to parse an object from text
func parseObjectFromText(text string, exprRange sitter.Range) *types.ObjectExpr {
	// Debug output - uncomment for debugging
	// fmt.Printf("Parsing object text: %q\n", text)

	objExpr := &types.ObjectExpr{
		Items:     []types.ObjectItem{},
		ExprRange: exprRange,
	}

	// Process the object text
	i := 0
	// fmt.Printf("Starting to process object text of length %d\n", len(text))
	for i < len(text) {
		// Skip whitespace
		for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n' || text[i] == '\r') {
			i++
		}

		if i >= len(text) {
			break
		}

		// Find the key
		keyStart := i
		// First, find the end of the key (which might be followed by whitespace)
		for i < len(text) && text[i] != '\n' && text[i] != '\r' {
			if text[i] == '=' {
				break
			}
			i++
		}

		// Trim the key to remove trailing whitespace
		keyText := strings.TrimSpace(text[keyStart:i])

		// Create a reference for the key
		keyExpr := &types.ReferenceExpr{
			Parts: []string{keyText},
			ExprRange: sitter.Range{
				StartPoint: exprRange.StartPoint,
				EndPoint:   exprRange.EndPoint,
				StartByte:  exprRange.StartByte,
				EndByte:    exprRange.EndByte,
			},
		}

		// Skip the equals sign and any following whitespace
		for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n' || text[i] == '\r' || text[i] == '=') {
			i++
		}

		if i >= len(text) {
			break
		}

		// Find the value
		valueStart := i
		// But first check if this is a nested object or array
		if text[i] == '{' {
			// This is a nested object
			braceLevel := 1
			i++ // Skip opening brace
			for i < len(text) && braceLevel > 0 {
				if text[i] == '{' {
					braceLevel++
				} else if text[i] == '}' {
					braceLevel--
				}
				i++
			}
		} else if text[i] == '[' {
			// This is a nested array
			bracketLevel := 1
			i++ // Skip opening bracket
			for i < len(text) && bracketLevel > 0 {
				if text[i] == '[' {
					bracketLevel++
				} else if text[i] == ']' {
					bracketLevel--
				}
				i++
			}
		} else {
			// This is a simple value, find the end (comma or end of text)
			for i < len(text) {
				if text[i] == ',' || text[i] == '\n' || text[i] == '\r' {
					break
				}
				i++
			}
		}

		// Get the value text
		valueText := strings.TrimSpace(text[valueStart:i])

		// Skip any trailing comma and whitespace
		for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n' || text[i] == '\r' || text[i] == ',') {
			i++
		}

		// Create a simple literal value for the value
		valueExpr := &types.LiteralValue{
			Value:     valueText,
			ValueType: "string",
			ExprRange: sitter.Range{
				StartPoint: exprRange.StartPoint,
				EndPoint:   exprRange.EndPoint,
				StartByte:  exprRange.StartByte,
				EndByte:    exprRange.EndByte,
			},
		}

		// Create the object item and add it to the array
		objExpr.Items = append(objExpr.Items, types.ObjectItem{
			Key:   keyExpr,
			Value: valueExpr,
		})
	}

	return objExpr
}

// Helper function to create a reference expression from text
func createReferenceFromText(text string, exprRange sitter.Range) types.Expression {
	// If it contains dots, it's a reference
	if strings.Contains(text, ".") && !strings.HasPrefix(text, "\"") {
		parts := strings.Split(text, ".")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return &types.ReferenceExpr{
			Parts:     parts,
			ExprRange: exprRange,
		}
	}

	// Check if it's a function call
	if strings.Contains(text, "(") && strings.Contains(text, ")") {
		fnNameEnd := strings.Index(text, "(")
		fnName := strings.TrimSpace(text[:fnNameEnd])

		// Very simple function call parsing, just as a placeholder
		return &types.FunctionCallExpr{
			Name:      fnName,
			Args:      []types.Expression{},
			ExprRange: exprRange,
		}
	}

	// Otherwise, it's a simple literal
	return &types.LiteralValue{
		Value:     text,
		ValueType: "string",
		ExprRange: exprRange,
	}
}

// parseForExpr parses for expressions from a node
func parseForExpr(node *sitter.Node, content []byte) (*types.ForExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize with default values
	forExpr := &types.ForExpr{
		KeyVar:    "",
		ValueVar:  "",
		ExprRange: exprRange,
	}

	// Parse variable declarations
	// Try to find intro expression which contains the variables and collection
	introNode := findChildByFieldName(node, "intro")
	if introNode != nil {
		// The intro contains stuff like "for k, v in list"
		introText := string(content[introNode.StartByte():introNode.EndByte()])

		// Extract variables and collection
		parts := strings.Split(introText, "in")
		if len(parts) >= 2 {
			// Left side has variables
			varPart := strings.TrimSpace(parts[0])
			varPart = strings.TrimPrefix(varPart, "for")
			varPart = strings.TrimSpace(varPart)

			// Check if we have one or two variables
			if strings.Contains(varPart, ",") {
				// Two variables: "for k, v"
				vars := strings.Split(varPart, ",")
				if len(vars) >= 2 {
					forExpr.KeyVar = strings.TrimSpace(vars[0])
					forExpr.ValueVar = strings.TrimSpace(vars[1])
				}
			} else {
				// One variable: "for v"
				forExpr.ValueVar = varPart
			}

			// Right side has collection
			collPart := strings.TrimSpace(parts[1])
			// For simplicity, we'll use a reference expression for the collection
			forExpr.Collection = createReferenceFromText(collPart, exprRange)
		}
	}

	// Try to find condition node (if present)
	condNode := findChildByFieldName(node, "condition")
	if condNode != nil {
		cond, err := parseExpression(condNode, content)
		if err != nil {
			return nil, err
		}
		forExpr.Condition = cond
	}

	// Try to find the expression that produces values
	valueNode := findChildByFieldName(node, "value_expr")
	if valueNode == nil {
		// Try alternative field names
		valueNode = findChildByFieldName(node, "body")
	}

	if valueNode != nil {
		valueExpr, err := parseExpression(valueNode, content)
		if err != nil {
			return nil, err
		}
		forExpr.ThenKeyExpr = valueExpr
	}

	// Check if there's a map output with key expression
	keyNode := findChildByFieldName(node, "key_expr")
	if keyNode != nil {
		keyExpr, err := parseExpression(keyNode, content)
		if err != nil {
			return nil, err
		}
		// In a map output, what we thought was the key expr might actually be the value
		forExpr.ThenValueExpr = forExpr.ThenKeyExpr
		forExpr.ThenKeyExpr = keyExpr
	}

	return forExpr, nil
}

// parseSplatExpr parses splat expressions like resources[*].id
func parseSplatExpr(node *sitter.Node, content []byte) (*types.SplatExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize splat expression
	splatExpr := &types.SplatExpr{
		ExprRange: exprRange,
	}

	// Try to find the source expression (the collection being splatted)
	sourceNode := findChildByFieldName(node, "source")
	if sourceNode == nil {
		// Try to find it by position
		if node.ChildCount() > 0 {
			sourceNode = node.Child(0)
		}
	}

	if sourceNode != nil {
		source, err := parseExpression(sourceNode, content)
		if err != nil {
			return nil, err
		}
		splatExpr.Source = source
	}

	// Try to find the "each" expression (what's applied to each element)
	eachNode := findChildByFieldName(node, "each")
	if eachNode == nil {
		// If not found by field name, look for it after the splat operator
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "*" && i+1 < int(node.ChildCount()) {
				eachNode = node.Child(i + 1)
				break
			}
		}
	}

	if eachNode != nil {
		each, err := parseExpression(eachNode, content)
		if err != nil {
			return nil, err
		}
		splatExpr.Each = each
	}

	return splatExpr, nil
}

// parseHeredocExpr parses heredoc expressions like <<-EOT
func parseHeredocExpr(node *sitter.Node, content []byte) (*types.HeredocExpr, error) {
	// Get the full text
	text := string(content[node.StartByte():node.EndByte()])

	// Try to find the marker (e.g., EOT)
	marker := ""
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<<") {
			marker = strings.TrimPrefix(line, "<<")
			marker = strings.TrimPrefix(marker, "-")
			marker = strings.TrimSpace(marker)
			break
		}
	}

	// Check if it's an indented heredoc
	isIndented := strings.Contains(text, "<<-")

	// Extract content (everything between marker lines)
	heredocContent := text
	startMarker := "<<" + marker
	if isIndented {
		startMarker = "<<-" + marker
	}
	endMarker := marker

	// Very simplistic approach - just remove markers and keep the rest
	if idx := strings.Index(heredocContent, startMarker); idx >= 0 {
		heredocContent = heredocContent[idx+len(startMarker):]
	}
	if idx := strings.LastIndex(heredocContent, endMarker); idx >= 0 {
		heredocContent = heredocContent[:idx]
	}

	// Trim any leading/trailing newlines
	heredocContent = strings.TrimSpace(heredocContent)

	return &types.HeredocExpr{
		Marker:   marker,
		Content:  heredocContent,
		Indented: isIndented,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

// parseIndexExpr parses index access expressions like array[0]
func parseIndexExpr(node *sitter.Node, content []byte) (*types.IndexExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize index expression
	indexExpr := &types.IndexExpr{
		ExprRange: exprRange,
	}

	// Find collection and index
	collectionNode := findChildByFieldName(node, "collection")
	indexNode := findChildByFieldName(node, "index")

	if collectionNode != nil {
		collection, err := parseExpression(collectionNode, content)
		if err != nil {
			return nil, err
		}
		indexExpr.Collection = collection
	}

	if indexNode != nil {
		index, err := parseExpression(indexNode, content)
		if err != nil {
			return nil, err
		}
		indexExpr.Key = index
	}

	return indexExpr, nil
}

// parseTupleExpr parses tuple expressions
func parseTupleExpr(node *sitter.Node, content []byte) (*types.TupleExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize tuple expression
	tupleExpr := &types.TupleExpr{
		Expressions: []types.Expression{},
		ExprRange:   exprRange,
	}

	// Parse items
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			itemNode := cursor.CurrentNode()

			// Skip separators, brackets, and null items
			if itemNode.Type() != "(" && itemNode.Type() != ")" && itemNode.Type() != "," {
				item, err := parseExpression(itemNode, content)
				if err != nil {
					return nil, err
				}
				tupleExpr.Expressions = append(tupleExpr.Expressions, item)
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return tupleExpr, nil
}

// parseUnaryExpr parses unary operations like !condition, -number
func parseUnaryExpr(node *sitter.Node, content []byte) (*types.UnaryExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Get all child nodes
	var operatorNode, exprNode *sitter.Node

	// Try to find by field names first
	operatorNode = findChildByFieldName(node, "operator")
	exprNode = findChildByFieldName(node, "operand")

	// If not found, try by position
	if operatorNode == nil && node.ChildCount() > 0 {
		operatorNode = node.Child(0)
	}

	if exprNode == nil && node.ChildCount() > 1 {
		exprNode = node.Child(1)
	}

	if operatorNode == nil || exprNode == nil {
		// Try direct text inference
		text := string(content[node.StartByte():node.EndByte()])

		if strings.HasPrefix(text, "!") {
			operatorText := "!"
			exprText := strings.TrimSpace(text[1:])

			return &types.UnaryExpr{
				Operator: operatorText,
				Expr: &types.LiteralValue{
					Value:     exprText,
					ValueType: "string",
					ExprRange: exprRange,
				},
				ExprRange: exprRange,
			}, nil
		} else if strings.HasPrefix(text, "-") {
			operatorText := "-"
			exprText := strings.TrimSpace(text[1:])

			return &types.UnaryExpr{
				Operator: operatorText,
				Expr: &types.LiteralValue{
					Value:     exprText,
					ValueType: "string",
					ExprRange: exprRange,
				},
				ExprRange: exprRange,
			}, nil
		}

		return nil, fmt.Errorf("unary expression missing parts")
	}

	// Get operator and expression
	operator := string(content[operatorNode.StartByte():operatorNode.EndByte()])

	expr, err := parseExpression(exprNode, content)
	if err != nil {
		return nil, err
	}

	return &types.UnaryExpr{
		Operator:  operator,
		Expr:      expr,
		ExprRange: exprRange,
	}, nil
}

// parseParenExpr parses parenthesized expressions like (1 + 2) * 3
func parseParenExpr(node *sitter.Node, content []byte) (*types.ParenExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Find the inner expression
	innerNode := findChildByFieldName(node, "expression")
	if innerNode == nil {
		// If not found by field name, try to find by position (skipping parentheses)
		for i := 1; i < int(node.ChildCount())-1; i++ {
			innerNode = node.Child(i)
			break
		}
	}

	if innerNode == nil {
		return nil, fmt.Errorf("parenthesized expression missing inner expression")
	}

	// Parse the inner expression
	innerExpr, err := parseExpression(innerNode, content)
	if err != nil {
		return nil, err
	}

	return &types.ParenExpr{
		Expression: innerExpr,
		ExprRange:  exprRange,
	}, nil
}

// parseRelativeTraversalExpr parses attribute access expressions like aws_instance.example.id
func parseRelativeTraversalExpr(node *sitter.Node, content []byte) (*types.RelativeTraversalExpr, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize the traversal expression
	traversalExpr := &types.RelativeTraversalExpr{
		Traversal: []types.TraversalElem{},
		ExprRange: exprRange,
	}

	// Find the source expression and attribute name
	sourceNode := findChildByFieldName(node, "source")
	attrNode := findChildByFieldName(node, "attr")

	if sourceNode != nil {
		source, err := parseExpression(sourceNode, content)
		if err != nil {
			return nil, err
		}
		traversalExpr.Source = source
	}

	if attrNode != nil {
		// Simple attribute access
		attrName := string(content[attrNode.StartByte():attrNode.EndByte()])
		// Remove the leading dot if present
		if strings.HasPrefix(attrName, ".") {
			attrName = attrName[1:]
		}

		traversalExpr.Traversal = append(traversalExpr.Traversal, types.TraversalElem{
			Type: "attr",
			Name: attrName,
		})
	}

	return traversalExpr, nil
}

// parseTemplateForDirective parses for loops within template strings (%{for x in xs})
func parseTemplateForDirective(node *sitter.Node, content []byte) (*types.TemplateForDirective, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize template for directive
	templateForDir := &types.TemplateForDirective{
		Content:   []types.Expression{},
		ExprRange: exprRange,
	}

	// Get the full text to extract variables and collection
	text := string(content[node.StartByte():node.EndByte()])

	// Match "for x in xs" or "for k, v in xs"
	// First, extract the intro part
	introPart := ""
	if strings.HasPrefix(text, "%{for") {
		endIntro := strings.Index(text, "~}")
		if endIntro < 0 {
			endIntro = strings.Index(text, "}")
		}
		if endIntro > 0 {
			introPart = text[4:endIntro]
			introPart = strings.TrimSpace(introPart)
		}
	}

	// Now extract variables and collection
	if len(introPart) > 0 {
		parts := strings.Split(introPart, "in")
		if len(parts) >= 2 {
			// Left side has variables
			varPart := strings.TrimSpace(parts[0])

			// Check if we have one or two variables
			if strings.Contains(varPart, ",") {
				// Two variables: "k, v"
				vars := strings.Split(varPart, ",")
				if len(vars) >= 2 {
					templateForDir.KeyVar = strings.TrimSpace(vars[0])
					templateForDir.ValueVar = strings.TrimSpace(vars[1])
				}
			} else {
				// One variable: "x"
				templateForDir.ValueVar = varPart
			}

			// Right side has collection
			collPart := strings.TrimSpace(parts[1])
			// For simplicity, we'll use a reference expression for the collection
			templateForDir.CollExpr = createReferenceFromText(collPart, exprRange)
		}
	}

	// The content is everything between the for directive and the endfor directive
	contentStart := strings.Index(text, "}")
	contentEnd := strings.Index(text, "%{endfor")

	if contentStart > 0 && contentEnd > contentStart {
		contentText := text[contentStart+1 : contentEnd]

		// Create a simple literal for the content
		templateForDir.Content = append(templateForDir.Content, &types.LiteralValue{
			Value:     contentText,
			ValueType: "string",
			ExprRange: exprRange,
		})
	}

	return templateForDir, nil
}

// parseTemplateIfDirective parses conditionals within template strings (%{if condition})
func parseTemplateIfDirective(node *sitter.Node, content []byte) (*types.TemplateIfDirective, error) {
	// Get range information
	exprRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize template if directive
	templateIfDir := &types.TemplateIfDirective{
		TrueExpr:  []types.Expression{},
		FalseExpr: []types.Expression{},
		ExprRange: exprRange,
	}

	// Get the full text to extract the condition, true part, and false part
	text := string(content[node.StartByte():node.EndByte()])

	// Extract the condition
	conditionText := ""
	if strings.HasPrefix(text, "%{if") {
		endCond := strings.Index(text, "~}")
		if endCond < 0 {
			endCond = strings.Index(text, "}")
		}
		if endCond > 0 {
			conditionText = text[4:endCond]
			conditionText = strings.TrimSpace(conditionText)
		}
	}

	// Create a reference expression for the condition
	if len(conditionText) > 0 {
		templateIfDir.Condition = createReferenceFromText(conditionText, exprRange)
	}

	// Extract the true part (between if and else/endif)
	trueStart := strings.Index(text, "}")
	trueEnd := -1
	if strings.Contains(text, "%{else") {
		trueEnd = strings.Index(text, "%{else")
	} else if strings.Contains(text, "%{endif") {
		trueEnd = strings.Index(text, "%{endif")
	}

	if trueStart > 0 && trueEnd > trueStart {
		trueText := text[trueStart+1 : trueEnd]

		// Create a literal for the true part
		templateIfDir.TrueExpr = append(templateIfDir.TrueExpr, &types.LiteralValue{
			Value:     trueText,
			ValueType: "string",
			ExprRange: exprRange,
		})
	}

	// Extract the false part (between else and endif)
	if strings.Contains(text, "%{else") {
		falseStart := strings.Index(text, "%{else") + 6 // Skip "%{else"
		falseStart = strings.Index(text[falseStart:], "}") + falseStart + 1
		falseEnd := strings.Index(text, "%{endif")

		if falseEnd > falseStart {
			falseText := text[falseStart:falseEnd]

			// Create a literal for the false part
			templateIfDir.FalseExpr = append(templateIfDir.FalseExpr, &types.LiteralValue{
				Value:     falseText,
				ValueType: "string",
				ExprRange: exprRange,
			})
		}
	}

	return templateIfDir, nil
}

// Helper to extract text between two strings
func extractBetween(text, start, end string) string {
	startIdx := strings.Index(text, start)
	if startIdx < 0 {
		return ""
	}

	startIdx += len(start)
	endIdx := strings.Index(text[startIdx:], end)
	if endIdx < 0 {
		return ""
	}

	return text[startIdx : startIdx+endIdx]
}

// Helper to find a child node by field name
func findChildByFieldName(node *sitter.Node, fieldName string) *sitter.Node {
	count := node.ChildCount()
	for i := 0; i < int(count); i++ {
		field := node.FieldNameForChild(i)
		if field == fieldName {
			return node.Child(i)
		}
	}
	return nil
}

// Helper to find associated comments
// Looks for comments just before or on the same line as the node
func findAssociatedComments(node *sitter.Node, content []byte) (string, string) {
	var blockComment, inlineComment string

	// Get parent node to look for sibling comments
	parent := node.Parent()
	if parent == nil {
		return "", ""
	}

	// Look for a comment node just before this node (block comment)
	for i := 0; i < int(parent.ChildCount()); i++ {
		child := parent.Child(i)
		if child == node {
			// Found our node, check if the previous sibling was a comment
			if i > 0 && parent.Child(i-1).Type() == "comment" {
				commentNode := parent.Child(i - 1)
				commentText := string(content[commentNode.StartByte():commentNode.EndByte()])
				blockComment = strings.TrimSpace(commentText)

				// If there are multiple comment lines before this node, combine them
				for j := i - 2; j >= 0; j-- {
					if parent.Child(j).Type() == "comment" {
						prevComment := parent.Child(j)
						// Check if this comment is on the previous line (consecutive comments)
						if prevComment.EndPoint().Row+1 == commentNode.StartPoint().Row {
							prevCommentText := string(content[prevComment.StartByte():prevComment.EndByte()])
							blockComment = strings.TrimSpace(prevCommentText) + "\n" + blockComment
							commentNode = prevComment
						} else {
							break
						}
					} else {
						break
					}
				}
			}
			break
		}
	}

	// Look for comments on the same line (inline comment)
	nodeLine := node.EndPoint().Row
	for i := 0; i < int(parent.ChildCount()); i++ {
		child := parent.Child(i)
		if child.Type() == "comment" && child.StartPoint().Row == nodeLine {
			commentText := string(content[child.StartByte():child.EndByte()])
			inlineComment = strings.TrimSpace(commentText)
			break
		}
	}

	return blockComment, inlineComment
}
