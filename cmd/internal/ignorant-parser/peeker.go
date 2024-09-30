package ignorant_parser

import "github.com/hashicorp/hcl/v2/hclwrite"

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
