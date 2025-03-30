package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/vahid-haghighat/terralint/parser"

	"github.com/vahid-haghighat/terralint/parser/types"
)

func Check(filePath string) error {
	extension := filepath.Ext(filePath)

	if extension != ".tf" && extension != ".tfvars" {
		return nil
	}

	original, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse the Terraform file to get the AST
	root, err := parser.ParseTerraformFile(filePath)
	if err != nil {
		return err
	}
	root = root

	// Uncomment for debugging
	printType(root, 0)

	formattedBytes, err := getFormattedContent(filePath)
	if err != nil {
		return err
	}

	return compare(original, formattedBytes)
}

func printType(v interface{}, indent int) {
	indentStr := strings.Repeat("    ", indent)
	nextIndentStr := strings.Repeat("    ", indent+1)

	if v == nil {
		fmt.Print("nil")
		return
	}

	switch v := v.(type) {
	case *types.LiteralValue:
		fmt.Printf("&types.LiteralValue{")
		if v.Value != nil {
			fmt.Printf("\n%sValue: ", nextIndentStr)
			switch val := v.Value.(type) {
			case string:
				fmt.Printf("%q", val)
			case bool, int, float64:
				fmt.Printf("%v", val)
			default:
				fmt.Printf("%v", val)
			}
			fmt.Printf(",\n")
		}
		if v.ValueType != "" {
			fmt.Printf("%sValueType: %q,\n", nextIndentStr, v.ValueType)
		}
		fmt.Printf("%s}", indentStr)

	case *types.ObjectExpr:
		fmt.Printf("&types.ObjectExpr{\n")
		fmt.Printf("%sItems: []types.ObjectItem{\n", nextIndentStr)
		for _, item := range v.Items {
			fmt.Printf("%s{\n", strings.Repeat("    ", indent+2))

			if item.Key != nil {
				fmt.Printf("%sKey: ", strings.Repeat("    ", indent+3))
				printType(item.Key, indent+3)
				fmt.Printf(",\n")
			}

			if item.Value != nil {
				fmt.Printf("%sValue: ", strings.Repeat("    ", indent+3))
				printType(item.Value, indent+3)
				fmt.Printf(",\n")
			}

			if item.InlineComment != "" {
				fmt.Printf("%sInlineComment: %q,\n", strings.Repeat("    ", indent+3), item.InlineComment)
			}

			if item.BlockComment != "" {
				fmt.Printf("%sBlockComment: %q,\n", strings.Repeat("    ", indent+3), item.BlockComment)
			}

			fmt.Printf("%s},\n", strings.Repeat("    ", indent+2))
		}
		fmt.Printf("%s},\n", nextIndentStr)
		fmt.Printf("%s}", indentStr)

	case *types.ArrayExpr:
		fmt.Printf("&types.ArrayExpr{\n")
		fmt.Printf("%sItems: []types.Expression{\n", nextIndentStr)
		for _, item := range v.Items {
			fmt.Printf("%s", strings.Repeat("    ", indent+2))
			printType(item, indent+2)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s},\n", nextIndentStr)
		fmt.Printf("%s}", indentStr)

	case *types.ReferenceExpr:
		fmt.Printf("&types.ReferenceExpr{\n")
		fmt.Printf("%sParts: %#v,\n", nextIndentStr, v.Parts)
		fmt.Printf("%s}", indentStr)

	case *types.FunctionCallExpr:
		fmt.Printf("&types.FunctionCallExpr{\n")
		fmt.Printf("%sName: %q,\n", nextIndentStr, v.Name)
		if len(v.Args) > 0 {
			fmt.Printf("%sArgs: []types.Expression{\n", nextIndentStr)
			for _, arg := range v.Args {
				fmt.Printf("%s", nextIndentStr)
				printType(arg, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}
		fmt.Printf("%s}", indentStr)

	case *types.TemplateExpr:
		fmt.Printf("&types.TemplateExpr{\n")
		fmt.Printf("%sParts: []types.Expression{\n", nextIndentStr)
		for _, part := range v.Parts {
			fmt.Printf("%s", strings.Repeat("    ", indent+2))
			printType(part, indent+2)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s},\n", nextIndentStr)
		fmt.Printf("%s}", indentStr)

	case *types.ConditionalExpr:
		fmt.Printf("&types.ConditionalExpr{\n")

		if v.Condition != nil {
			fmt.Printf("%sCondition: ", nextIndentStr)
			printType(v.Condition, indent+1)
			fmt.Printf(",\n")
		}

		if v.TrueExpr != nil {
			fmt.Printf("%sTrueExpr: ", nextIndentStr)
			printType(v.TrueExpr, indent+1)
			fmt.Printf(",\n")
		}

		if v.FalseExpr != nil {
			fmt.Printf("%sFalseExpr: ", nextIndentStr)
			printType(v.FalseExpr, indent+1)
			fmt.Printf(",\n")
		}

		fmt.Printf("%s}", indentStr)

	case *types.BinaryExpr:
		fmt.Printf("&types.BinaryExpr{\n")

		if v.Left != nil {
			fmt.Printf("%sLeft: ", nextIndentStr)
			printType(v.Left, indent+1)
			fmt.Printf(",\n")
		}

		if v.Operator != "" {
			fmt.Printf("%sOperator: %q,\n", nextIndentStr, v.Operator)
		}

		if v.Right != nil {
			fmt.Printf("%sRight: ", nextIndentStr)
			printType(v.Right, indent+1)
			fmt.Printf(",\n")
		}

		fmt.Printf("%s}", indentStr)

	case *types.ForArrayExpr, *types.ForMapExpr:
		// Handle for expressions
		var collection, condition types.Expression

		switch v := v.(type) {
		case *types.ForArrayExpr:
			collection = v.Collection
			condition = v.Condition
			fmt.Printf("&types.ForArrayExpr{\n")
			fmt.Printf("%sValueVar: %q,\n", nextIndentStr, v.ValueVar)
			if v.KeyVar != "" {
				fmt.Printf("%sKeyVar: %q,\n", nextIndentStr, v.KeyVar)
			}
		case *types.ForMapExpr:
			collection = v.Collection
			condition = v.Condition
			fmt.Printf("&types.ForMapExpr{\n")
			fmt.Printf("%sValueVar: %q,\n", nextIndentStr, v.ValueVar)
			if v.KeyVar != "" {
				fmt.Printf("%sKeyVar: %q,\n", nextIndentStr, v.KeyVar)
			}
		}

		// Print collection
		if collection != nil {
			fmt.Printf("%sCollection: ", nextIndentStr)
			printType(collection, indent+1)
			fmt.Printf(",\n")
		}

		// Print condition if present
		if condition != nil {
			fmt.Printf("%sCondition: ", nextIndentStr)
			printType(condition, indent+1)
			fmt.Printf(",\n")
		}

		fmt.Printf("%s}", indentStr)

	case *types.SplatExpr:
		fmt.Printf("&types.SplatExpr{\n")
		if v.Source != nil {
			fmt.Printf("%sSource: ", nextIndentStr)
			printType(v.Source, indent+1)
			fmt.Printf(",\n")
		}
		if v.Each != nil {
			fmt.Printf("%sEach: ", nextIndentStr)
			printType(v.Each, indent+1)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s}", indentStr)

	case *types.HeredocExpr:
		fmt.Printf("&types.HeredocExpr{\n")

		if v.Marker != "" {
			fmt.Printf("%sMarker: %q,\n", nextIndentStr, v.Marker)
		}

		if v.Content != "" {
			fmt.Printf("%sContent: %q,\n", nextIndentStr, v.Content)
		}

		if v.Indented {
			fmt.Printf("%sIndented: true,\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	case *types.IndexExpr:
		fmt.Printf("&types.IndexExpr{\n")

		if v.Collection != nil {
			fmt.Printf("%sCollection: ", nextIndentStr)
			printType(v.Collection, indent+1)
			fmt.Printf(",\n")
		}

		if v.Key != nil {
			fmt.Printf("%sKey: ", nextIndentStr)
			printType(v.Key, indent+1)
			fmt.Printf(",\n")
		}

		fmt.Printf("%s}", indentStr)

	case *types.TupleExpr:
		fmt.Printf("&types.TupleExpr{\n")
		fmt.Printf("%sExpressions: []types.Expression{\n", nextIndentStr)
		for _, expr := range v.Expressions {
			fmt.Printf("%s", strings.Repeat("    ", indent+2))
			printType(expr, indent+2)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s},\n", nextIndentStr)
		fmt.Printf("%s}", indentStr)

	case *types.UnaryExpr:
		fmt.Printf("&types.UnaryExpr{\n")

		if v.Operator != "" {
			fmt.Printf("%sOperator: %q,\n", nextIndentStr, v.Operator)
		}

		if v.Expr != nil {
			fmt.Printf("%sExpr: ", nextIndentStr)
			printType(v.Expr, indent+1)
			fmt.Printf(",\n")
		}

		fmt.Printf("%s}", indentStr)

	case *types.ParenExpr:
		fmt.Printf("&types.ParenExpr{\n")
		if v.Expression != nil {
			fmt.Printf("%sExpression: ", nextIndentStr)
			printType(v.Expression, indent+1)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s}", indentStr)

	case *types.RelativeTraversalExpr:
		fmt.Printf("&types.RelativeTraversalExpr{\n")

		if v.Source != nil {
			fmt.Printf("%sSource: ", nextIndentStr)
			printType(v.Source, indent+1)
			fmt.Printf(",\n")
		}

		if len(v.Traversal) > 0 {
			fmt.Printf("%sTraversal: []types.TraversalElem{\n", nextIndentStr)
			for _, elem := range v.Traversal {
				fmt.Printf("%s{\n", strings.Repeat("    ", indent+2))
				fmt.Printf("%sType: %q,\n", strings.Repeat("    ", indent+3), elem.Type)

				if elem.Name != "" {
					fmt.Printf("%sName: %q,\n", strings.Repeat("    ", indent+3), elem.Name)
				}

				if elem.Index != nil {
					fmt.Printf("%sIndex: ", strings.Repeat("    ", indent+3))
					printType(elem.Index, indent+3)
					fmt.Printf(",\n")
				}

				fmt.Printf("%s},\n", strings.Repeat("    ", indent+2))
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	case *types.TemplateForDirective:
		fmt.Printf("&types.TemplateForDirective{\n")

		if v.KeyVar != "" {
			fmt.Printf("%sKeyVar: %q,\n", nextIndentStr, v.KeyVar)
		}

		if v.ValueVar != "" {
			fmt.Printf("%sValueVar: %q,\n", nextIndentStr, v.ValueVar)
		}

		if v.CollExpr != nil {
			fmt.Printf("%sCollExpr: ", nextIndentStr)
			printType(v.CollExpr, indent+1)
			fmt.Printf(",\n")
		}

		if len(v.Content) > 0 {
			fmt.Printf("%sContent: []types.Expression{\n", nextIndentStr)
			for _, expr := range v.Content {
				fmt.Printf("%s", strings.Repeat("    ", indent+2))
				printType(expr, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	case *types.TemplateIfDirective:
		fmt.Printf("&types.TemplateIfDirective{\n")

		if v.Condition != nil {
			fmt.Printf("%sCondition: ", nextIndentStr)
			printType(v.Condition, indent+1)
			fmt.Printf(",\n")
		}

		if len(v.TrueExpr) > 0 {
			fmt.Printf("%sTrueExpr: []types.Expression{\n", nextIndentStr)
			for _, expr := range v.TrueExpr {
				fmt.Printf("%s", strings.Repeat("    ", indent+2))
				printType(expr, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		if len(v.FalseExpr) > 0 {
			fmt.Printf("%sFalseExpr: []types.Expression{\n", nextIndentStr)
			for _, expr := range v.FalseExpr {
				fmt.Printf("%s", strings.Repeat("    ", indent+2))
				printType(expr, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	// Handle Body types as well
	case *types.Block:
		fmt.Printf("&types.Block{\n")

		if v.Type != "" {
			fmt.Printf("%sType: %q,\n", nextIndentStr, v.Type)
		}

		if len(v.Labels) > 0 {
			fmt.Printf("%sLabels: %#v,\n", nextIndentStr, v.Labels)
		}

		if v.InlineComment != "" {
			fmt.Printf("%sInlineComment: %q,\n", nextIndentStr, v.InlineComment)
		}

		if v.BlockComment != "" {
			fmt.Printf("%sBlockComment: %q,\n", nextIndentStr, v.BlockComment)
		}

		if len(v.Children) > 0 {
			fmt.Printf("%sChildren: []types.Body{\n", nextIndentStr)
			for _, child := range v.Children {
				fmt.Printf("%s", strings.Repeat("    ", indent+2))
				printType(child, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	case *types.Attribute:
		fmt.Printf("&types.Attribute{\n")

		if v.Name != "" {
			fmt.Printf("%sName: %q,\n", nextIndentStr, v.Name)
		}

		if v.Value != nil {
			fmt.Printf("%sValue: ", nextIndentStr)
			printType(v.Value, indent+1)
			fmt.Printf(",\n")
		}

		if v.InlineComment != "" {
			fmt.Printf("%sInlineComment: %q,\n", nextIndentStr, v.InlineComment)
		}

		if v.BlockComment != "" {
			fmt.Printf("%sBlockComment: %q,\n", nextIndentStr, v.BlockComment)
		}

		fmt.Printf("%s}", indentStr)

	case *types.Root:
		fmt.Printf("&types.Root{\n")

		if len(v.Children) > 0 {
			fmt.Printf("%sChildren: []types.Body{\n", nextIndentStr)
			for _, child := range v.Children {
				fmt.Printf("%s", strings.Repeat("    ", indent+2))
				printType(child, indent+2)
				fmt.Printf(",\n")
			}
			fmt.Printf("%s},\n", nextIndentStr)
		}

		fmt.Printf("%s}", indentStr)

	case *types.FormatDirective:
		fmt.Printf("&types.FormatDirective{\n")

		if v.DirectiveType != "" {
			fmt.Printf("%sDirectiveType: %q,\n", nextIndentStr, v.DirectiveType)
		}

		if len(v.Parameters) > 0 {
			fmt.Printf("%sParameters: %#v,\n", nextIndentStr, v.Parameters)
		}

		fmt.Printf("%s}", indentStr)

	default:
		// For any other types, use reflection as a fallback
		fmt.Printf("%T(%+v)", v, v)
	}
}

// printASTHumanReadable prints the AST in a readable format
func printASTHumanReadable(node interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *types.Root:
		fmt.Printf("%sRoot with %d children\n", indentStr, len(n.Children))
		for _, child := range n.Children {
			if child != nil {
				printASTHumanReadable(child, indent+1)
			}
		}
	case *types.Block:
		fmt.Printf("%sBlock: Type=%s, Labels=%v, %d children\n",
			indentStr, n.Type, n.Labels, len(n.Children))
		if n.BlockComment != "" {
			fmt.Printf("%s  BlockComment: %s\n", indentStr, n.BlockComment)
		}
		if n.InlineComment != "" {
			fmt.Printf("%s  InlineComment: %s\n", indentStr, n.InlineComment)
		}
		for _, child := range n.Children {
			printASTHumanReadable(child, indent+1)
		}
	case *types.Attribute:
		fmt.Printf("%sAttribute: Name=%s, Value=%v\n",
			indentStr, n.Name, getExpressionSummary(n.Value))
		if n.BlockComment != "" {
			fmt.Printf("%s  BlockComment: %s\n", indentStr, n.BlockComment)
		}
		if n.InlineComment != "" {
			fmt.Printf("%s  InlineComment: %s\n", indentStr, n.InlineComment)
		}
	default:
		fmt.Printf("%sUnknown node type: %T\n", indentStr, n)
	}
}

// getExpressionSummary returns a summary of an expression
func getExpressionSummary(expr types.Expression) string {
	if expr == nil {
		return "nil"
	}

	switch e := expr.(type) {
	case *types.LiteralValue:
		return fmt.Sprintf("Literal(%v, %s)", e.Value, e.ValueType)
	case *types.ObjectExpr:
		summary := fmt.Sprintf("Object with %d items", len(e.Items))
		if len(e.Items) > 0 {
			summary += " ["
			for i, item := range e.Items {
				if i > 0 {
					summary += ", "
				}

				// Check if we have a nested reference pattern (Reference(a): Reference(b))
				// and flatten it to Reference(a.b) for better readability
				if refKey, isRefKey := item.Key.(*types.ReferenceExpr); isRefKey {
					if refValue, isRefValue := item.Value.(*types.ReferenceExpr); isRefValue {
						// This is a simple case like github.here = github.here
						summary += fmt.Sprintf("Reference(%s): Reference(%s)",
							strings.Join(refKey.Parts, "."),
							strings.Join(refValue.Parts, "."))
					} else if objValue, isObjValue := item.Value.(*types.ObjectExpr); isObjValue && len(objValue.Items) == 1 {
						// This is a nested case like a = { n = 1 }
						// Try to flatten it to a.n = 1
						nestedItem := objValue.Items[0]
						if nestedKey, isNestedKeyRef := nestedItem.Key.(*types.ReferenceExpr); isNestedKeyRef {
							// Create a flattened reference
							flattenedParts := append(refKey.Parts, nestedKey.Parts...)
							summary += fmt.Sprintf("Reference(%s): %s",
								strings.Join(flattenedParts, "."),
								getExpressionSummary(nestedItem.Value))
						} else {
							// Fall back to normal format
							summary += fmt.Sprintf("%s: %s",
								getExpressionSummary(item.Key),
								getExpressionSummary(item.Value))
						}
					} else {
						// Normal key-value pair
						summary += fmt.Sprintf("%s: %s",
							getExpressionSummary(item.Key),
							getExpressionSummary(item.Value))
					}
				} else {
					// Normal key-value pair
					summary += fmt.Sprintf("%s: %s",
						getExpressionSummary(item.Key),
						getExpressionSummary(item.Value))
				}

				if item.BlockComment != "" || item.InlineComment != "" {
					summary += " (has comments)"
				}
			}
			summary += "]"
		}
		return summary
	case *types.ArrayExpr:
		summary := fmt.Sprintf("Array with %d items", len(e.Items))
		if len(e.Items) > 0 {
			summary += " ["
			for i, item := range e.Items {
				if i > 0 {
					summary += ", "
				}
				summary += getExpressionSummary(item)
			}
			summary += "]"
		}
		return summary
	case *types.TupleExpr:
		summary := fmt.Sprintf("Tuple with %d items", len(e.Expressions))
		if len(e.Expressions) > 0 {
			summary += " ["
			for i, item := range e.Expressions {
				if i > 0 {
					summary += ", "
				}
				summary += getExpressionSummary(item)
			}
			summary += "]"
		}
		return summary
	case *types.ReferenceExpr:
		return fmt.Sprintf("Reference(%s)", strings.Join(e.Parts, "."))
	case *types.FunctionCallExpr:
		return fmt.Sprintf("FunctionCall(%s, %d args)", e.Name, len(e.Args))
	case *types.ForArrayExpr, *types.ForMapExpr:
		// Handle for expressions
		var collection, condition types.Expression

		switch e := e.(type) {
		case *types.ForArrayExpr:
			collection = e.Collection
			condition = e.Condition
			return fmt.Sprintf("for %s in %s%s",
				e.ValueVar,
				getExpressionSummary(collection),
				conditionSummary(condition))
		case *types.ForMapExpr:
			collection = e.Collection
			condition = e.Condition
			return fmt.Sprintf("for %s, %s in %s%s",
				e.KeyVar,
				e.ValueVar,
				getExpressionSummary(collection),
				conditionSummary(condition))
		}
		return fmt.Sprintf("%T", e)
	default:
		return fmt.Sprintf("%T", e)
	}
}

func conditionSummary(condition types.Expression) string {
	if condition != nil {
		return fmt.Sprintf(" if %s", getExpressionSummary(condition))
	}
	return ""
}

func compare(original []byte, formatted []byte) error {
	originalHash, formattedHash := generateHash(original, formatted)

	if formattedHash == originalHash {
		return nil
	}

	return compareContent(original, formatted)
}

func generateHash(original []byte, formatted []byte) (string, string) {
	hasher := sha1.New()
	hasher.Write(original)

	originalHash := hex.EncodeToString(hasher.Sum(nil))

	hasher.Reset()

	hasher.Write(formatted)

	formattedHash := hex.EncodeToString(hasher.Sum(nil))
	return originalHash, formattedHash
}

func compareContent(original []byte, formatted []byte) error {
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = time.Hour
	src := string(original)
	dst := string(formatted)

	wSrc, wDst, warray := dmp.DiffLinesToRunes(src, dst)
	diffs := dmp.DiffMainRunes(wSrc, wDst, false)
	diffs = dmp.DiffCharsToLines(diffs, warray)

	var notEquals []diffmatchpatch.Diff
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			notEquals = append(notEquals, diff)
		}
	}

	if notEquals == nil || len(notEquals) == 0 {
		return nil
	}

	var errorText strings.Builder
	errorText.WriteString("\n")
	errorText.WriteString(dmp.DiffPrettyText(diffs))
	return errors.New(errorText.String())
}
