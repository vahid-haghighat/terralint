package types

type Section struct {
	StandaloneComment []string
	InlineComment     string
	Body              *Body
	Children          []*Section
	Text              string
}
