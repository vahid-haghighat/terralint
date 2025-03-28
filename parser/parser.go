package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	sitterhcl "github.com/smacker/go-tree-sitter/hcl"
	"github.com/vahid-haghighat/terralint/parser/types"
)

// Enable/disable debug printing
var debugEnabled = true

// debugPrint prints a debug message if debugging is enabled
func debugPrint(format string, args ...interface{}) {
	if debugEnabled {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// truncateString truncates a string to maxLen and adds "..." if needed
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

	case *types.ForExpr:
		var builder strings.Builder
		builder.WriteString("for ")
		if e.KeyVar != "" {
			builder.WriteString(e.KeyVar + ", ")
		}
		builder.WriteString(e.ValueVar + " in ")
		renderedCollection, err := renderExpression(e.Collection)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedCollection)
		builder.WriteString(" : ")

		if e.ThenKeyExpr != nil {
			renderedKey, err := renderExpression(e.ThenKeyExpr)
			if err != nil {
				return "", err
			}
			builder.WriteString(renderedKey)
		}

		if e.ThenValueExpr != nil {
			builder.WriteString(" => ")
			renderedValue, err := renderExpression(e.ThenValueExpr)
			if err != nil {
				return "", err
			}
			builder.WriteString(renderedValue)
		}

		return builder.String(), nil

	case *types.ConditionalExpr:
		var builder strings.Builder
		builder.WriteString("if ")
		renderedCondition, err := renderExpression(e.Condition)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedCondition)
		builder.WriteString(" then ")
		renderedTrueExpr, err := renderExpression(e.TrueExpr)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedTrueExpr)
		builder.WriteString(" else ")
		renderedFalseExpr, err := renderExpression(e.FalseExpr)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedFalseExpr)
		builder.WriteString("\n")
		return builder.String(), nil

	case *types.BinaryExpr:
		var builder strings.Builder
		renderedLeft, err := renderExpression(e.Left)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedLeft)
		builder.WriteString(" ")
		builder.WriteString(e.Operator)
		builder.WriteString(" ")
		renderedRight, err := renderExpression(e.Right)
		if err != nil {
			return "", err
		}
		builder.WriteString(renderedRight)
		return builder.String(), nil

	case *types.ParenExpr:
		renderedExpr, err := renderExpression(e.Expression)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", renderedExpr), nil

	default:
		return "", fmt.Errorf("unhandled expression type: %T", expr)
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
		// Log unknown node types for debugging
		fmt.Printf("Unknown node type: %s\n", nodeType)
		// Skip non-expression nodes like punctuation
		if nodeType == "=" || nodeType == "," || nodeType == "." || nodeType == "[" || nodeType == "]" || nodeType == "{" || nodeType == "}" || nodeType == "(" || nodeType == ")" {
			return nil, nil
		}
		// For any other type of node, try to parse it as an expression
		expr, err := parseExpression(node, content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse unknown node as expression: %w", err)
		}
		// Create an attribute with the expression value
		return &types.Attribute{
			Name:  "value",
			Value: expr,
			Range: sitter.Range{
				StartPoint: node.StartPoint(),
				EndPoint:   node.EndPoint(),
				StartByte:  node.StartByte(),
				EndByte:    node.EndByte(),
			},
		}, nil
	}
}

