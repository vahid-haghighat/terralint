# Terraform Parser Tests

This directory contains test files for the Terraform parser. These files are designed to test the parser's ability to handle various Terraform syntax constructs.

## Test Files

### simple_test.tf

A simple Terraform file with basic constructs:
- Resource block
- Variable block
- Output block
- Data source block
- Provider block
- Locals block
- Module block
- Filter blocks (nested in data source)

This file is used to test the basic functionality of the parser.

### complex_terraform_test.tf

A complex Terraform file with advanced constructs:
- Nested expressions
- Conditional expressions
- For expressions
- Binary expressions
- Template expressions
- Heredoc strings
- Complex function calls
- Dynamic blocks
- Complex type constraints
- Validation blocks

This file is used to test the parser's ability to handle complex Terraform syntax. It's designed to be a comprehensive test of all the features that the parser should support.

### edge_cases_test.tf

A Terraform file with edge cases and unusual syntax patterns:
- Empty blocks
- Blocks with only comments
- Extremely long lines
- Unusual whitespace
- Nested conditionals
- Multiple blocks with the same name
- Unicode characters
- Escaped sequences
- Empty strings and special values
- Unusual numbers
- Special characters in identifiers
- Unusual block nesting
- Unusual function calls
- Comments in unusual places
- Unusual heredoc
- Unusual for expressions
- Unusual splat expressions

This file is used to test the parser's ability to handle edge cases and unusual syntax patterns.

### template_directives_test.tf

A Terraform file with complex template directives and interpolation:
- Complex templates with multiple interpolations
- Templates with nested expressions
- Templates with function calls
- Templates with math expressions
- Templates with references to other resources
- Templates with for expressions
- Templates with conditional expressions
- Templates with strip markers
- Templates with indentation
- Templates with escaping
- Templates with special characters
- Templates with unicode
- Templates with newlines and tabs

This file is used to test the parser's ability to handle complex template directives and interpolation.

### modules_test/

A directory containing a Terraform module structure:
- main.tf: The main Terraform configuration file
- variables.tf: Variable definitions
- modules/vpc/: A VPC module
  - main.tf: The main module configuration file
  - variables.tf: Module variable definitions
  - outputs.tf: Module outputs

This directory is used to test the parser's ability to handle module structures and dependencies.

## Running the Tests

The tests are run using Go's testing framework. To run the tests, use the following command:

```bash
cd parser
go test -v
```

## Adding New Tests

When adding new tests, consider the following:
1. Create a new test file with a descriptive name
2. Add a test case to the `TestParser` function in `parser/parser_test.go`
3. Define the expected blocks, attributes, and expression types
4. Run the tests to verify that the parser handles the new test file correctly