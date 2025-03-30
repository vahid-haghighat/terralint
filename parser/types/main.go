package types

import "github.com/hashicorp/hcl/v2"

// Body is the base interface that all Terraform elements must implement
type Body interface {
	BodyType() string
}

// Root represents the top-level HCL document
type Root struct {
	Children []Body
}

func (r *Root) BodyType() string {
	return "root"
}

// FormatDirective represents formatter and linter directives like # tflint-ignore
type FormatDirective struct {
	DirectiveType string   // The directive type (e.g., "tflint-ignore")
	Parameters    []string // Any parameters for the directive
	Range         hcl.Range
}

func (f *FormatDirective) BodyType() string {
	return "format_directive"
}

// Block represents a Terraform block (resource, data, module, etc.)
type Block struct {
	Type          string    // The type of block (resource, data, variable, etc.)
	Labels        []string  // Labels/identifiers for the block (e.g., "aws_instance" "example")
	Range         hcl.Range // Source code position information
	InlineComment string    // Inline comment if present
	BlockComment  string    // Block comment if present
	Children      []Body    // Nested blocks and attributes
}

func (b *Block) BodyType() string {
	return "block"
}

// Attribute represents a key-value pair in HCL
type Attribute struct {
	Name          string     // The name of the attribute
	Value         Expression // The value of the attribute
	Range         hcl.Range
	InlineComment string
	BlockComment  string
}

func (a *Attribute) BodyType() string {
	return "attribute"
}
