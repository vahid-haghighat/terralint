package ignorant_parser

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"regexp"
	"slices"
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

var listForRegex, _ = regexp.Compile("for\\s.*(?:,\\s.*)?\\sin\\s.*(\\s)?:(\\s)?.*")

type Section struct {
	Type     string
	Labels   []string
	Value    []hclwrite.Tokens
	Comments hclwrite.Tokens
}

func (s *Section) FlattenValue() hclwrite.Tokens {
	var flattened hclwrite.Tokens
	for _, tokens := range s.Value {
		flattened = append(flattened, tokens...)
	}
	return slices.Clip(flattened)
}

func (s *Section) LineCounts() int {
	count := 0
	for _, token := range s.FlattenValue() {
		if token.Type == hclsyntax.TokenNewline {
			count++
		}
	}
	return count
}

func (s *Section) HasValue() bool {
	for _, tokens := range s.Value {
		if len(tokens) > 0 {
			return true
		}
	}
	return s.Type != "" || len(s.Labels) > 0 || len(s.Comments) > 0
}

func (s *Section) IsAttribute() bool {
	return isAttribute(s.FlattenValue())
}

func (s *Section) IsBlock() bool {
	return isBlock(s.FlattenValue())
}

func (s *Section) IsList() bool {
	return isList(s.Value)
}

func (s *Section) Tokens() hclwrite.Tokens {
	if !s.HasValue() {
		return nil
	}

	file := hclwrite.NewEmptyFile()
	file.Body().AppendUnstructuredTokens(s.Comments)
	if s.Value == nil {
		return file.BuildTokens(nil)
	}
	var tokens hclwrite.Tokens
	sectionTokens := GetSectionBody(s.FlattenValue())
	if s.IsAttribute() {
		if s.IsBlock() {
			if sectionTokens == nil || len(sectionTokens) == 0 {
				tokens = hclwrite.Tokens{
					&tokenOBrace,
					&tokenCBrace,
					&tokenNewLine,
				}
			} else {
				tokens = append(tokens, hclwrite.Tokens{
					&tokenOBrace,
					&tokenNewLine,
				}...)
				tokens = append(tokens, sectionTokens...)
				if !tokenIsNewline(tokens[len(tokens)-1]) {
					tokens = append(tokens, &tokenNewLine)
				}
				tokens = append(tokens, hclwrite.Tokens{
					&tokenCBrace,
					&tokenNewLine,
				}...)
			}
		} else {
			tokens = sectionTokens
		}
		file.Body().SetAttributeRaw(s.Type, tokens)
		tokens = file.BuildTokens(nil)
		tokens = tokens[:len(tokens)-1]
	} else if s.IsBlock() {
		if sectionTokens == nil || len(sectionTokens) == 0 {
			emptyBlockTokens := hclwrite.Tokens{
				&hclwrite.Token{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(s.Type),
					SpacesBefore: 0,
				},
			}
			for _, label := range s.Labels {
				emptyBlockTokens = append(emptyBlockTokens, hclwrite.Tokens{
					&tokenOQuote,
					&hclwrite.Token{
						Type:         hclsyntax.TokenIdent,
						Bytes:        []byte(label),
						SpacesBefore: 0,
					},
					&tokenCQuote,
				}...)
			}
			emptyBlockTokens = append(emptyBlockTokens, hclwrite.Tokens{
				&tokenOBrace,
				&tokenCBrace,
				&tokenNewLine,
			}...)
			file.Body().AppendUnstructuredTokens(emptyBlockTokens)
		} else {
			block := hclwrite.NewBlock(s.Type, s.Labels)
			block.Body().AppendUnstructuredTokens(GetSectionBody(s.FlattenValue()))
			file.Body().AppendBlock(block)
		}
		file.Body().AppendNewline()
		tokens = file.BuildTokens(nil)
	} else {
		tokens = file.BuildTokens(s.FlattenValue())
	}
	tokens = append(tokens, &tokenNewLine)
	tokens = trimEndNewLines(tokens)
	return tokens
}

func (s *Section) ListCount() int {
	return len(s.Value) - 2
}

func trimEndNewLines(tokens hclwrite.Tokens) hclwrite.Tokens {
	if tokens == nil {
		return nil
	}

	if len(tokens) == 1 {
		if tokenIsNewline(tokens[0]) {
			return nil
		}
		return tokens
	}

	index := len(tokens) - 1
	for index >= 0 && tokenIsNewline(tokens[index]) {
		index--
	}

	if index+2 > len(tokens) {
		return tokens
	}
	return tokens[:index+2]
}

func isSomething(tokens hclwrite.Tokens, tokenType hclsyntax.TokenType, ignoreTokens []hclsyntax.TokenType) bool {
	if tokens == nil {
		return false
	}

	index := 0
	for index < len(tokens) && !tokenIsNewline(tokens[index]) {
		index++
	}

	for i := 0; i < index; i++ {
		if tokens[i].Type == tokenType {
			return true
		} else if !contains(tokens[i].Type, ignoreTokens) {
			return false
		}
	}
	return false
}

func isAttribute(tokens hclwrite.Tokens) bool {
	return isSomething(tokens, hclsyntax.TokenEqual, []hclsyntax.TokenType{})
}

func isBlock(tokens hclwrite.Tokens) bool {
	return isSomething(tokens, hclsyntax.TokenOBrace, []hclsyntax.TokenType{hclsyntax.TokenEqual})
}

func isList(tokens []hclwrite.Tokens) bool {
	return isSomething(tokens[0], hclsyntax.TokenOBrack, []hclsyntax.TokenType{hclsyntax.TokenEqual, hclsyntax.TokenNewline})
}

func contains(token hclsyntax.TokenType, tokens []hclsyntax.TokenType) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

func convertToHclwrite(tokens hclsyntax.Tokens) hclwrite.Tokens {
	if len(tokens) == 0 {
		return nil
	}
	result := make(hclwrite.Tokens, len(tokens))
	for index, token := range tokens {
		result[index] = &hclwrite.Token{
			Type:         token.Type,
			Bytes:        token.Bytes,
			SpacesBefore: 0,
		}
	}

	return result
}
