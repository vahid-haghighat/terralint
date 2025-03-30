package types

import sitter "github.com/smacker/go-tree-sitter"

// Expression interface represents any valid Terraform expression
type Expression interface {
	ExpressionType() string
	Range() sitter.Range
}

// LiteralValue represents primitive values (string, number, bool)
type LiteralValue struct {
	Value     interface{} // The actual value
	ValueType string      // Type of the value (string, number, bool)
	ExprRange sitter.Range
}

func (l *LiteralValue) ExpressionType() string {
	return "literal"
}

func (l *LiteralValue) Range() sitter.Range {
	return l.ExprRange
}

// ObjectExpr represents object/map expressions
type ObjectExpr struct {
	Items     []ObjectItem
	ExprRange sitter.Range
}

type ObjectItem struct {
	Key           Expression
	Value         Expression
	InlineComment string
	BlockComment  string // Added to capture block comments
	ExprRange     sitter.Range
}

func (o *ObjectItem) ExpressionType() string {
	return "object_item"
}

func (o *ObjectItem) Range() sitter.Range {
	return o.ExprRange
}

func (o *ObjectExpr) ExpressionType() string {
	return "object"
}

func (o *ObjectExpr) Range() sitter.Range {
	return o.ExprRange
}

// ArrayExpr represents array/list expressions
type ArrayExpr struct {
	Items     []Expression
	ExprRange sitter.Range
}

func (a *ArrayExpr) ExpressionType() string {
	return "array"
}

func (a *ArrayExpr) Range() sitter.Range {
	return a.ExprRange
}

// ReferenceExpr represents variable references and attribute access
type ReferenceExpr struct {
	Parts     []string // e.g., ["var", "environment"] for var.environment
	ExprRange sitter.Range
}

func (r *ReferenceExpr) ExpressionType() string {
	return "reference"
}

func (r *ReferenceExpr) Range() sitter.Range {
	return r.ExprRange
}

// FunctionCallExpr represents function calls
type FunctionCallExpr struct {
	Name      string
	Args      []Expression
	ExprRange sitter.Range
}

func (f *FunctionCallExpr) ExpressionType() string {
	return "function_call"
}

func (f *FunctionCallExpr) Range() sitter.Range {
	return f.ExprRange
}

// TemplateExpr represents string interpolation
type TemplateExpr struct {
	Parts     []Expression // Mix of LiteralValue, TemplateForDirective, TemplateIfDirective, and other expressions
	ExprRange sitter.Range
}

func (t *TemplateExpr) ExpressionType() string {
	return "template"
}

func (t *TemplateExpr) Range() sitter.Range {
	return t.ExprRange
}

// ConditionalExpr represents conditional expressions (condition ? true_val : false_val)
type ConditionalExpr struct {
	Condition Expression
	TrueExpr  Expression
	FalseExpr Expression
	ExprRange sitter.Range
}

func (c *ConditionalExpr) ExpressionType() string {
	return "conditional"
}

func (c *ConditionalExpr) Range() sitter.Range {
	return c.ExprRange
}

// BinaryExpr represents binary operations
type BinaryExpr struct {
	Left      Expression
	Operator  string // Should support: ==, !=, <, >, <=, >=, &&, ||, +, -, *, /, %, etc.
	Right     Expression
	ExprRange sitter.Range
}

func (b *BinaryExpr) ExpressionType() string {
	return "binary"
}

func (b *BinaryExpr) Range() sitter.Range {
	return b.ExprRange
}

// ForMapExpr represents for expressions: [for x in xs : upper(x)] or {for k, v in map : k => v}
type ForMapExpr struct {
	// Iterator variables
	ValueVar string // The value variable name (e.g. "v" in "for k, v in map")
	KeyVar   string // The key variable for maps (e.g., "k" in "for k, v in map" or "x" in "for x in xs")

	// Collection being iterated
	Collection Expression // The collection being iterated over (e.g., "xs" in "for x in xs")

	// Result expressions
	ThenKeyExpr   Expression // The value expression (e.g., "k" in "for k, v in map : k => v")
	ThenValueExpr Expression // The value expression for map outputs (e.g., "v" in "for k, v in map : k => v")

	// Filtering and grouping
	Condition Expression // Optional "if" condition (e.g., "x != null" in "for x in xs : x if x != null")

	// Source location
	ExprRange sitter.Range
}

func (f *ForMapExpr) ExpressionType() string {
	return "for_map"
}

func (f *ForMapExpr) Range() sitter.Range {
	return f.ExprRange
}

