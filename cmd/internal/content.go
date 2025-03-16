package internal

import (
	"context"
	"encoding/json"
	"log"

	//"context"
	"fmt"
	"strconv"
	"strings"

	//"github.com/hashicorp/hcl/v2"
	//"github.com/hashicorp/hcl/v2/hclwrite"
	sitter "github.com/smacker/go-tree-sitter"
	sitterhcl "github.com/smacker/go-tree-sitter/hcl"
	"github.com/vahid-haghighat/terralint/cmd/internal/types"
	"os"
)

func getFormattedContent(filePath string) ([]byte, error) {
	root, err := ParseTerraformFile(filePath)

	if err != nil {
		log.Println(err)
	}
	root = root
	err = err
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

			// Parse the node based on its type
			body, err := parseNode(node, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing node at line %d: %w",
					node.StartPoint().Row+1, err)
			}

			if body != nil {
				if node.Type() == "body" {
					root.Children = body.(*types.Root).Children
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
		return parseAttribute(node, content)
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

	// Extract block comments and inline comments
	blockComment, inlineComment := findAssociatedComments(node, content)
	block.BlockComment = blockComment
	block.InlineComment = inlineComment

	// Process the block's direct children to find identifier (type) and string_literals (labels)
	var body *sitter.Node

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
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil, fmt.Errorf("attribute without name")
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Find expression node
	exprNode := findChildByFieldName(node, "expression")
	if exprNode == nil {
		return nil, fmt.Errorf("attribute without expression")
	}

	expr, err := parseExpression(exprNode, content)
	if err != nil {
		return nil, err
	}

	// Find inline comment if any
	_, inlineComment := findAssociatedComments(node, content)

	attribute := &types.Attribute{
		Name:          name,
		Value:         expr,
		Range:         nodeRange,
		InlineComment: inlineComment,
	}

	return attribute, nil
}

// parseExpression parses an expression node
func parseExpression(node *sitter.Node, content []byte) (types.Expression, error) {
	nodeType := node.Type()

	switch nodeType {
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
	default:
		return nil, fmt.Errorf("unknown expression type: %s", nodeType)
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

	// Find object items
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			itemNode := cursor.CurrentNode()

			if itemNode.Type() == "object_elem" {
				keyNode := findChildByFieldName(itemNode, "key")
				valueNode := findChildByFieldName(itemNode, "value")

				if keyNode != nil && valueNode != nil {
					key, err := parseExpression(keyNode, content)
					if err != nil {
						return nil, err
					}

					value, err := parseExpression(valueNode, content)
					if err != nil {
						return nil, err
					}

					// Find inline comment
					_, inlineComment := findAssociatedComments(itemNode, content)

					item := types.ObjectItem{
						Key:           key,
						Value:         value,
						InlineComment: inlineComment,
					}

					obj.Items = append(obj.Items, item)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}

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
	// Extract the reference parts
	text := string(content[node.StartByte():node.EndByte()])
	parts := strings.Split(text, ".")

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
	nameNode := findChildByFieldName(node, "function")
	if nameNode == nil {
		return nil, fmt.Errorf("function call without name")
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	function := &types.FunctionCallExpr{
		Name: name,
		Args: []types.Expression{},
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}

	// Parse arguments
	argsNode := findChildByFieldName(node, "arguments")
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
						return nil, err
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

func parseForExpr(node *sitter.Node, content []byte) (*types.ForExpr, error) {
	// Extract the key and value variables
	keyVarNode := findChildByFieldName(node, "key_var")
	valueVarNode := findChildByFieldName(node, "value_var")

	var keyVar, valueVar string
	if keyVarNode != nil {
		keyVar = string(content[keyVarNode.StartByte():keyVarNode.EndByte()])
	}

	if valueVarNode != nil {
		valueVar = string(content[valueVarNode.StartByte():valueVarNode.EndByte()])
	}

	// Get collection, value expression, and condition
	collNode := findChildByFieldName(node, "collection")
	if collNode == nil {
		return nil, fmt.Errorf("for expression missing collection")
	}

	collExpr, err := parseExpression(collNode, content)
	if err != nil {
		return nil, err
	}

	// Value expression
	valueExprNode := findChildByFieldName(node, "expr")
	if valueExprNode == nil {
		return nil, fmt.Errorf("for expression missing value expression")
	}

	valueExpr, err := parseExpression(valueExprNode, content)
	if err != nil {
		return nil, err
	}

	// Optional key expression
	var keyExpr types.Expression
	keyExprNode := findChildByFieldName(node, "key_expr")
	if keyExprNode != nil {
		keyExpr, err = parseExpression(keyExprNode, content)
		if err != nil {
			return nil, err
		}
	}

	// Optional condition
	var condition types.Expression
	condNode := findChildByFieldName(node, "condition")
	if condNode != nil {
		condition, err = parseExpression(condNode, content)
		if err != nil {
			return nil, err
		}
	}

	// Check if this is a grouped expression
	isGrouped := false
	if strings.Contains(string(content[node.StartByte():node.EndByte()]), "=>") && strings.Contains(string(content[node.StartByte():node.EndByte()]), "...") {
		isGrouped = true
	}

	return &types.ForExpr{
		KeyVar:    keyVar,
		ValueVar:  valueVar,
		CollExpr:  collExpr,
		KeyExpr:   keyExpr,
		ValueExpr: valueExpr,
		Condition: condition,
		IsGrouped: isGrouped,
		ExprRange: sitter.Range{
			StartPoint: node.StartPoint(),
			EndPoint:   node.EndPoint(),
			StartByte:  node.StartByte(),
			EndByte:    node.EndByte(),
		},
	}, nil
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

			// Skip parentheses and commas
			if exprNode.Type() != "[" && exprNode.Type() != "]" && exprNode.Type() != "," {
				expr, err := parseExpression(exprNode, content)
				if err != nil {
					return nil, err
				}
				tuple.Expressions = append(tuple.Expressions, expr)
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

	// Check for block comment (immediately before the node)
	prevSibling := node.PrevSibling()
	if prevSibling != nil && prevSibling.Type() == "comment" {
		commentText := string(content[prevSibling.StartByte():prevSibling.EndByte()])
		blockComment = strings.TrimSpace(commentText)
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
