package ignorant_parser

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"regexp"
)

var tokenOBrace = hclwrite.Token{
	Type:         hclsyntax.TokenOBrace,
	Bytes:        []byte("{"),
	SpacesBefore: 0,
}

var tokenCBrace = hclwrite.Token{
	Type:         hclsyntax.TokenCBrace,
	Bytes:        []byte("}"),
	SpacesBefore: 0,
}

var tokenEqual = hclwrite.Token{
	Type:         hclsyntax.TokenEqual,
	Bytes:        []byte("="),
	SpacesBefore: 0,
}

var tokenNewLine = hclwrite.Token{
	Type:         hclsyntax.TokenNewline,
	Bytes:        []byte("\n"),
	SpacesBefore: 0,
}

var tokenOBrack = hclwrite.Token{
	Type:         hclsyntax.TokenOBrack,
	Bytes:        []byte("["),
	SpacesBefore: 0,
}

var tokenCBrack = hclwrite.Token{
	Type:         hclsyntax.TokenCBrack,
	Bytes:        []byte("]"),
	SpacesBefore: 0,
}

var tokenOQuote = hclwrite.Token{
	Type:         hclsyntax.TokenOQuote,
	Bytes:        []byte("\""),
	SpacesBefore: 0,
}

var tokenCQuote = hclwrite.Token{
	Type:         hclsyntax.TokenCQuote,
	Bytes:        []byte("\""),
	SpacesBefore: 0,
}

var tokenComma = hclwrite.Token{
	Type:         hclsyntax.TokenComma,
	Bytes:        []byte(","),
	SpacesBefore: 0,
}

var tokenFor = hclwrite.Token{
	Type:         hclsyntax.TokenIdent,
	Bytes:        []byte("for"),
	SpacesBefore: 0,
}

var listForRegex, _ = regexp.Compile("for\\s.*(?:,\\s.*)?\\sin\\s.*(\\s)?:(\\s)?.*")

type expressionType int

const (
	unsetExpression expressionType = iota
	attributeExpression
	blockExpression
	blockAttributeExpression
	listExpression
	listForExpression
	blockForExpression
)

type Section struct {
	Type           string
	Labels         []string
	Value          []hclwrite.Tokens
	Comments       hclwrite.Tokens
	expressionType expressionType
}

type ForExpression struct {
	Key   hclwrite.Tokens
	Value hclwrite.Tokens
	Body  *ForExpressionBody
}

type ForExpressionBody struct {
	Key  hclwrite.Tokens
	Body hclwrite.Tokens
}
