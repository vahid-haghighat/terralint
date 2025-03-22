# Terraform Linter Parser Tests

This directory contains complex Terraform test files designed to test the robustness and capabilities of the Terraform linter parser. These tests include edge cases, complex expressions, and unusual syntax patterns to ensure the parser can handle a wide variety of Terraform code.

## Test Files

1. **complex_terraform_test.tf**
   - Contains extremely complex Terraform code with nested expressions, conditionals, and for loops
   - Tests the parser's ability to handle deeply nested structures and complex interpolations

2. **edge_cases_test.tf**
   - Contains edge cases and unusual syntax patterns
   - Tests the parser's robustness against unexpected but valid Terraform code

3. **template_directives_test.tf**
   - Focuses on complex template directives and interpolation
   - Tests the parser's ability to handle complex string templates with directives

4. **modules_test/**
   - Contains a complex module structure with multiple dependencies
   - Tests the parser's ability to handle multi-file Terraform projects
   - Includes main.tf, variables.tf, and a VPC module

## Running the Tests

To run the parser tests, use the following command from the test directory:

```bash
go run run_parser_tests.go
```

This will attempt to parse each test file using the Terraform linter parser and report the results.

## Test Structure

Each test file is designed to exercise specific aspects of the Terraform language:

- **Complex expressions**: Nested function calls, conditionals, for expressions
- **Dynamic blocks**: Complex dynamic block structures with nested blocks
- **Template directives**: String interpolation with complex directives
- **Module structure**: Multi-file projects with dependencies
- **Edge cases**: Unusual but valid syntax patterns

## Adding New Tests

When adding new test files:

1. Create a new .tf file with the test cases
2. Add the file to the `testFiles` list in `run_parser_tests.go`
3. Run the tests to ensure the parser can handle the new cases

## Expected Results

The parser should be able to successfully parse all test files without errors. Any parsing failures indicate potential issues with the parser implementation that need to be addressed.

## Using These Tests for Development

These tests are valuable during development to ensure that changes to the parser don't break existing functionality. Run these tests after making changes to the parser to verify that it can still handle complex Terraform code correctly.