package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vahid-haghighat/terralint/parser/types"
)

// TestCase represents a test case for the parser
type TestCase struct {
	Name        string
	FilePath    string
	Description string
	Expected    types.Body
}

func TestParser(t *testing.T) {
	// Define test cases
	testCases := []TestCase{
		{
			Name:        "Simple Terraform",
			FilePath:    "test_files/simple_test.tf",
			Description: "Simple Terraform file with basic constructs",
			Expected:    createSimpleTerraformExpected(),
		},
		{
			Name:        "Complex Terraform",
			FilePath:    "test_files/complex_terraform_test.tf",
			Description: "Complex Terraform with nested expressions, conditionals, and for loops",
			Expected:    createComplexTerraformExpected(),
		},
		{
			Name:        "Module Structure",
			FilePath:    "test_files/modules_test/main.tf",
			Description: "Module structure with nested modules",
			Expected:    createModuleExpected(),
		},
		{
			Name:        "Edge Cases",
			FilePath:    "test_files/edge_cases_test.tf",
			Description: "Edge cases and unusual syntax patterns",
			Expected:    createEdgeCasesExpected(),
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Get the full path to the test file
			fullPath, err := filepath.Abs(tc.FilePath)
			if err != nil {
				t.Fatalf("Failed to get absolute path for %s: %v", tc.FilePath, err)
			}

			// Check if the file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Fatalf("Test file does not exist: %s", fullPath)
			}

			// Parse the file
			root, err := ParseTerraformFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tc.FilePath, err)
			}

			// Verify the root is not nil
			if root == nil {
				t.Fatalf("Parsed root is nil for %s", tc.FilePath)
			}

			// Compare the parsed structure with the expected structure
			compareStructures(t, tc.Expected, root)
		})
	}
}

// compareStructures recursively compares the expected and actual structures
func compareStructures(t *testing.T, expected, actual types.Body) {
	// Check if types match
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("Type mismatch: expected %T, got %T", expected, actual)
		return
	}

	// Compare based on type
	switch exp := expected.(type) {
	case *types.Root:
		act, ok := actual.(*types.Root)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.Root, got %T", actual)
			return
		}
		compareRoots(t, exp, act)
	case *types.Block:
		act, ok := actual.(*types.Block)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.Block, got %T", actual)
			return
		}
		compareBlocks(t, exp, act)
	case *types.Attribute:
		act, ok := actual.(*types.Attribute)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.Attribute, got %T", actual)
			return
		}
		compareAttributes(t, exp, act)
	case *types.DynamicBlock:
		act, ok := actual.(*types.DynamicBlock)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.DynamicBlock, got %T", actual)
			return
		}
		compareDynamicBlocks(t, exp, act)
	default:
		t.Errorf("Unsupported type for comparison: %T", expected)
	}
}

// compareRoots compares two Root structures
func compareRoots(t *testing.T, expected, actual *types.Root) {
	// Check if the number of children matches
	if len(expected.Children) != len(actual.Children) {
		t.Errorf("Root children count mismatch: expected %d, got %d",
			len(expected.Children), len(actual.Children))
		return
	}

	// Compare each child
	for i, expChild := range expected.Children {
		if i >= len(actual.Children) {
			t.Errorf("Missing child at index %d in actual", i)
			continue
		}
		compareStructures(t, expChild, actual.Children[i])
	}
}

// compareBlocks compares two Block structures
func compareBlocks(t *testing.T, expected, actual *types.Block) {
	// Check if the block type matches
	if expected.Type != actual.Type {
		t.Errorf("Block type mismatch: expected %s, got %s", expected.Type, actual.Type)
	}

	// Check if the labels match
	// Special case: if both are empty (nil or empty slice), consider them equal
	if len(expected.Labels) == 0 && len(actual.Labels) == 0 {
		// Both are empty, so they're considered equal
	} else if !reflect.DeepEqual(expected.Labels, actual.Labels) {
		t.Errorf("Block labels mismatch: expected %v (type: %T), got %v (type: %T)",
			expected.Labels, expected.Labels, actual.Labels, actual.Labels)

		// Print more details about the labels
		t.Logf("Expected labels length: %d", len(expected.Labels))
		t.Logf("Actual labels length: %d", len(actual.Labels))

		if len(expected.Labels) > 0 && len(actual.Labels) > 0 {
			t.Logf("First expected label: %q (type: %T)", expected.Labels[0], expected.Labels[0])
			t.Logf("First actual label: %q (type: %T)", actual.Labels[0], actual.Labels[0])
		}
	}

	// Check if the number of children matches
	if len(expected.Children) != len(actual.Children) {
		t.Errorf("Block %s %v children count mismatch: expected %d, got %d",
			expected.Type, expected.Labels, len(expected.Children), len(actual.Children))

		// Log the actual children for debugging
		t.Logf("Actual children:")
		for i, child := range actual.Children {
			switch c := child.(type) {
			case *types.Block:
				t.Logf("  %d: Block %s %v", i, c.Type, c.Labels)
			case *types.Attribute:
				t.Logf("  %d: Attribute %s", i, c.Name)
			default:
				t.Logf("  %d: %T", i, child)
			}
		}
		return
	}

	// Compare each child
	for i, expChild := range expected.Children {
		if i >= len(actual.Children) {
			t.Errorf("Missing child at index %d in actual", i)
			continue
		}
		compareStructures(t, expChild, actual.Children[i])
	}
}

