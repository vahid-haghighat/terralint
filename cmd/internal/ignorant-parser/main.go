package ignorant_parser

import (
	"errors"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
)

func ParseConfigFromFile(filePath string) ([]*Section, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	parsedFile, diags := hclwrite.ParseConfig(file, filePath, hcl.Pos{})

	if diags != nil && diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	if parsedFile == nil {
		return nil, errors.New("failed to parse file")
	}

	p := &parser{
		peeker: &peeker{
			Tokens:    parsedFile.BuildTokens(nil),
			NextIndex: 0,
		},
	}

	sections, err := p.ParseBody()

	if err != nil {
		return nil, err
	}

	return sections, nil
}

// ParseSectionConfig This method assumes the correctness of the syntax. Use with care
func ParseSectionConfig(tokens hclwrite.Tokens) ([]*Section, error) {
	tokens = GetSectionBody(tokens)
	if tokens == nil || len(tokens) == 0 {
		return nil, nil
	}

	if tokens[len(tokens)-1].Type != hclsyntax.TokenEOF {
		tokens = append(tokens, &hclwrite.Token{
			Type:         hclsyntax.TokenEOF,
			Bytes:        []byte{},
			SpacesBefore: 0,
		})
	}

	p := &parser{
		peeker: &peeker{
			Tokens:    tokens,
			NextIndex: 0,
		},
	}
	sections, err := p.ParseBody()

	if err != nil {
		return nil, err
	}
	return sections, nil
}
