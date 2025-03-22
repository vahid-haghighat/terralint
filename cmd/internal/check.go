package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/vahid-haghighat/terralint/parser"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	fmt.Println("Parsed Terraform AST with Comments:")
	printAST(root, 0)

	formattedBytes, err := getFormattedContent(filePath)
	if err != nil {
		return err
	}

	return compare(original, formattedBytes)
}

// printAST prints the AST in a readable format
func printAST(node interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *types.Root:
		fmt.Printf("%sRoot with %d children\n", indentStr, len(n.Children))
		for _, child := range n.Children {
			if child != nil {
				printAST(child, indent+1)
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
			printAST(child, indent+1)
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
	case *types.DynamicBlock:
		fmt.Printf("%sDynamicBlock: Labels=%v, Iterator=%s\n",
			indentStr, n.Labels, n.Iterator)
		fmt.Printf("%s  ForEach: %v\n", indentStr, getExpressionSummary(n.ForEach))
		for _, child := range n.Content {
			printAST(child, indent+1)
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
	case *types.ForExpr:
		summary := "ForExpr["
		// Show the variable and collection
		if e.KeyVar != "" {
			summary += fmt.Sprintf("for %s, %s in %s",
				e.KeyVar, e.ValueVar, getExpressionSummary(e.Collection))
		} else {
			summary += fmt.Sprintf("for %s in %s",
				e.ValueVar, getExpressionSummary(e.Collection))
		}

		// Show the value expression
		if e.ThenKeyExpr != nil {
			summary += fmt.Sprintf(": %s", getExpressionSummary(e.ThenKeyExpr))
		}

		// Show key expression if present (for map outputs)
		if e.ThenValueExpr != nil {
			summary += fmt.Sprintf(" => %s", getExpressionSummary(e.ThenValueExpr))
		}

		// Show condition if present
		if e.Condition != nil {
			summary += fmt.Sprintf(" if %s", getExpressionSummary(e.Condition))
		}

		summary += "]"
		return summary
	default:
		return fmt.Sprintf("%T", e)
	}
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
