package ignorant_parser

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type peeker struct {
	Tokens    hclwrite.Tokens
	NextIndex int
}

func (p *peeker) Peek() *hclwrite.Token {
	ret, _ := p.nextToken()
	return ret
}

func (p *peeker) nextToken() (*hclwrite.Token, int) {
	if p.NextIndex < len(p.Tokens)-1 {
		return p.Tokens[p.NextIndex], p.NextIndex + 1
	}
	return p.Tokens[len(p.Tokens)-1], len(p.Tokens)
}

func (p *peeker) Read() *hclwrite.Token {
	ret, nextIdx := p.nextToken()
	p.NextIndex = nextIdx
	return ret
}

func (p *peeker) LookAheadFor(tokens hclwrite.Tokens) bool {
	if len(tokens) == 0 || len(p.Tokens) == 0 {
		return false
	}

	index := p.NextIndex
	order := 0

	for index < len(p.Tokens) && p.Tokens[index].Type == hclsyntax.TokenNewline {
		index++
	}

	if index >= len(p.Tokens) {
		return false
	}

	for order < len(tokens) && tokens[order].Type == hclsyntax.TokenNewline {
		order++
	}

	if order >= len(tokens) {
		return false
	}

	// Not enough tokens left to match
	if len(p.Tokens)-index < len(tokens)-order {
		return false
	}

	for i := order; i < len(tokens); i++ {
		if index >= len(p.Tokens) || string(p.Tokens[index].Bytes) != string(tokens[i].Bytes) {
			return false
		}
		index++
	}
	return true
}