func parseBodyNode(node *sitter.Node, content []byte) (*types.Root, error) {
	// Create an appropriate container for the body contents
	// This could be a Block, Attribute, or some other type depending on your needs
	body := &types.Root{
		Children: []types.Body{},
	}

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child, err := parseNode(node.NamedChild(i), content)
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
	if node == nil {
		return nil, fmt.Errorf("nil node passed to parseBlock")
	}

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
	if blockComment != "" {
		block.BlockComment = blockComment
	}
	block.InlineComment = inlineComment

	// Create a tree-sitter cursor
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

		// Process remaining children to extract labels and the block body
		for cursor.GoToNextSibling() {
			childNode := cursor.CurrentNode()
			switch childNode.Type() {
			case "string_lit", "quoted_template":
				labelText := string(content[childNode.StartByte():childNode.EndByte()])
				// Remove quotes if present
				labelText = strings.Trim(labelText, "\"'")
				block.Labels = append(block.Labels, labelText)
			case "body":
				// Parse the block body
				bodyCursor := sitter.NewTreeCursor(childNode)
				defer bodyCursor.Close()

				if bodyCursor.GoToFirstChild() {
					for {
						bodyChildNode := bodyCursor.CurrentNode()
						if bodyChildNode == nil {
							break
						}

						// Skip comments and punctuation
						if bodyChildNode.Type() == "comment" || bodyChildNode.Type() == "{" || bodyChildNode.Type() == "}" {
							if !bodyCursor.GoToNextSibling() {
								break
							}
							continue
						}

						var childBody types.Body
						var err error

						switch bodyChildNode.Type() {
						case "block":
							childBody, err = parseBlock(bodyChildNode, content)
						case "attribute":
							childBody, err = parseAttribute(bodyChildNode, content)
						default:
							// Try to parse as a generic node
							childBody, err = parseNode(bodyChildNode, content)
						}

						if err != nil {
							return nil, fmt.Errorf("error parsing block child: %w", err)
						}

						if childBody != nil {
							block.Children = append(block.Children, childBody)
						}

						if !bodyCursor.GoToNextSibling() {
							break
						}
					}
				}
			}
		}
	} else {
		return nil, fmt.Errorf("block at line %d has no children", node.StartPoint().Row+1)
	}

	return block, nil
}

// parseAttribute converts a tree-sitter attribute node to our Attribute type
func parseAttribute(node *sitter.Node, content []byte) (types.Body, error) {
	if node == nil {
		return nil, fmt.Errorf("nil node passed to parseAttribute")
	}

	// Get range information
	nodeRange := sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}

	// Initialize the attribute with empty values
	attr := &types.Attribute{
		Name:  "",
		Range: nodeRange,
	}

	// Extract attribute comments and inline comments
	blockComment, inlineComment := findAssociatedComments(node, content)
	if blockComment != "" {
		attr.BlockComment = blockComment
	}
	attr.InlineComment = inlineComment

	// Create a tree-sitter cursor
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		// First child should be the identifier (attribute name)
		nameNode := cursor.CurrentNode()
		if nameNode.Type() == "identifier" {
			attr.Name = string(content[nameNode.StartByte():nameNode.EndByte()])
		} else {
			return nil, fmt.Errorf("attribute at line %d doesn't start with an identifier", node.StartPoint().Row+1)
		}

		// Move to the value
		if cursor.GoToNextSibling() {
			// Skip the equals token
			if cursor.CurrentNode().Type() == "=" {
				cursor.GoToNextSibling()
			}
			valueNode := cursor.CurrentNode()
			if valueNode != nil {
				value, err := parseExpression(valueNode, content)
				if err != nil {
					return nil, fmt.Errorf("error parsing attribute value: %w", err)
				}
				attr.Value = value
			}
		}
	} else {
		return nil, fmt.Errorf("attribute at line %d has no children", node.StartPoint().Row+1)
	}

	return attr, nil
}

// Helper function to create an expression range from a node
func exprRange(node *sitter.Node) sitter.Range {
	return sitter.Range{
		StartPoint: node.StartPoint(),
		EndPoint:   node.EndPoint(),
		StartByte:  node.StartByte(),
		EndByte:    node.EndByte(),
	}
}