// ForArrayExpr represents for expressions: [for x in xs : upper(x)]
type ForArrayExpr struct {
	// Iterator variables
	ValueVar string // The value variable name (e.g. "v" in "for k, v in map")
	KeyVar   string // The key variable for maps (e.g., "k" in "for k, v in map" or "x" in "for x in xs")

	// Collection being iterated
	Collection Expression // The collection being iterated over (e.g., "xs" in "for x in xs")

	// Result expressions
	ThenValueExpr Expression // The value expression for map outputs (e.g., "v" in "for k, v in map : k => v")

	// Filtering and grouping
	Condition Expression // Optional "if" condition (e.g., "x != null" in "for x in xs : x if x != null")

	// Source location
	ExprRange sitter.Range
}

func (f *ForArrayExpr) ExpressionType() string {
	return "for_array"
}

func (f *ForArrayExpr) Range() sitter.Range {
	return f.ExprRange
}

// SplatExpr represents splat expressions: aws_instance.server[*].id
type SplatExpr struct {
	Source    Expression // Expression being splattered
	Each      Expression // Expression to evaluate for each element
	ExprRange sitter.Range
}

func (s *SplatExpr) ExpressionType() string {
	return "splat"
}

func (s *SplatExpr) Range() sitter.Range {
	return s.ExprRange
}

// HeredocExpr represents heredoc strings
type HeredocExpr struct {
	Marker    string // The heredoc marker (e.g., "EOT")
	Content   string // The content of the heredoc
	Indented  bool   // Whether it's an indented heredoc (<<-)
	ExprRange sitter.Range
}

func (h *HeredocExpr) ExpressionType() string {
	return "heredoc"
}

func (h *HeredocExpr) Range() sitter.Range {
	return h.ExprRange
}

// IndexExpr represents index access operations: list[0], map["key"]
type IndexExpr struct {
	Collection Expression
	Key        Expression
	ExprRange  sitter.Range
}

func (i *IndexExpr) ExpressionType() string {
	return "index"
}

func (i *IndexExpr) Range() sitter.Range {
	return i.ExprRange
}

// TupleExpr represents tuple expressions
type TupleExpr struct {
	Expressions []Expression
	ExprRange   sitter.Range
}

func (t *TupleExpr) ExpressionType() string {
	return "tuple"
}

func (t *TupleExpr) Range() sitter.Range {
	return t.ExprRange
}

// Add UnaryExpr for operations like ! and -
type UnaryExpr struct {
	Operator  string // !, -
	Expr      Expression
	ExprRange sitter.Range
}

func (u *UnaryExpr) ExpressionType() string {
	return "unary"
}

func (u *UnaryExpr) Range() sitter.Range {
	return u.ExprRange
}

// ParenExpr for grouped expressions (a + b) * c
type ParenExpr struct {
	Expression Expression
	ExprRange  sitter.Range
}

func (p *ParenExpr) ExpressionType() string {
	return "paren"
}

func (p *ParenExpr) Range() sitter.Range {
	return p.ExprRange
}

// RelativeTraversalExpr for attribute access like aws_instance.example.id
type RelativeTraversalExpr struct {
	Source    Expression
	Traversal []TraversalElem
	ExprRange sitter.Range
}

type TraversalElem struct {
	Type  string     // "attr" or "index"
	Name  string     // For attribute access
	Index Expression // For index access
}

func (r *RelativeTraversalExpr) ExpressionType() string {
	return "relative_traversal"
}

func (r *RelativeTraversalExpr) Range() sitter.Range {
	return r.ExprRange
}

// TemplateForDirective represents for loops within template strings
type TemplateForDirective struct {
	KeyVar    string       // Optional key variable for maps
	ValueVar  string       // Value variable
	CollExpr  Expression   // Collection to iterate over
	Content   []Expression // Content to repeat for each iteration
	ExprRange sitter.Range
}

func (t *TemplateForDirective) ExpressionType() string {
	return "template_for"
}

func (t *TemplateForDirective) Range() sitter.Range {
	return t.ExprRange
}

// TemplateIfDirective represents conditional logic within template strings
type TemplateIfDirective struct {
	Condition Expression
	TrueExpr  []Expression
	FalseExpr []Expression
	ExprRange sitter.Range
}

func (t *TemplateIfDirective) ExpressionType() string {
	return "template_if"
}

func (t *TemplateIfDirective) Range() sitter.Range {
	return t.ExprRange
}

// TypeExpr represents type constraints like list(string)
type TypeExpr struct {
	TypeName   string       // Base type name (e.g., "list", "map")
	Parameters []Expression // Type parameters if any
	ExprRange  sitter.Range
}

func (t *TypeExpr) ExpressionType() string {
	return "type"
}

func (t *TypeExpr) Range() sitter.Range {
	return t.ExprRange
}
