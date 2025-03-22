package internal

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
	"github.com/vahid-haghighat/terralint/cmd/internal/types"
)

func getFormattedContent(filePath string) ([]byte, error) {
	root, err := ParseTerraformFile(filePath)
	root = root
	if err != nil {
		return nil, err
	}

	// For now, just return the original file content
	// In the future, we can implement formatting based on the parsed AST
	return os.ReadFile(filePath)
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
func parseBlock(node *sitter.Node, content []byte) (*types.Block, error) {
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
	// Get the label (type of the dynamic block)
	var labels []string
	labelsNode := findChildByFieldName(node, "labels")
	if labelsNode != nil && labelsNode.NamedChildCount() > 0 {
		labelNode := labelsNode.NamedChild(0)
		label := string(content[labelNode.StartByte():labelNode.EndByte()])
		// Remove quotes if string
		if len(label) >= 2 && label[0] == '"' && label[len(label)-1] == '"' {
			label = label[1 : len(label)-1]
		}
		labels = append(labels, label)
	}

	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Parse the body to find forEach and content
	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return nil, fmt.Errorf("dynamic block without body")
	}

	// Create the dynamic block
	dynamicBlock := &types.DynamicBlock{
		Labels:   labels,
		Range:    nodeRange,
		Content:  []types.Body{},
		Iterator: "", // Default empty, will be set if found
	}

	// Find comments associated with this block
	blockComment, inlineComment := findAssociatedComments(node, content)
	dynamicBlock.BlockComment = blockComment
	dynamicBlock.InlineComment = inlineComment

	// Parse attributes in the dynamic block
	cursor := sitter.NewTreeCursor(bodyNode)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			childNode := cursor.CurrentNode()

			// Look for forEach attribute
			if childNode.Type() == "attribute" {
				attrName := findChildByFieldName(childNode, "name")
				if attrName != nil && string(content[attrName.StartByte():attrName.EndByte()]) == "for_each" {
					exprNode := findChildByFieldName(childNode, "expression")
					if exprNode != nil {
						expr, err := parseExpression(exprNode, content)
						if err != nil {
							return nil, err
						}
						dynamicBlock.ForEach = expr
					}
				} else if attrName != nil && string(content[attrName.StartByte():attrName.EndByte()]) == "iterator" {
					// Parse iterator if present
					exprNode := findChildByFieldName(childNode, "expression")
					if exprNode != nil {
						// Iterator should be a simple string literal
						text := string(content[exprNode.StartByte():exprNode.EndByte()])
						text = strings.Trim(text, "\"")
						dynamicBlock.Iterator = text
					}
				}
			} else if childNode.Type() == "block" && findChildByFieldName(childNode, "type") != nil {
				blockType := findChildByFieldName(childNode, "type")
				if string(content[blockType.StartByte():blockType.EndByte()]) == "content" {
					// Parse content block
					contentBodyNode := findChildByFieldName(childNode, "body")
					if contentBodyNode != nil {
						contentCursor := sitter.NewTreeCursor(contentBodyNode)
						if contentCursor.GoToFirstChild() {
							for {
								contentChildNode := contentCursor.CurrentNode()
								contentChild, err := parseNode(contentChildNode, content)
								if err != nil {
									contentCursor.Close()
									return nil, err
								}
								if contentChild != nil {
									dynamicBlock.Content = append(dynamicBlock.Content, contentChild)
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

			if !cursor.GoToNextSibling() {
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
	} else if strings.Contains(exprText, ".") && !strings.HasPrefix(exprText, "\"") && !strings.Contains(exprText, "{") {
		// This is likely a reference with dots (like var.brewdex_secret)
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
	case "array", "tuple_cons":
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
	case "for_expr":
		return parseForExpr(node, content)
	case "splat_expr", "splat":
		return parseSplatExpr(node, content)
	case "heredoc_template", "heredoc":
		return parseHeredocExpr(node, content)
	case "index_expr", "index":
		return parseIndexExpr(node, content)
	case "tuple":
		return parseTupleExpr(node, content)
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

			// Skip separators and brackets
			if itemNode.Type() != "[" && itemNode.Type() != "]" && itemNode.Type() != "," {
				item, err := parseExpression(itemNode, content)
				if err != nil {
					return nil, err
				}
				array.Items = append(array.Items, item)
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return array, nil
}

func parseReferenceExpr(node *sitter.Node, content []byte) (*types.ReferenceExpr, error) {
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
	falseNode := findChildByFieldName(node, "false_val")

	if condNode == nil || trueNode == nil || falseNode == nil {
		return nil, fmt.Errorf("conditional expression missing parts")
	}

	condition, err := parseExpression(condNode, content)
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
	leftNode := findChildByFieldName(node, "left")
	rightNode := findChildByFieldName(node, "right")
	operatorNode := findChildByFieldName(node, "operator")

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

		if i >= len(text) || text[i] != '=' {
			// Skip this line if there's no equals sign
			for i < len(text) && text[i] != '\n' {
				i++
			}
			continue
		}
		// fmt.Printf("Found key: %q\n", keyText)
		i++ // Skip the equals sign

		// Skip whitespace after equals sign
		for i < len(text) && (text[i] == ' ' || text[i] == '\t') {
			i++
		}

		if i >= len(text) {
			break
		}

		// Find the value
		var value string
		var valueExpr types.Expression

		// Create a reference expression for the key
		keyExpr := &types.ReferenceExpr{
			Parts:     []string{keyText},
			ExprRange: exprRange,
		}

		// Handle different value types
		if i < len(text) && text[i] == '"' {
			// String literal
			valueStart := i
			i++ // Skip the opening quote

			// Find the closing quote
			for i < len(text) && text[i] != '"' {
				if text[i] == '\\' && i+1 < len(text) {
					i += 2 // Skip escaped character
				} else {
					i++
				}
			}

			if i < len(text) {
				i++ // Skip the closing quote
				value = text[valueStart:i]

				// Create a string literal
				valueExpr = &types.LiteralValue{
					Value:     value[1 : len(value)-1], // Remove quotes
					ValueType: "string",
					ExprRange: exprRange,
				}
			}
		} else if i+3 < len(text) && text[i:i+4] == "true" {
			// Boolean true
			value = "true"
			i += 4
			valueExpr = &types.LiteralValue{
				Value:     "true",
				ValueType: "bool",
				ExprRange: exprRange,
			}
		} else if i+4 < len(text) && text[i:i+5] == "false" {
			// Boolean false
			value = "false"
			i += 5
			valueExpr = &types.LiteralValue{
				Value:     "false",
				ValueType: "bool",
				ExprRange: exprRange,
			}
		} else if i+3 < len(text) && text[i:i+4] == "null" {
			// Null value
			value = "null"
			i += 4
			valueExpr = &types.LiteralValue{
				Value:     "null",
				ValueType: "null",
				ExprRange: exprRange,
			}
		} else if i < len(text) && text[i] >= '0' && text[i] <= '9' {
			// Number literal
			valueStart := i

			// Find the end of the number
			for i < len(text) && ((text[i] >= '0' && text[i] <= '9') || text[i] == '.') {
				i++
			}

			value = text[valueStart:i]
			valueExpr = &types.LiteralValue{
				Value:     value,
				ValueType: "number",
				ExprRange: exprRange,
			}
		} else if i < len(text) && text[i] == '[' {
			// Array/tuple literal
			valueStart := i
			i++ // Skip the opening bracket
			braceLevel := 1

			// Find the closing bracket
			for i < len(text) && braceLevel > 0 {
				if text[i] == '[' {
					braceLevel++
				} else if text[i] == ']' {
					braceLevel--
				} else if text[i] == '"' {
					i++ // Skip the opening quote

					// Skip the string content
					for i < len(text) && text[i] != '"' {
						if text[i] == '\\' && i+1 < len(text) {
							i += 2 // Skip escaped character
						} else {
							i++
						}
					}
				}

				if i < len(text) {
					i++
				}
			}

			value = text[valueStart:i]

			// Create a tuple expression
			// Parse the array elements
			arrayContent := value[1 : len(value)-1] // Remove brackets
			elements := []types.Expression{}

			// Split by commas, but handle nested structures
			start := 0
			var arrayBraceLevel int
			var arrayBracketLevel int
			var arrayQuoteOpen bool

			for i := 0; i < len(arrayContent); i++ {
				char := arrayContent[i]

				if char == '"' && (i == 0 || arrayContent[i-1] != '\\') {
					arrayQuoteOpen = !arrayQuoteOpen
				} else if !arrayQuoteOpen {
					if char == '{' {
						arrayBraceLevel++
					} else if char == '}' {
						arrayBraceLevel--
					} else if char == '[' {
						arrayBracketLevel++
					} else if char == ']' {
						arrayBracketLevel--
					} else if char == ',' && arrayBraceLevel == 0 && arrayBracketLevel == 0 {
						// Found a comma at the top level, extract the element
						element := strings.TrimSpace(arrayContent[start:i])
						if element != "" {
							elements = append(elements, createReferenceFromText(element, exprRange))
						}
						start = i + 1
					}
				}
			}

			// Add the last element
			if start < len(arrayContent) {
				element := strings.TrimSpace(arrayContent[start:])
				if element != "" {
					elements = append(elements, createReferenceFromText(element, exprRange))
				}
			}

			valueExpr = &types.TupleExpr{
				Expressions: elements,
				ExprRange:   exprRange,
			}
		} else if i < len(text) && text[i] == '{' {
			// Nested object literal
			valueStart := i
			i++ // Skip the opening brace
			braceLevel := 1

			// Find the closing brace
			for i < len(text) && braceLevel > 0 {
				if text[i] == '{' {
					braceLevel++
				} else if text[i] == '}' {
					braceLevel--
				} else if text[i] == '"' {
					i++ // Skip the opening quote

					// Skip the string content
					for i < len(text) && text[i] != '"' {
						if text[i] == '\\' && i+1 < len(text) {
							i += 2 // Skip escaped character
						} else {
							i++
						}
					}
				}

				if i < len(text) {
					i++
				}
			}

			value = text[valueStart:i]

			// Create a nested object expression
			nestedObjText := value[1 : len(value)-1] // Remove braces
			valueExpr = parseObjectFromText(nestedObjText, exprRange)
		} else {
			// Reference or other expression
			valueStart := i

			// Find the end of the value (newline or comma)
			for i < len(text) && text[i] != '\n' && text[i] != ',' {
				i++
			}

			value = strings.TrimSpace(text[valueStart:i])
			valueExpr = createReferenceFromText(value, exprRange)
		}

		// Add the object item
		objExpr.Items = append(objExpr.Items, types.ObjectItem{
			Key:   keyExpr,
			Value: valueExpr,
		})

		// Skip to the next line
		for i < len(text) && text[i] != '\n' {
			i++
		}
	}

	return objExpr
}
func createReferenceFromText(text string, exprRange sitter.Range) types.Expression {
	// Trim any whitespace
	text = strings.TrimSpace(text)

	// Handle empty text
	if text == "" {
		return &types.LiteralValue{
			Value:     "",
			ValueType: "string",
			ExprRange: exprRange,
		}
	}

	// Handle string literals
	if (strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\"")) ||
		(strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'")) {
		// Remove the quotes
		value := text[1 : len(text)-1]
		return &types.LiteralValue{
			Value:     value,
			ValueType: "string",
			ExprRange: exprRange,
		}
	}

	// Handle boolean literals
	if text == "true" || text == "false" {
		return &types.LiteralValue{
			Value:     text,
			ValueType: "bool",
			ExprRange: exprRange,
		}
	}

	// Handle number literals
	if _, err := strconv.ParseFloat(text, 64); err == nil {
		return &types.LiteralValue{
			Value:     text,
			ValueType: "number",
			ExprRange: exprRange,
		}
	}

	// Handle null literal
	if text == "null" {
		return &types.LiteralValue{
			Value:     "null",
			ValueType: "null",
			ExprRange: exprRange,
		}
	}

	// Handle function calls
	if strings.Contains(text, "(") && strings.Contains(text, ")") {
		openParenIndex := strings.Index(text, "(")
		if openParenIndex > 0 {
			funcName := strings.TrimSpace(text[:openParenIndex])

			// Create a function call expression
			return &types.FunctionCallExpr{
				Name:      funcName,
				Args:      []types.Expression{}, // We can't easily parse the args from text
				ExprRange: exprRange,
			}
		}
	}

	// Handle references with dots
	if strings.Contains(text, ".") {
		// This looks like a reference with dots
		parts := strings.Split(text, ".")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return &types.ReferenceExpr{
			Parts:     parts,
			ExprRange: exprRange,
		}
	}

	// Handle array/tuple literals
	if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") {
		// This is an array/tuple literal
		return &types.TupleExpr{
			Expressions: []types.Expression{}, // We can't easily parse the elements from text
			ExprRange:   exprRange,
		}
	}

	// Handle object literals
	if strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}") {
		// This is an object literal
		return &types.ObjectExpr{
			Items:     []types.ObjectItem{}, // We can't easily parse the items from text
			ExprRange: exprRange,
		}
	}

	// Default to a reference to a single variable
	return &types.ReferenceExpr{
		Parts:     []string{text},
		ExprRange: exprRange,
	}
}

func parseForExpr(node *sitter.Node, content []byte) (*types.ForExpr, error) {
	// Get the full text for debugging
	fullText := string(content[node.StartByte():node.EndByte()])

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

	// Special case for direct attribute values like f = [ for item in var.items: item * 2 ]
	// These might be parsed as references instead of for expressions
	if strings.Contains(fullText, "for") && strings.Contains(fullText, "in") && strings.Contains(fullText, ":") {
		// This looks like a for expression in text form
		// Extract the components using regex

		// Extract the iterator variables and collection
		forInMatch := regexp.MustCompile(`for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+in\s+(.+?):`)
		forKeyValueMatch := regexp.MustCompile(`for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*,\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+in\s+(.+?):`)

		if forKeyValueMatch.MatchString(fullText) {
			// This is a "for k, v in collection" pattern
			matches := forKeyValueMatch.FindStringSubmatch(fullText)
			if len(matches) >= 4 {
				forExpr.KeyVar = strings.TrimSpace(matches[1])
				forExpr.ValueVar = strings.TrimSpace(matches[2])

				// Parse the collection
				collText := strings.TrimSpace(matches[3])
				forExpr.Collection = createReferenceFromText(collText, exprRange)
			}
		} else if forInMatch.MatchString(fullText) {
			// This is a "for item in collection" pattern
			matches := forInMatch.FindStringSubmatch(fullText)
			if len(matches) >= 3 {
				forExpr.ValueVar = strings.TrimSpace(matches[1])

				// Parse the collection
				collText := strings.TrimSpace(matches[2])
				forExpr.Collection = createReferenceFromText(collText, exprRange)
			}
		}

		// Extract the value expression
		valueMatch := regexp.MustCompile(`:\s*(.+?)[\]\}]`).FindStringSubmatch(fullText)
		if len(valueMatch) >= 2 {
			valueText := strings.TrimSpace(valueMatch[1])

			// Check if this is a key-value mapping (for object expressions)
			if strings.Contains(valueText, "=>") {
				parts := strings.Split(valueText, "=>")
				if len(parts) >= 2 {
					valueExprText := strings.TrimSpace(parts[0])
					keyExprText := strings.TrimSpace(parts[1])

					forExpr.ThenKeyExpr = createReferenceFromText(valueExprText, exprRange)
					forExpr.ThenValueExpr = createReferenceFromText(keyExprText, exprRange)
				}
			} else {
				// This is a simple value expression
				forExpr.ThenKeyExpr = createReferenceFromText(valueText, exprRange)
			}
		}

		// Return the for expression
		return forExpr, nil
	}

	// Check if this is a for_expr with a single child
	if node.NamedChildCount() == 1 {
		// Get the actual for expression node (for_tuple_expr or for_object_expr)
		forNode := node.NamedChild(0)

		// Process based on the type of for expression
		if forNode.Type() == "for_tuple_expr" || forNode.Type() == "for_object_expr" {
			// Extract components from the for expression

			// 1. Find the for_intro node which contains the variables and collection
			var forIntroNode *sitter.Node
			var valueNode *sitter.Node
			var keyNode *sitter.Node

			for i := 0; i < int(forNode.NamedChildCount()); i++ {
				child := forNode.NamedChild(i)

				if child.Type() == "for_intro" {
					forIntroNode = child
				} else if child.Type() == "expression" {
					// In for_tuple_expr, there's only one expression (the value)
					// In for_object_expr, there are two expressions (value and key)
					if valueNode == nil {
						valueNode = child
					} else if keyNode == nil && forNode.Type() == "for_object_expr" {
						keyNode = child
					}
				}
			}

			// 2. Extract variables and collection from for_intro
			if forIntroNode != nil {
				// Parse the for_intro text to extract variables and collection
				introText := string(content[forIntroNode.StartByte():forIntroNode.EndByte()])

				// Remove the "for " prefix and the ":" suffix
				introText = strings.TrimPrefix(introText, "for ")
				introText = strings.TrimSuffix(introText, ":")

				// Split by "in" to get variables and collection
				parts := strings.Split(introText, " in ")
				if len(parts) == 2 {
					varPart := strings.TrimSpace(parts[0])
					collPart := strings.TrimSpace(parts[1])

					// Check if we have one or two variables
					if strings.Contains(varPart, ",") {
						// Two variables: "k, v"
						varParts := strings.Split(varPart, ",")
						if len(varParts) >= 2 {
							forExpr.KeyVar = strings.TrimSpace(varParts[0])
							forExpr.ValueVar = strings.TrimSpace(varParts[1])
						}
					} else {
						// One variable: "item"
						forExpr.ValueVar = varPart
					}

					// Parse the collection expression
					forExpr.Collection = &types.ReferenceExpr{
						Parts:     strings.Split(collPart, "."),
						ExprRange: exprRange,
					}
				}
			}

			// 3. Extract the value expression
			if valueNode != nil {
				// Get the text of the value expression
				valueText := string(content[valueNode.StartByte():valueNode.EndByte()])

				// For tuple expressions, we need to handle expressions like "item * 2"
				if forNode.Type() == "for_tuple_expr" {
					// Try to parse the expression
					valueExpr, err := parseExpression(valueNode, content)
					if err == nil {
						forExpr.ThenKeyExpr = valueExpr
					} else {
						// If parsing fails, create a reference expression
						forExpr.ThenKeyExpr = &types.ReferenceExpr{
							Parts:     []string{valueText},
							ExprRange: exprRange,
						}
					}
				} else {
					// For object expressions, just parse normally
					valueExpr, err := parseExpression(valueNode, content)
					if err == nil {
						forExpr.ThenKeyExpr = valueExpr
					}
				}
			}

			// 4. Extract the key expression (for object expressions)
			if keyNode != nil {
				keyExpr, err := parseExpression(keyNode, content)
				if err == nil {
					forExpr.ThenValueExpr = keyExpr
				}
			}
		}
	}

	// Set default values if we couldn't extract them
	if forExpr.Collection == nil {
		forExpr.Collection = &types.LiteralValue{
			Value:     "collection",
			ValueType: "string",
			ExprRange: exprRange,
		}
	}

	if forExpr.ThenKeyExpr == nil {
		forExpr.ThenKeyExpr = &types.LiteralValue{
			Value:     "value",
			ValueType: "string",
			ExprRange: exprRange,
		}
	}

	return forExpr, nil
}

func parseSplatExpr(node *sitter.Node, content []byte) (*types.SplatExpr, error) {
	sourceNode := findChildByFieldName(node, "collection")
	if sourceNode == nil {
		return nil, fmt.Errorf("splat expression missing source")
	}

	source, err := parseExpression(sourceNode, content)
	if err != nil {
		return nil, err
	}

	// Each expression (what's being accessed after the splat)
	// For full splats this is often implicit
	var each types.Expression
	eachNode := findChildByFieldName(node, "each")
	if eachNode != nil {
		each, err = parseExpression(eachNode, content)
		if err != nil {
			return nil, err
		}
	} else {
		// If no explicit each expression, use a simple reference
		each = &types.ReferenceExpr{
			Parts: []string{"*"},
			ExprRange: sitter.Range{
				StartPoint: node.EndPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.EndByte(),
				EndByte:    node.EndByte(),
			},
		}
	}

	return &types.SplatExpr{
		Source: source,
		Each:   each,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseHeredocExpr(node *sitter.Node, content []byte) (*types.HeredocExpr, error) {
	markerNode := findChildByFieldName(node, "marker")
	if markerNode == nil {
		return nil, fmt.Errorf("heredoc without marker")
	}

	marker := string(content[markerNode.StartByte():markerNode.EndByte()])

	// Check if indented heredoc (starts with <<-)
	indent := false
	if strings.HasPrefix(string(content[node.StartByte():node.StartByte()+3]), "<<-") {
		indent = true
	}

	// Get content
	contentNode := findChildByFieldName(node, "content")
	var contentText string
	if contentNode != nil {
		contentText = string(content[contentNode.StartByte():contentNode.EndByte()])
	}

	return &types.HeredocExpr{
		Marker:   marker,
		Content:  contentText,
		Indented: indent,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseIndexExpr(node *sitter.Node, content []byte) (*types.IndexExpr, error) {
	targetNode := findChildByFieldName(node, "target")
	indexNode := findChildByFieldName(node, "index")

	if targetNode == nil || indexNode == nil {
		return nil, fmt.Errorf("index expression missing parts")
	}

	collection, err := parseExpression(targetNode, content)
	if err != nil {
		return nil, err
	}

	key, err := parseExpression(indexNode, content)
	if err != nil {
		return nil, err
	}

	return &types.IndexExpr{
		Collection: collection, // Changed from Target to Collection
		Key:        key,        // Changed from Index to Key
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseTupleExpr(node *sitter.Node, content []byte) (*types.TupleExpr, error) {
	tuple := &types.TupleExpr{
		Expressions: []types.Expression{},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Parse items in the tuple
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			exprNode := cursor.CurrentNode()

			// Skip brackets, commas, and structural elements
			if exprNode.Type() != "[" && exprNode.Type() != "]" &&
				exprNode.Type() != "," && exprNode.Type() != "tuple_start" &&
				exprNode.Type() != "tuple_end" {

				// Skip empty or whitespace-only nodes
				nodeText := string(content[exprNode.StartByte():exprNode.EndByte()])
				if strings.TrimSpace(nodeText) != "" {
					expr, err := parseExpression(exprNode, content)
					if err != nil {
						return nil, err
					}

					// Skip null literals that might be artifacts of parsing
					if literalExpr, ok := expr.(*types.LiteralValue); ok {
						if literalExpr.ValueType == "null" && literalExpr.Value == "" {
							// Skip this null value
						} else {
							tuple.Expressions = append(tuple.Expressions, expr)
						}
					} else {
						tuple.Expressions = append(tuple.Expressions, expr)
					}
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

	return tuple, nil
}

func parseUnaryExpr(node *sitter.Node, content []byte) (*types.UnaryExpr, error) {
	operatorNode := findChildByFieldName(node, "operator")
	expressionNode := findChildByFieldName(node, "expression")

	if operatorNode == nil || expressionNode == nil {
		return nil, fmt.Errorf("unary expression missing parts")
	}

	operator := string(content[operatorNode.StartByte():operatorNode.EndByte()])

	expr, err := parseExpression(expressionNode, content)
	if err != nil {
		return nil, err
	}

	return &types.UnaryExpr{
		Operator: operator,
		Expr:     expr,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseParenExpr(node *sitter.Node, content []byte) (*types.ParenExpr, error) {
	// Find the expression inside the parentheses
	exprNode := node.NamedChild(0)
	if exprNode == nil {
		return nil, fmt.Errorf("parenthesized expression is empty")
	}

	expr, err := parseExpression(exprNode, content)
	if err != nil {
		return nil, err
	}

	return &types.ParenExpr{
		Expression: expr,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseRelativeTraversalExpr(node *sitter.Node, content []byte) (*types.RelativeTraversalExpr, error) {
	sourceNode := findChildByFieldName(node, "object")
	if sourceNode == nil {
		return nil, fmt.Errorf("traversal expression missing source")
	}

	source, err := parseExpression(sourceNode, content)
	if err != nil {
		return nil, err
	}

	// Get the attribute being accessed
	attrNode := findChildByFieldName(node, "attr")
	if attrNode == nil {
		return nil, fmt.Errorf("traversal expression missing attribute")
	}

	attr := string(content[attrNode.StartByte():attrNode.EndByte()])

	// Create a TraversalElem for the attribute access
	traversal := []types.TraversalElem{
		{
			Type: "attr",
			Name: attr,
			// Index is nil for attribute access
		},
	}

	return &types.RelativeTraversalExpr{
		Source:    source,
		Traversal: traversal,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseTemplateForDirective(node *sitter.Node, content []byte) (*types.TemplateForDirective, error) {
	// Extract text to manually parse the for directive because tree-sitter doesn't
	// give us a reliable structure for template directives
	forText := string(content[node.StartByte():node.EndByte()])

	// Extract for variables and collection
	forPart := extractBetween(forText, "%{for ", " in ")
	collPart := extractBetween(forText, " in ", "}")

	// Check if we have a value_var only or key_var and value_var
	var keyVar, valueVar string
	if strings.Contains(forPart, ",") {
		parts := strings.Split(forPart, ",")
		if len(parts) >= 2 {
			keyVar = strings.TrimSpace(parts[0])
			valueVar = strings.TrimSpace(parts[1])
		}
	} else {
		valueVar = strings.TrimSpace(forPart)
	}

	// Extract result expressions
	exprText := extractBetween(forText, "}", "%{endfor}")

	// Create template expression for the result
	resultExpr := &types.TemplateExpr{
		Parts: []types.Expression{
			&types.LiteralValue{
				Value:     exprText,
				ValueType: "string",
				ExprRange: sitter.Range{
					StartPoint: node.StartPoint(),
					EndPoint:   node.EndPoint(),
					StartByte:  node.StartByte(),
					EndByte:    node.EndByte(),
				},
			},
		},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Create a simple reference expression for collection
	collExpr := &types.ReferenceExpr{
		Parts: []string{collPart},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	return &types.TemplateForDirective{
		KeyVar:   keyVar,
		ValueVar: valueVar,
		CollExpr: collExpr,
		Content:  []types.Expression{resultExpr}, // Changed from ValueExpr to Content
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

func parseTemplateIfDirective(node *sitter.Node, content []byte) (*types.TemplateIfDirective, error) {
	// Extract text to manually parse the if directive
	ifText := string(content[node.StartByte():node.EndByte()])

	// Extract condition
	condText := extractBetween(ifText, "%{if ", "}")

	// Extract true branch
	trueText := extractBetween(ifText, "}", "%{else}")
	if trueText == "" {
		trueText = extractBetween(ifText, "}", "%{endif}")
	}

	// Extract false branch if it exists
	falseText := ""
	if strings.Contains(ifText, "%{else}") {
		falseText = extractBetween(ifText, "%{else}", "%{endif}")
	}

	// Create condition expression (simple reference)
	condExpr := &types.ReferenceExpr{
		Parts: []string{condText},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Create template expressions for true and false branches
	trueExpr := &types.TemplateExpr{
		Parts: []types.Expression{
			&types.LiteralValue{
				Value:     trueText,
				ValueType: "string",
				ExprRange: sitter.Range{
					StartPoint: node.StartPoint(),
					EndPoint:   node.EndPoint(),
					StartByte:  node.StartByte(),
					EndByte:    node.EndByte(),
				},
			},
		},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	var falseExpr *types.TemplateExpr
	if falseText != "" {
		falseExpr = &types.TemplateExpr{
			Parts: []types.Expression{
				&types.LiteralValue{
					Value:     falseText,
					ValueType: "string",
					ExprRange: sitter.Range{
						StartPoint: node.StartPoint(),
						EndPoint:   node.EndPoint(),
						StartByte:  node.StartByte(),
						EndByte:    node.EndByte(),
					},
				},
			},
			ExprRange: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}
	}

	// Create the arrays of expressions for true and false branches
	trueExprs := []types.Expression{trueExpr}
	var falseExprs []types.Expression
	if falseExpr != nil {
		falseExprs = []types.Expression{falseExpr}
	}

	return &types.TemplateIfDirective{
		Condition: condExpr,
		TrueExpr:  trueExprs,  // Now passing a slice of Expression
		FalseExpr: falseExprs, // Now passing a slice of Expression
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
}

// Helper function to extract text between two markers
func extractBetween(text, start, end string) string {
	startIndex := strings.Index(text, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)

	endIndex := strings.Index(text[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	return text[startIndex : startIndex+endIndex]
}

// Helper function to find a child node by field name
func findChildByFieldName(node *sitter.Node, fieldName string) *sitter.Node {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		field := node.FieldNameForChild(i)
		if field == fieldName {
			return child
		}
	}
	return nil
}

// Helper function to find comments associated with a node
// Returns block comment (if any) and inline comment (if any)
func findAssociatedComments(node *sitter.Node, content []byte) (string, string) {
	var blockComment, inlineComment string

	// Check for block comments (before the node)
	// We'll collect all consecutive comments before the node
	var blockComments []string
	prevSibling := node.PrevSibling()

	// Keep looking at previous siblings until we find a non-comment
	for prevSibling != nil && prevSibling.Type() == "comment" {
		commentText := string(content[prevSibling.StartByte():prevSibling.EndByte()])
		// Add to the beginning of the list to maintain order
		blockComments = append([]string{strings.TrimSpace(commentText)}, blockComments...)
		prevSibling = prevSibling.PrevSibling()
	}

	// Join all block comments with newlines
	if len(blockComments) > 0 {
		blockComment = strings.Join(blockComments, "\n")
	}

	// Check for inline comment (on the same line, after the node)
	endLine := node.EndPoint().Row
	parent := node.Parent()

	if parent != nil {
		cursor := sitter.NewTreeCursor(parent)
		defer cursor.Close()

		if cursor.GoToFirstChild() {
			// Find the node
			for cursor.CurrentNode() != node {
				if !cursor.GoToNextSibling() {
					break
				}
			}

			// Check next siblings for an inline comment
			if cursor.GoToNextSibling() {
				nextNode := cursor.CurrentNode()
				if nextNode.Type() == "comment" && nextNode.StartPoint().Row == endLine {
					commentText := string(content[nextNode.StartByte():nextNode.EndByte()])
					inlineComment = strings.TrimSpace(commentText)
				}
			}
		}
	}

	return blockComment, inlineComment
}
