package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TestFile represents a Terraform file to be tested
type TestFile struct {
	Path        string
	Description string
}

func main() {
	// Define test files
	testFiles := []TestFile{
		{
			Path:        "complex_terraform_test.tf",
			Description: "Complex Terraform with nested expressions, conditionals, and for loops",
		},
		{
			Path:        "edge_cases_test.tf",
			Description: "Edge cases and unusual syntax patterns",
		},
		{
			Path:        "template_directives_test.tf",
			Description: "Complex template directives and interpolation",
		},
		{
			Path:        "modules_test/main.tf",
			Description: "Complex module structure with multiple dependencies",
		},
		{
			Path:        "modules_test/variables.tf",
			Description: "Complex variable definitions with validations",
		},
		{
			Path:        "modules_test/modules/vpc/main.tf",
			Description: "VPC module implementation",
		},
		{
			Path:        "modules_test/modules/vpc/variables.tf",
			Description: "VPC module variables",
		},
		{
			Path:        "modules_test/modules/vpc/outputs.tf",
			Description: "VPC module outputs",
		},
	}

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure we're in the test directory
	if !strings.HasSuffix(currentDir, "test") {
		fmt.Println("This script should be run from the test directory")
		os.Exit(1)
	}

	fmt.Println("=== Running Terraform Parser Tests ===")
	fmt.Println()

	totalFiles := len(testFiles)
	successCount := 0
	failureCount := 0
	var failures []string

	// Process each test file
	for i, testFile := range testFiles {
		fmt.Printf("[%d/%d] Testing: %s\n", i+1, totalFiles, testFile.Path)
		fmt.Printf("Description: %s\n", testFile.Description)

		// Get the full path to the test file
		fullPath := filepath.Join(currentDir, testFile.Path)

		// Measure parsing time
		startTime := time.Now()

		// Use our debug_parser.go program to parse the file
		cmd := exec.Command("go", "run", "./parser/debug_parser.go", fullPath)
		output, err := cmd.CombinedOutput()
		elapsedTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			failureCount++
			failures = append(failures, testFile.Path)
		} else {
			fmt.Printf("✅ SUCCESS\n")
			fmt.Printf("Output: %s\n", string(output))
			successCount++
		}

		fmt.Printf("Parsing time: %v\n", elapsedTime)
		fmt.Println()
	}

	// Print summary
	fmt.Println("=== Test Summary ===")
	fmt.Printf("Total files tested: %d\n", totalFiles)
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)

	if failureCount > 0 {
		fmt.Println("\nFailed files:")
		for _, failure := range failures {
			fmt.Printf("- %s\n", failure)
		}
		os.Exit(1)
	}
}