// findAssociatedComments extracts block and inline comments associated with a node
func findAssociatedComments(node *sitter.Node, content []byte) (string, string) {
	var blockComment, inlineComment string

	// Look for comments before the node
	if node.PrevSibling() != nil && node.PrevSibling().Type() == "comment" {
		blockComment = string(content[node.PrevSibling().StartByte():node.PrevSibling().EndByte()])
	}

	// Look for comments after the node on the same line
	if node.NextSibling() != nil && node.NextSibling().Type() == "comment" {
		inlineComment = string(content[node.NextSibling().StartByte():node.NextSibling().EndByte()])
	}

	return blockComment, inlineComment
}

// parseExpression converts a tree-sitter expression node to our Expression type
func parseExpression(node *sitter.Node, content []byte) (types.Expression, error) {
	if debugEnabled {
		log.Printf("Parsing expression node type: %s, text: %s, child count: %d", node.Type(), node.Content(content), node.NamedChildCount())
		for i := uint32(0); i < node.NamedChildCount(); i++ {
			child := node.NamedChild(int(i))
			log.Printf("  Child %d: type=%s, text=%s", i, child.Type(), child.Content(content))
		}
	}

	switch node.Type() {
	case "expression":
		if node.NamedChildCount() > 0 {
			return parseExpression(node.NamedChild(0), content)
		}
		return nil, fmt.Errorf("empty expression node")
	case "operation":
		if node.NamedChildCount() > 0 {
			return parseExpression(node.NamedChild(0), content)
		}
		return nil, fmt.Errorf("empty operation node")
	case "binary_operation":
		if node.NamedChildCount() < 3 {
			return nil, fmt.Errorf("binary operation must have at least 3 children")
		}
		left, err := parseExpression(node.NamedChild(0), content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse left operand: %w", err)
		}
		operator := node.NamedChild(1).Content(content)
		right, err := parseExpression(node.NamedChild(2), content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse right operand: %w", err)
		}
		return &types.BinaryExpr{
			Left:      left,
			Operator:  operator,
			Right:     right,
			ExprRange: exprRange(node),
		}, nil
	case "string_lit":
		return parseStringLiteral(node, content)
	case "literal_value":
		return parseLiteralValue(node, content)
	case "collection_value":
		return parseCollectionValue(node, content)
	case "variable_expr":
		return parseVariableExpr(node, content)
	case "get_attr":
		return parseGetAttr(node, content)
	case "function_call":
		return parseFunctionCall(node, content)
	case "template_expr":
		return parseTemplateExpr(node, content)
	case "binary_expr":
		return parseBinaryExpr(node, content)
	case "paren_expr":
		return parseParenExpr(node, content)
	case "for_expr":
		return parseForExpr(node, content)
	case "conditional_expr":
		return parseConditionalExpr(node, content)
	case "conditional":
		if node.NamedChildCount() < 3 {
			return nil, fmt.Errorf("conditional must have at least 3 children")
		}
		condition, err := parseExpression(node.NamedChild(0), content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse condition: %w", err)
		}
		trueExpr, err := parseExpression(node.NamedChild(1), content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse true expression: %w", err)
		}
		falseExpr, err := parseExpression(node.NamedChild(2), content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse false expression: %w", err)
		}
		return &types.ConditionalExpr{
			Condition:  condition,
			TrueExpr:   trueExpr,
			FalseExpr:  falseExpr,
			ExprRange:  exprRange(node),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported expression type: %s", node.Type())
	}
}

func parseBinaryExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 3 {
		return nil, fmt.Errorf("binary expression must have at least 3 children")
	}

	left, err := parseExpression(node.NamedChild(0), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse left operand: %w", err)
	}

	operator := node.NamedChild(1).Content(content)

	right, err := parseExpression(node.NamedChild(2), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse right operand: %w", err)
	}

	return &types.BinaryExpr{
		Left:      left,
		Operator:  operator,
		Right:     right,
		ExprRange: exprRange(node),
	}, nil
}

func parseParenExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 1 {
		return nil, fmt.Errorf("parenthesized expression must have at least 1 child")
	}

	expr, err := parseExpression(node.NamedChild(0), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inner expression: %w", err)
	}

	return &types.ParenExpr{
		Expression: expr,
		ExprRange: exprRange(node),
	}, nil
}

func parseForExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 1 {
		return nil, fmt.Errorf("for expression must have at least 1 child")
	}

	// Find the for_tuple_expr or for_object_expr node
	var forNode *sitter.Node
	for i := uint32(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(int(i))
		if child.Type() == "for_tuple_expr" || child.Type() == "for_object_expr" {
			forNode = child
			break
		}
	}

	if forNode == nil {
		return nil, fmt.Errorf("for expression missing tuple or object expression")
	}

	var keyVar, valueVar string
	var collection, condition types.Expression
	var thenKeyExpr, thenValueExpr types.Expression

	// Parse the for expression parts
	for i := uint32(0); i < forNode.NamedChildCount(); i++ {
		child := forNode.NamedChild(int(i))
		switch child.Type() {
		case "for_intro":
			// Parse key and value variables
			for j := uint32(0); j < child.NamedChildCount(); j++ {
				introChild := child.NamedChild(int(j))
				if introChild.Type() == "identifier" {
					if valueVar == "" {
						valueVar = introChild.Content(content)
					} else {
						keyVar = valueVar
						valueVar = introChild.Content(content)
					}
				}
			}
		case "expression":
			// Parse collection and result expressions
			if collection == nil {
				var err error
				collection, err = parseExpression(child, content)
				if err != nil {
					return nil, fmt.Errorf("failed to parse collection: %w", err)
				}
			} else {
				// This is the result expression
				if forNode.Type() == "for_tuple_expr" {
					var err error
					thenValueExpr, err = parseExpression(child, content)
					if err != nil {
						return nil, fmt.Errorf("failed to parse then value expression: %w", err)
					}
				} else {
					// For object expressions, we need both key and value expressions
					if child.NamedChildCount() >= 2 {
						var err error
						thenKeyExpr, err = parseExpression(child.NamedChild(0), content)
						if err != nil {
							return nil, fmt.Errorf("failed to parse then key expression: %w", err)
						}
						thenValueExpr, err = parseExpression(child.NamedChild(1), content)
						if err != nil {
							return nil, fmt.Errorf("failed to parse then value expression: %w", err)
						}
					}
				}
			}
		case "binary_operation":
			// This might be the condition
			var err error
			condition, err = parseBinaryExpr(child, content)
			if err != nil {
				return nil, fmt.Errorf("failed to parse condition: %w", err)
			}
		}
	}

	forExpr := &types.ForExpr{
		KeyVar:       keyVar,
		ValueVar:     valueVar,
		Collection:   collection,
		Condition:    condition,
		ThenKeyExpr:  thenKeyExpr,
		ThenValueExpr: thenValueExpr,
		ExprRange:    exprRange(node),
	}

	// If this is a tuple for expression (array context), wrap it in an ArrayExpr
	if forNode.Type() == "for_tuple_expr" {
		return &types.ArrayExpr{
			Items:     []types.Expression{forExpr},
			ExprRange: exprRange(node),
		}, nil
	}

	// For object expressions, return the ForExpr directly
	return forExpr, nil
}

func parseConditionalExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 3 {
		return nil, fmt.Errorf("conditional expression must have at least 3 children")
	}

	condition, err := parseExpression(node.NamedChild(0), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse condition: %w", err)
	}

	trueExpr, err := parseExpression(node.NamedChild(1), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse true expression: %w", err)
	}

	falseExpr, err := parseExpression(node.NamedChild(2), content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse false expression: %w", err)
	}

	return &types.ConditionalExpr{
		Condition:  condition,
		TrueExpr:   trueExpr,
		FalseExpr:  falseExpr,
		ExprRange:  exprRange(node),
	}, nil
}

