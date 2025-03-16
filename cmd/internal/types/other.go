package types

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
}

// Range represents a range in the source code
type Range struct {
	Start Position
	End   Position
}

//// FormatDirective represents formatter directives like # format: off
//type FormatDirective struct {
//	Directive string
//	Range     Range
//}