// compareAttributes compares two Attribute structures
func compareAttributes(t *testing.T, expected, actual *types.Attribute) {
	// Check if the attribute name matches
	if expected.Name != actual.Name {
		t.Errorf("Attribute name mismatch: expected %s, got %s", expected.Name, actual.Name)
	}

	// Compare the attribute values
	compareExpressions(t, expected.Value, actual.Value)
}

// compareDynamicBlocks compares two DynamicBlock structures
func compareDynamicBlocks(t *testing.T, expected, actual *types.DynamicBlock) {
	// Check if the labels match
	if !reflect.DeepEqual(expected.Labels, actual.Labels) {
		t.Errorf("Dynamic block labels mismatch: expected %v, got %v", expected.Labels, actual.Labels)
	}

	// Compare the for_each expressions
	compareExpressions(t, expected.ForEach, actual.ForEach)

	// Check if the iterator matches
	if expected.Iterator != actual.Iterator {
		t.Errorf("Dynamic block iterator mismatch: expected %s, got %s", expected.Iterator, actual.Iterator)
	}

	// Check if the number of content blocks matches
	if len(expected.Content) != len(actual.Content) {
		t.Errorf("Dynamic block content count mismatch: expected %d, got %d",
			len(expected.Content), len(actual.Content))
		return
	}

	// Compare each content block
	for i, expContent := range expected.Content {
		if i >= len(actual.Content) {
			t.Errorf("Missing content at index %d in actual", i)
			continue
		}
		compareStructures(t, expContent, actual.Content[i])
	}
}

// compareExpressions compares two Expression structures
func compareExpressions(t *testing.T, expected, actual types.Expression) {
	// Check if both are nil
	if expected == nil && actual == nil {
		return
	}

	// Check if one is nil but the other isn't
	if (expected == nil && actual != nil) || (expected != nil && actual == nil) {
		t.Errorf("Expression nil mismatch: expected %v, got %v", expected, actual)
		return
	}

	// Check if the expression types match
	if expected.ExpressionType() != actual.ExpressionType() {
		t.Errorf("Expression type mismatch: expected %s, got %s",
			expected.ExpressionType(), actual.ExpressionType())
		return
	}

	// Compare based on expression type
	switch exp := expected.(type) {
	case *types.LiteralValue:
		act, ok := actual.(*types.LiteralValue)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.LiteralValue, got %T", actual)
			return
		}
		// Compare value type
		if exp.ValueType != act.ValueType {
			t.Errorf("Literal value type mismatch: expected %s, got %s", exp.ValueType, act.ValueType)
		}
		// For string literals, compare the string value
		if exp.ValueType == "string" {
			expStr, expOk := exp.Value.(string)
			actStr, actOk := act.Value.(string)
			if expOk && actOk && expStr != actStr {
				t.Errorf("String literal value mismatch: expected %s, got %s", expStr, actStr)
			}
		}
	case *types.ObjectExpr:
		act, ok := actual.(*types.ObjectExpr)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.ObjectExpr, got %T", actual)
			return
		}
		// Check if the number of items matches
		if len(exp.Items) != len(act.Items) {
			t.Errorf("Object items count mismatch: expected %d, got %d", len(exp.Items), len(act.Items))
			return
		}
		// Compare each item
		for i, expItem := range exp.Items {
			if i >= len(act.Items) {
				t.Errorf("Missing object item at index %d in actual", i)
				continue
			}
			compareExpressions(t, expItem.Key, act.Items[i].Key)
			compareExpressions(t, expItem.Value, act.Items[i].Value)
		}
	case *types.ArrayExpr:
		act, ok := actual.(*types.ArrayExpr)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.ArrayExpr, got %T", actual)
			return
		}
		// Check if the number of items matches
		if len(exp.Items) != len(act.Items) {
			t.Errorf("Array items count mismatch: expected %d, got %d", len(exp.Items), len(act.Items))
			return
		}
		// Compare each item
		for i, expItem := range exp.Items {
			if i >= len(act.Items) {
				t.Errorf("Missing array item at index %d in actual", i)
				continue
			}
			compareExpressions(t, expItem, act.Items[i])
		}
	case *types.ReferenceExpr:
		act, ok := actual.(*types.ReferenceExpr)
		if !ok {
			t.Errorf("Type assertion failed: expected *types.ReferenceExpr, got %T", actual)
			return
		}
		// Check if the parts match
		if !reflect.DeepEqual(exp.Parts, act.Parts) {
			t.Errorf("Reference parts mismatch: expected %v, got %v", exp.Parts, act.Parts)
		}
	}
	// Add more expression type comparisons as needed
}

// Helper function to compare string slices
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
