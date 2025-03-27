# Complex Terraform Test Files

This directory contains a breakdown of the complex Terraform test file (`complex_terraform_test.tf`) into smaller, more focused files to facilitate easier and more manageable unit testing of the parser.

## Purpose

The original file contains extremely complex Terraform code with nested expressions, complex interpolations, and edge cases. By breaking it down into smaller pieces, we can:

1. Test specific Terraform constructs in isolation
2. Create more targeted unit tests for the parser
3. Make it easier to identify and debug parsing issues
4. Improve test coverage by focusing on specific syntax patterns

## Files Structure

The complex Terraform code has been split into the following files:

1. `01_complex_module.tf` - Tests parsing of complex module definitions with various attribute types
2. `02_complex_resource.tf` - Tests parsing of resources with complex dynamic blocks and for_each
3. `03_complex_locals.tf` - Tests parsing of locals with nested expressions and transformations
4. `04_complex_data_source.tf` - Tests parsing of data sources with dynamic blocks and complex expressions
5. `05_complex_variable.tf` - Tests parsing of variables with complex type constraints and validations
6. `06_complex_output.tf` - Tests parsing of outputs with nested expressions and dependencies
7. `07_complex_provider.tf` - Tests parsing of provider configurations with complex expressions
8. `08_complex_terraform_config.tf` - Tests parsing of Terraform configuration blocks

## Usage for Unit Testing

These files can be used to create more focused unit tests for the parser. For example:

```go
func TestParseComplexModule(t *testing.T) {
    content, err := ioutil.ReadFile("test_files/complex_terraform_split/01_complex_module.tf")
    if err != nil {
        t.Fatalf("Failed to read test file: %v", err)
    }
    
    // Parse the content
    result, err := parser.Parse(string(content))
    if err != nil {
        t.Fatalf("Failed to parse: %v", err)
    }
    
    // Assert specific aspects of the module parsing
    if len(result.Modules) != 1 {
        t.Errorf("Expected 1 module, got %d", len(result.Modules))
    }
    
    module := result.Modules[0]
    if module.Name != "complex_module" {
        t.Errorf("Expected module name 'complex_module', got '%s'", module.Name)
    }
    
    // Test specific attributes and expressions within the module
    // ...
}
```

Similar tests can be created for each of the other files, focusing on the specific Terraform constructs they contain.

## Notes

- These files are not meant to be valid Terraform configurations that can be applied. They are specifically designed to test the parser's ability to handle complex Terraform syntax.
- The files contain references to variables and resources that are not declared, which is expected and intentional for testing purposes.
- When writing unit tests, focus on the parser's ability to correctly identify and parse the syntax structures, not on the validity of the Terraform configuration itself.