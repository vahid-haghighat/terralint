package parser

import (
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	. "github.com/hashicorp/hcl/v2/hclsyntax"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vahid-haghighat/terralint/parser/types"
	"github.com/zclconf/go-cty/cty"
)

// ParseTerraformFile reads a Terraform file and parses it into an AST
func ParseTerraformFile(filePath string) (*types.Root, error) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a new HCL parser
	parser := hclparse.NewParser()

	// Parse the file
	file, diags := parser.ParseHCL(content, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	// Create root node
	root := &types.Root{
		Children: make([]types.Body, 0),
	}

	// Get the root body
	body := file.Body.(*hclsyntax.Body)

	// Process attributes and blocks in order by their source position
	var items []hclsyntax.Node
	for name, attr := range body.Attributes {
		items = append(items, &hclsyntax.Attribute{
			Name:        name,
			Expr:        attr.Expr,
			SrcRange:    attr.SrcRange,
			NameRange:   attr.NameRange,
			EqualsRange: attr.EqualsRange,
		})
	}
	for _, block := range body.Blocks {
		items = append(items, block)
	}

	// Sort items by their source position
	sort.Slice(items, func(i, j int) bool {
		return items[i].Range().Start.Byte < items[j].Range().Start.Byte
	})

	// Process items in order
	for _, item := range items {
		switch item := item.(type) {
		case *hclsyntax.Attribute:
			attribute, err := convertAttribute(item.Name, item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert attribute %s: %w", item.Name, err)
			}
			root.Children = append(root.Children, attribute)
		case *hclsyntax.Block:
			block, err := convertBlock(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert block: %w", err)
			}
			root.Children = append(root.Children, block)
		}
	}

	return root, nil
}

func convertBlock(block *hclsyntax.Block) (*types.Block, error) {
	b := &types.Block{
		Type:     block.Type,
		Labels:   block.Labels,
		Range:    block.Range(),
		Children: make([]types.Body, 0),
	}

	// Process attributes and blocks in order by their source position
	var items []hclsyntax.Node
	for name, attr := range block.Body.Attributes {
		items = append(items, &hclsyntax.Attribute{
			Name:        name,
			Expr:        attr.Expr,
			SrcRange:    attr.SrcRange,
			NameRange:   attr.NameRange,
			EqualsRange: attr.EqualsRange,
		})
	}
	for _, block := range block.Body.Blocks {
		items = append(items, block)
	}

	// Sort items by their source position
	sort.Slice(items, func(i, j int) bool {
		return items[i].Range().Start.Byte < items[j].Range().Start.Byte
	})

	// Process items in order
	for _, item := range items {
		switch item := item.(type) {
		case *hclsyntax.Attribute:
			attribute, err := convertAttribute(item.Name, item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert attribute %s: %w", item.Name, err)
			}
			b.Children = append(b.Children, attribute)
		case *hclsyntax.Block:
			nested, err := convertBlock(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert nested block: %w", err)
			}
			b.Children = append(b.Children, nested)
		}
	}

	return b, nil
}

func convertAttribute(name string, attr *hclsyntax.Attribute) (*types.Attribute, error) {
	expr, err := convertExpression(attr.Expr)
	if err != nil {
		return nil, err
	}

	return &types.Attribute{
		Name:  name,
		Value: expr,
		Range: attr.Range(),
	}, nil
}

