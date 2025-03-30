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
			Name:        "Complex Module",
			FilePath:    "test_files/complex_terraform_split/01_complex_module.tf",
			Description: "Complex module with nested expressions, conditionals, and for loops",
			Expected:    createComplexModuleExpected(),
		},
		{
			Name:        "Complex Resource",
			FilePath:    "test_files/complex_terraform_split/02_complex_resource.tf",
			Description: "Resource with complex dynamic blocks and for_each",
			Expected:    createComplexResourceExpected(),
		},
		{
			Name:        "Complex Locals",
			FilePath:    "test_files/complex_terraform_split/03_complex_locals.tf",
			Description: "Complex locals with nested expressions",
			Expected:    createComplexLocalsExpected(),
		},
		{
			Name:        "Complex Data Source",
			FilePath:    "test_files/complex_terraform_split/04_complex_data_source.tf",
			Description: "Data source with complex expressions",
			Expected:    createComplexDataSourceExpected(),
		},
		{
			Name:        "Complex Variable",
			FilePath:    "test_files/complex_terraform_split/05_complex_variable.tf",
			Description: "Variable with complex type constraints and validations",
			Expected:    createComplexVariableExpected(),
		},
		{
			Name:        "Complex Output",
			FilePath:    "test_files/complex_terraform_split/06_complex_output.tf",
			Description: "Output with complex expressions",
			Expected:    createComplexOutputExpected(),
		},
		{
			Name:        "Complex Provider",
			FilePath:    "test_files/complex_terraform_split/07_complex_provider.tf",
			Description: "Provider configuration with complex expressions",
			Expected:    createComplexProviderExpected(),
		},
		{
			Name:        "Complex Terraform Config",
			FilePath:    "test_files/complex_terraform_split/08_complex_terraform_config.tf",
			Description: "Terraform configuration with complex expressions",
			Expected:    createComplexTerraformConfigExpected(),
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

	// Check if the block comment matches
	if expected.BlockComment != actual.BlockComment {
		t.Errorf("Block comment mismatch for block %s: expected %q, got %q",
			expected.Type, expected.BlockComment, actual.BlockComment)
	}

	// Check if the inline comment matches
	if expected.InlineComment != actual.InlineComment {
		t.Errorf("Inline comment mismatch for block %s: expected %q, got %q",
			expected.Type, expected.InlineComment, actual.InlineComment)
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

	// Check if the block comment matches
	if expected.BlockComment != actual.BlockComment {
		t.Errorf("Block comment mismatch for attribute %s: expected %q, got %q",
			expected.Name, expected.BlockComment, actual.BlockComment)
	}

	// Check if the inline comment matches
	if expected.InlineComment != actual.InlineComment {
		t.Errorf("Inline comment mismatch for attribute %s: expected %q, got %q",
			expected.Name, expected.InlineComment, actual.InlineComment)
	}

	// Compare the attribute values
	compareExpressions(t, expected.Value, actual.Value)
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
		// Print more details about the expressions
		t.Errorf("Expression type mismatch: expected %s, got %s. Expected: %T, Actual: %T",
			expected.ExpressionType(), actual.ExpressionType(), expected, actual)

		// Print more details for reference expressions
		if ref, ok := expected.(*types.ReferenceExpr); ok {
			t.Logf("Expected reference parts: %v", ref.Parts)
		}
		if ref, ok := actual.(*types.ReferenceExpr); ok {
			t.Logf("Actual reference parts: %v", ref.Parts)
		}

		// Print more details for array expressions
		if arr, ok := expected.(*types.ArrayExpr); ok {
			t.Logf("Expected array items count: %d", len(arr.Items))
		}
		if arr, ok := actual.(*types.ArrayExpr); ok {
			t.Logf("Actual array items count: %d", len(arr.Items))
		}

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
			// Compare comments for object items
			if expItem.BlockComment != act.Items[i].BlockComment {
				t.Errorf("Object item block comment mismatch at index %d: expected %q, got %q",
					i, expItem.BlockComment, act.Items[i].BlockComment)
			}
			if expItem.InlineComment != act.Items[i].InlineComment {
				t.Errorf("Object item inline comment mismatch at index %d: expected %q, got %q",
					i, expItem.InlineComment, act.Items[i].InlineComment)
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

			// Print the array items for debugging
			t.Logf("Expected array items:")
			for i, item := range exp.Items {
				switch v := item.(type) {
				case *types.LiteralValue:
					t.Logf("  %d: %v (%s)", i, v.Value, v.ValueType)
				default:
					t.Logf("  %d: %T", i, item)
				}
			}

			t.Logf("Actual array items:")
			for i, item := range act.Items {
				switch v := item.(type) {
				case *types.LiteralValue:
					t.Logf("  %d: %v (%s)", i, v.Value, v.ValueType)
				default:
					t.Logf("  %d: %T", i, item)
				}
			}

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