func parseFunctionCall(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 2 {
		return nil, fmt.Errorf("function call must have at least 2 children")
	}

	nameNode := node.NamedChild(0)
	if nameNode == nil || nameNode.Type() != "identifier" {
		return nil, fmt.Errorf("function call must have an identifier as its first child")
	}

	args := make([]types.Expression, 0)
	for i := uint32(1); i < node.NamedChildCount(); i++ {
		arg := node.NamedChild(int(i))
		if expr, err := parseExpression(arg, content); err == nil {
			args = append(args, expr)
		}
	}

	return &types.FunctionCallExpr{
		Name:      nameNode.Content(content),
		Args:      args,
		ExprRange: exprRange(node),
	}, nil
}

func parseStringLiteral(node *sitter.Node, content []byte) (types.Expression, error) {
	text := node.Content(content)
	// Strip quotes from string literals
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}
	return &types.LiteralValue{
		Value:     text,
		ValueType: "string",
		ExprRange: exprRange(node),
	}, nil
}

func parseLiteralValue(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 1 {
		return nil, fmt.Errorf("literal value node has no children")
	}

	firstChild := node.NamedChild(0)
	nodeType := firstChild.Type()
	nodeContent := firstChild.Content(content)

	switch nodeType {
	case "string_lit":
		// Strip quotes from string literals
		value := strings.Trim(nodeContent, "\"")
		return &types.LiteralValue{
			Value:     value,
			ValueType: "string",
			ExprRange: exprRange(node),
		}, nil
	case "bool_lit":
		value := nodeContent == "true"
		return &types.LiteralValue{
			Value:     value,
			ValueType: "bool",
			ExprRange: exprRange(node),
		}, nil
	case "null_lit":
		return &types.LiteralValue{
			Value:     nil,
			ValueType: "null",
			ExprRange: exprRange(node),
		}, nil
	default:
		// Try to parse as number, if fails return as string
		if val, err := strconv.ParseFloat(nodeContent, 64); err == nil {
			return &types.LiteralValue{
				Value:     val,
				ValueType: "number",
				ExprRange: exprRange(node),
			}, nil
		}
		return &types.LiteralValue{
			Value:     nodeContent,
			ValueType: "string",
			ExprRange: exprRange(node),
		}, nil
	}
}

func parseCollectionValue(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 1 {
		return nil, fmt.Errorf("collection value node has no children")
	}

	firstChild := node.NamedChild(0)
	nodeType := firstChild.Type()

	switch nodeType {
	case "[", "tuple":
		items := make([]types.Expression, 0)
		for i := uint32(0); i < firstChild.NamedChildCount(); i++ {
			child := firstChild.NamedChild(int(i))
			if expr, err := parseExpression(child, content); err == nil {
				items = append(items, expr)
			}
		}
		return &types.ArrayExpr{
			Items:     items,
			ExprRange: exprRange(node),
		}, nil
	case "{", "object":
		items := make([]types.ObjectItem, 0)
		for i := uint32(0); i < firstChild.NamedChildCount(); i++ {
			child := firstChild.NamedChild(int(i))
			if child.Type() == "object_elem" {
				key := child.NamedChild(0)
				value := child.NamedChild(1)
				if key != nil && value != nil {
					keyExpr, err := parseExpression(key, content)
					if err != nil {
						continue
					}
					valueExpr, err := parseExpression(value, content)
					if err != nil {
						continue
					}
					items = append(items, types.ObjectItem{
						Key:   keyExpr,
						Value: valueExpr,
					})
				}
			}
		}
		return &types.ObjectExpr{
			Items:     items,
			ExprRange: exprRange(node),
		}, nil
	default:
		// If not a collection, treat as a literal value
		return parseLiteralValue(node, content)
	}
}

func parseVariableExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.ChildCount() == 0 {
		return nil, fmt.Errorf("variable expression has no children")
	}

	// Get the identifier from the first child
	firstChild := node.Child(0)
	if firstChild.Type() != "identifier" {
		return nil, fmt.Errorf("first child of variable expression is not an identifier")
	}

	// Initialize parts with the identifier
	parts := []string{firstChild.Content(content)}

	// Process siblings (get_attr and index nodes)
	current := node.NextNamedSibling()
	for current != nil {
		switch current.Type() {
		case "get_attr":
			// Skip the leading dot in attribute name
			attrName := strings.TrimPrefix(current.Content(content), ".")
			parts = append(parts, attrName)
		case "index_expr":
			// Handle index expressions
			if current.ChildCount() > 0 {
				keyNode := current.Child(0)
				if keyNode.Type() == "string_lit" {
					// Strip quotes from string literals
					str := strings.Trim(keyNode.Content(content), "\"")
					parts = append(parts, str)
				}
			}
		}
		current = current.NextNamedSibling()
	}

	// If there are multiple parts, this is a traversal expression
	if len(parts) > 1 {
		traversal := make([]types.TraversalElem, len(parts)-1)
		for i, part := range parts[1:] {
			traversal[i] = types.TraversalElem{
				Name: part,
			}
		}
		return &types.RelativeTraversalExpr{
			Source:    nil,
			Traversal: traversal,
			ExprRange: exprRange(node),
		}, nil
	}

	// Otherwise, it's a simple reference
	return &types.ReferenceExpr{
		Parts:     parts,
		ExprRange: exprRange(node),
	}, nil
}

func parseGetAttr(node *sitter.Node, content []byte) (types.Expression, error) {
	// Skip the leading dot in attribute name
	attrName := strings.TrimPrefix(node.Content(content), ".")

	// Initialize parts with the attribute name
	parts := []string{attrName}

	// Process the source expression
	parent := node.Parent()
	if parent != nil && parent.Type() == "variable_expr" {
		// If the source is a variable expression, parse it and combine parts
		sourceExpr, err := parseVariableExpr(parent, content)
		if err != nil {
			return nil, err
		}

		if refExpr, ok := sourceExpr.(*types.ReferenceExpr); ok {
			// Combine the source parts with the attribute name
			parts = append(refExpr.Parts, attrName)

			// If there are multiple parts, this is a traversal expression
			if len(parts) > 1 {
				traversal := make([]types.TraversalElem, len(parts)-1)
				for i, part := range parts[1:] {
					traversal[i] = types.TraversalElem{
						Name: part,
					}
				}
				return &types.RelativeTraversalExpr{
					Source:    nil,
					Traversal: traversal,
					ExprRange: exprRange(node),
				}, nil
			}
		}
	}

	// If we can't process the source expression, just return a reference with the attribute name
	return &types.ReferenceExpr{
		Parts:     parts,
		ExprRange: exprRange(node),
	}, nil
}

func parseTemplateExpr(node *sitter.Node, content []byte) (types.Expression, error) {
	if node.NamedChildCount() < 1 {
		return nil, fmt.Errorf("template expression node has no children")
	}

	quotedTemplate := node.NamedChild(0)
	if quotedTemplate.Type() != "quoted_template" {
		return nil, fmt.Errorf("expected quoted_template, got %s", quotedTemplate.Type())
	}

	parts := make([]types.Expression, 0)
	for i := uint32(0); i < quotedTemplate.NamedChildCount(); i++ {
		child := quotedTemplate.NamedChild(int(i))
		switch child.Type() {
		case "template_literal":
			literal := &types.LiteralValue{
				Value:     child.Content(content),
				ValueType: "string",
				ExprRange: exprRange(child),
			}
			parts = append(parts, literal)
		case "template_interpolation":
			if child.NamedChildCount() > 0 {
				expr, err := parseExpression(child.NamedChild(0), content)
				if err == nil {
					parts = append(parts, expr)
				}
			}
		}
	}

	return &types.TemplateExpr{
		Parts:     parts,
		ExprRange: exprRange(node),
	}, nil
}
