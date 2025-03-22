package types

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type RawSection struct {
	StandaloneComment []*sitter.Node
	InlineComment     *sitter.Node
	Body              *sitter.Node
	Children          []*RawSection
	Text              string
}

func (s *RawSection) IsEmpty() bool {
	return len(s.StandaloneComment) == 0 && s.InlineComment == nil && s.Body == nil && len(s.Children) == 0
}