func convertExpression(expr hclsyntax.Expression) (types.Expression, error) {
	switch e := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		var value interface{}
		switch {
		case e.Val.Type() == cty.String:
			value = e.Val.AsString()
		case e.Val.Type() == cty.Bool:
			value = e.Val.True()
		case e.Val.Type() == cty.Number:
			bf := e.Val.AsBigFloat()
			if bf.IsInt() {
				i, _ := bf.Int64()
				value = i
			} else {
				f, _ := bf.Float64()
				value = f
			}
		default:
			value = e.Val.GoString()
		}
		return &types.LiteralValue{
			Value:     value,
			ValueType: e.Val.Type().FriendlyName(),
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.TemplateExpr:
		// Check if this is a single literal value
		if len(e.Parts) == 1 {
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				var value interface{}
				switch {
				case lit.Val.Type() == cty.String:
					value = lit.Val.AsString()
				case lit.Val.Type() == cty.Bool:
					value = lit.Val.True()
				case lit.Val.Type() == cty.Number:
					bf := lit.Val.AsBigFloat()
					if bf.IsInt() {
						i, _ := bf.Int64()
						value = i
					} else {
						f, _ := bf.Float64()
						value = f
					}
				default:
					value = lit.Val.GoString()
				}
				return &types.LiteralValue{
					Value:     value,
					ValueType: lit.Val.Type().FriendlyName(),
					ExprRange: e.Range(),
				}, nil
			}
		}

		// Check if this is a heredoc by looking at the source code
		srcRange := e.Range()
		if srcRange.Start.Line != srcRange.End.Line && len(e.Parts) > 0 {
			// This is a multi-line template, treat it as a heredoc
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return &types.HeredocExpr{
					Marker:    "EOT", // Default marker
					Content:   lit.Val.AsString(),
					Indented:  false,
					ExprRange: e.Range(),
				}, nil
			}
		}

		parts := make([]types.Expression, len(e.Parts))
		for i, part := range e.Parts {
			converted, err := convertExpression(part)
			if err != nil {
				return nil, err
			}
			parts[i] = converted
		}
		return &types.TemplateExpr{
			Parts:     parts,
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.ScopeTraversalExpr:
		// Check if this is a single part that's a quoted string
		if len(e.Traversal) == 1 {
			if root, ok := e.Traversal[0].(hcl.TraverseRoot); ok {
				if strings.HasPrefix(root.Name, "\"") && strings.HasSuffix(root.Name, "\"") {
					// This is a quoted string, treat it as a literal
					value := strings.Trim(root.Name, "\"")
					return &types.LiteralValue{
						Value:     value,
						ValueType: "string",
						ExprRange: e.Range(),
					}, nil
				}
			}
		}

		// Build parts for traversal
		parts := make([]string, len(e.Traversal))
		for i, traverser := range e.Traversal {
			switch t := traverser.(type) {
			case hcl.TraverseRoot:
				parts[i] = t.Name
			case hcl.TraverseAttr:
				parts[i] = t.Name
			case hcl.TraverseIndex:
				switch t.Key.Type() {
				case cty.Number:
					bf := t.Key.AsBigFloat()
					if i, acc := bf.Int64(); acc == big.Exact {
						parts[i] = fmt.Sprintf("[%d]", i)
					} else {
						parts[i] = fmt.Sprintf("[%s]", bf.String())
					}
				case cty.String:
					parts[i] = fmt.Sprintf("[%q]", t.Key.AsString())
				}
			}
		}

		// If any part is an index expression, return a relative traversal
		hasIndex := false
		for _, part := range parts {
			if strings.HasPrefix(part, "[") {
				hasIndex = true
				break
			}
		}

		if hasIndex {
			// Split parts into source and traversal
			source := &types.ReferenceExpr{
				Parts:     parts[:1],
				ExprRange: e.Range(),
			}
			traversal := make([]types.TraversalElem, len(parts)-1)
			for i, part := range parts[1:] {
				if strings.HasPrefix(part, "[") {
					traversal[i] = types.TraversalElem{
						Type: "index",
						Name: part,
					}
				} else {
					traversal[i] = types.TraversalElem{
						Type: "attr",
						Name: part,
					}
				}
			}
			return &types.RelativeTraversalExpr{
				Source:    source,
				Traversal: traversal,
				ExprRange: e.Range(),
			}, nil
		}

		// Otherwise, return a reference
		return &types.ReferenceExpr{
			Parts:     parts,
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.FunctionCallExpr:
		args := make([]types.Expression, len(e.Args))
		for i, arg := range e.Args {
			converted, err := convertExpression(arg)
			if err != nil {
				return nil, err
			}
			args[i] = converted
		}

		// Create a function call expression
		return &types.FunctionCallExpr{
			Name:      e.Name,
			Args:      args,
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.ObjectConsExpr:
		items := make([]types.ObjectItem, len(e.Items))
		for i, item := range e.Items {
			key, err := convertExpression(item.KeyExpr)
			if err != nil {
				return nil, err
			}
			value, err := convertExpression(item.ValueExpr)
			if err != nil {
				return nil, err
			}
			items[i] = types.ObjectItem{
				Key:   key,
				Value: value,
			}
		}
		return &types.ObjectExpr{
			Items: items,
		}, nil
	case *hclsyntax.ObjectConsKeyExpr:
		// For object keys, always create a ReferenceExpr
		if traversal, ok := e.Wrapped.(*hclsyntax.ScopeTraversalExpr); ok {
			parts := make([]string, len(traversal.Traversal))
			for i, traverser := range traversal.Traversal {
				switch t := traverser.(type) {
				case hcl.TraverseRoot:
					parts[i] = t.Name
				case hcl.TraverseAttr:
					parts[i] = t.Name
				}
			}
			return &types.ReferenceExpr{
				Parts: parts,
			}, nil
		}
		return convertExpression(e.Wrapped)
	case *hclsyntax.TupleConsExpr:
		items := make([]types.Expression, len(e.Exprs))
		for i, expr := range e.Exprs {
			converted, err := convertExpression(expr)
			if err != nil {
				return nil, err
			}
			items[i] = converted
		}
		return &types.ArrayExpr{
			Items: items,
		}, nil
	case *hclsyntax.BinaryOpExpr:
		left, err := convertExpression(e.LHS)
		if err != nil {
			return nil, err
		}
		right, err := convertExpression(e.RHS)
		if err != nil {
			return nil, err
		}
		// Map the operation to a string based on its token type
		var operator string
		switch e.Op.Type {
		case cty.Number:
			switch e.Op {
			case OpAdd:
				operator = "+"
			case OpSubtract:
				operator = "-"
			case OpMultiply:
				operator = "*"
			case OpDivide:
				operator = "/"
			case OpModulo:
				operator = "%"
			default:
				return nil, fmt.Errorf("unsupported numeric operator")
			}
		case cty.Bool:
			switch e.Op {
			case OpEqual:
				operator = "=="
			case OpNotEqual:
				operator = "!="
			case OpGreaterThan:
				operator = ">"
			case OpGreaterThanOrEqual:
				operator = ">="
			case OpLessThan:
				operator = "<"
			case OpLessThanOrEqual:
				operator = "<="
			case OpLogicalAnd:
				operator = "&&"
			case OpLogicalOr:
				operator = "||"
			default:
				return nil, fmt.Errorf("unsupported boolean operator")
			}
		case cty.DynamicPseudoType:
			// For dynamic types, we need to determine the operator based on the operation itself
			switch e.Op {
			case OpEqual:
				operator = "=="
			case OpNotEqual:
				operator = "!="
			case OpGreaterThan:
				operator = ">"
			case OpGreaterThanOrEqual:
				operator = ">="
			case OpLessThan:
				operator = "<"
			case OpLessThanOrEqual:
				operator = "<="
			case OpLogicalAnd:
				operator = "&&"
			case OpLogicalOr:
				operator = "||"
			case OpAdd:
				operator = "+"
			case OpSubtract:
				operator = "-"
			case OpMultiply:
				operator = "*"
			case OpDivide:
				operator = "/"
			case OpModulo:
				operator = "%"
			default:
				return nil, fmt.Errorf("unsupported dynamic operator")
			}
		default:
			return nil, fmt.Errorf("unsupported operator type: %s", e.Op.Type.FriendlyName())
		}
		return &types.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}, nil
	case *hclsyntax.UnaryOpExpr:
		expr, err := convertExpression(e.Val)
		if err != nil {
			return nil, err
		}
		// Map the operation to a string based on its type
		var operator string
		switch e.Op.Impl.Params()[0].Type {
		case cty.Number:
			operator = "-"
		case cty.Bool:
			operator = "!"
		default:
			return nil, fmt.Errorf("unsupported unary operator type: %s", e.Op.Impl.Params()[0].Type.FriendlyName())
		}
		return &types.UnaryExpr{
			Operator: operator,
			Expr:     expr,
		}, nil
	case *hclsyntax.ConditionalExpr:
		condition, err := convertExpression(e.Condition)
		if err != nil {
			return nil, err
		}
		trueResult, err := convertExpression(e.TrueResult)
		if err != nil {
			return nil, err
		}
		falseResult, err := convertExpression(e.FalseResult)
		if err != nil {
			return nil, err
		}
		return &types.ConditionalExpr{
			Condition: condition,
			TrueExpr:  trueResult,
			FalseExpr: falseResult,
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.ForExpr:
		collection, err := convertExpression(e.CollExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert collection: %w", err)
		}
		thenValue, err := convertExpression(e.ValExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value expression: %w", err)
		}
		var condition types.Expression
		if e.CondExpr != nil {
			condition, err = convertExpression(e.CondExpr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert condition: %w", err)
			}
		}
		if e.KeyExpr != nil {
			thenKey, err := convertExpression(e.KeyExpr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert key expression: %w", err)
			}
			return &types.ForMapExpr{
				KeyVar:        e.KeyVar,
				ValueVar:      e.ValVar,
				Collection:    collection,
				ThenKeyExpr:   thenKey,
				ThenValueExpr: thenValue,
				Condition:     condition,
			}, nil
		}
		return &types.ForArrayExpr{
			KeyVar:        e.KeyVar,
			ValueVar:      e.ValVar,
			Collection:    collection,
			ThenValueExpr: thenValue,
			Condition:     condition,
		}, nil
	case *hclsyntax.SplatExpr:
		source, err := convertExpression(e.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to convert splat source: %w", err)
		}
		each, err := convertExpression(e.Each)
		if err != nil {
			return nil, fmt.Errorf("failed to convert splat each: %w", err)
		}
		return &types.SplatExpr{
			Source: source,
			Each:   each,
		}, nil
	case *hclsyntax.IndexExpr:
		collection, err := convertExpression(e.Collection)
		if err != nil {
			return nil, err
		}
		key, err := convertExpression(e.Key)
		if err != nil {
			return nil, err
		}
		return &types.IndexExpr{
			Collection: collection,
			Key:        key,
			ExprRange:  e.Range(),
		}, nil
	case *hclsyntax.ParenthesesExpr:
		expression, err := convertExpression(e.Expression)
		if err != nil {
			return nil, err
		}
		return &types.ParenExpr{
			Expression: expression,
			ExprRange:  e.Range(),
		}, nil
	case *hclsyntax.TemplateJoinExpr:
		tuple, err := convertExpression(e.Tuple)
		if err != nil {
			return nil, err
		}
		return &types.TemplateExpr{
			Parts:     []types.Expression{tuple},
			ExprRange: e.Range(),
		}, nil
	case *hclsyntax.TemplateWrapExpr:
		wrapped, err := convertExpression(e.Wrapped)
		if err != nil {
			return nil, fmt.Errorf("failed to convert wrapped expression: %w", err)
		}
		return wrapped, nil
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func convertExpressions(exprs []hclsyntax.Expression) ([]types.Expression, error) {
	result := make([]types.Expression, len(exprs))
	for i, expr := range exprs {
		converted, err := convertExpression(expr)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func convertLiteralValue(val cty.Value, rng hcl.Range) (*types.LiteralValue, error) {
	var value interface{}
	var valueType string

	switch {
	case val.Type() == cty.String:
		value = val.AsString()
		valueType = "string"
	case val.Type() == cty.Number:
		bf := val.AsBigFloat()
		if f64, acc := bf.Float64(); acc == 0 {
			value = f64
			valueType = "number"
		} else {
			value = bf.String()
			valueType = "number"
		}
	case val.Type() == cty.Bool:
		value = val.True()
		valueType = "bool"
	case val.IsNull():
		value = nil
		valueType = "null"
	default:
		return nil, fmt.Errorf("unsupported literal value type: %s", val.Type().GoString())
	}

	return &types.LiteralValue{
		Value:     value,
		ValueType: valueType,
		ExprRange: rng,
	}, nil
}

func convertRange(r hcl.Range) sitter.Range {
	return sitter.Range{
		StartPoint: sitter.Point{
			Row:    uint32(r.Start.Line - 1),
			Column: uint32(r.Start.Column - 1),
		},
		EndPoint: sitter.Point{
			Row:    uint32(r.End.Line - 1),
			Column: uint32(r.End.Column - 1),
		},
	}
}

// getOperationSymbol returns the string representation of an operation
func getOperationSymbol(op *hclsyntax.Operation) string {
	switch op {
	// Arithmetic
	case hclsyntax.OpAdd:
		return "+"
	case hclsyntax.OpSubtract:
		return "-"
	case hclsyntax.OpMultiply:
		return "*"
	case hclsyntax.OpDivide:
		return "/"
	case hclsyntax.OpModulo:
		return "%"
	case hclsyntax.OpNegate:
		return "-"

	// Logical
	case hclsyntax.OpLogicalAnd:
		return "&&"
	case hclsyntax.OpLogicalOr:
		return "||"
	case hclsyntax.OpLogicalNot:
		return "!"

	// Comparison
	case hclsyntax.OpEqual:
		return "=="
	case hclsyntax.OpNotEqual:
		return "!="
	case hclsyntax.OpGreaterThan:
		return ">"
	case hclsyntax.OpGreaterThanOrEqual:
		return ">="
	case hclsyntax.OpLessThan:
		return "<"
	case hclsyntax.OpLessThanOrEqual:
		return "<="

	default:
		return "unknown"
	}
}

func convertTraversal(traversal []hcl.Traverser) []string {
	parts := make([]string, len(traversal))
	for i, traverser := range traversal {
		switch t := traverser.(type) {
		case hcl.TraverseRoot:
			parts[i] = t.Name
		case hcl.TraverseAttr:
			parts[i] = t.Name
		case hcl.TraverseIndex:
			// Convert the index value to a string
			switch {
			case t.Key.Type() == cty.Number:
				bf := t.Key.AsBigFloat()
				if bf.IsInt() {
					i, _ := bf.Int64()
					parts[i] = fmt.Sprintf("[%d]", i)
				} else {
					f, _ := bf.Float64()
					parts[i] = fmt.Sprintf("[%f]", f)
				}
			case t.Key.Type() == cty.String:
				parts[i] = fmt.Sprintf("[%q]", t.Key.AsString())
			default:
				parts[i] = fmt.Sprintf("[%s]", t.Key.GoString())
			}
		}
	}
	return parts
}

func convertTraversalParts(traversal hcl.Traversal) []string {
	parts := make([]string, len(traversal))
	for i, traverser := range traversal {
		switch t := traverser.(type) {
		case hcl.TraverseRoot:
			parts[i] = t.Name
		case hcl.TraverseAttr:
			parts[i] = t.Name
		case hcl.TraverseIndex:
			// Convert the index value to a string
			switch {
			case t.Key.Type() == cty.Number:
				bf := t.Key.AsBigFloat()
				if bf.IsInt() {
					i, _ := bf.Int64()
					parts[i] = fmt.Sprintf("[%d]", i)
				} else {
					f, _ := bf.Float64()
					parts[i] = fmt.Sprintf("[%f]", f)
				}
			case t.Key.Type() == cty.String:
				parts[i] = fmt.Sprintf("[%q]", t.Key.AsString())
			default:
				parts[i] = fmt.Sprintf("[%s]", t.Key.GoString())
			}
		}
	}
	return parts
}

// isReferenceExpr checks if a traversal should be treated as a reference expression
func isReferenceExpr(parts []string) bool {
	if len(parts) == 0 {
		return false
	}

	// If it's a single part, it's a reference if it's a simple identifier
	if len(parts) == 1 {
		// Check if it's a simple identifier (not an index)
		return !strings.HasPrefix(parts[0], "[")
	}

	// If it's a traversal with multiple parts, it's a reference if:
	// 1. The first part is a simple identifier (not an index)
	// 2. The subsequent parts are either simple identifiers or indices
	return !strings.HasPrefix(parts[0], "[")
}
