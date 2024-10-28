package ignorant_parser

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"strings"
)

type parser struct {
	*peeker
}

var openingTokens = []hclsyntax.TokenType{
	hclsyntax.TokenOBrace,
	hclsyntax.TokenOBrack,
	hclsyntax.TokenOParen,
	hclsyntax.TokenOHeredoc,
	hclsyntax.TokenTemplateControl,
	hclsyntax.TokenTemplateInterp,
}

var newLineTokens = []hclsyntax.TokenType{
	hclsyntax.TokenNewline,
	hclsyntax.TokenComment,
}

func oppositeBracket(ty hclsyntax.TokenType) hclsyntax.TokenType {
	switch ty {

	case hclsyntax.TokenOBrace:
		return hclsyntax.TokenCBrace
	case hclsyntax.TokenOBrack:
		return hclsyntax.TokenCBrack
	case hclsyntax.TokenOParen:
		return hclsyntax.TokenCParen
	case hclsyntax.TokenOQuote:
		return hclsyntax.TokenCQuote
	case hclsyntax.TokenOHeredoc:
		return hclsyntax.TokenCHeredoc

	case hclsyntax.TokenCBrace:
		return hclsyntax.TokenOBrace
	case hclsyntax.TokenCBrack:
		return hclsyntax.TokenOBrack
	case hclsyntax.TokenCParen:
		return hclsyntax.TokenOParen
	case hclsyntax.TokenCQuote:
		return hclsyntax.TokenOQuote
	case hclsyntax.TokenCHeredoc:
		return hclsyntax.TokenOHeredoc

	case hclsyntax.TokenTemplateControl:
		return hclsyntax.TokenTemplateSeqEnd
	case hclsyntax.TokenTemplateInterp:
		return hclsyntax.TokenTemplateSeqEnd
	case hclsyntax.TokenTemplateSeqEnd:
		return hclsyntax.TokenTemplateInterp

	default:
		return hclsyntax.TokenNil
	}
}

func isOpeningToken(tokenType hclsyntax.TokenType) bool {
	for _, item := range openingTokens {
		if tokenType == item {
			return true
		}
	}
	return false
}

func isClosingBracket(ty hclsyntax.TokenType) bool {
	return isOpeningToken(oppositeBracket(ty))
}

func isEnd(token *hclwrite.Token, tokens []hclsyntax.TokenType) bool {
	if token.Type == hclsyntax.TokenEOF {
		return true
	}

	for _, item := range tokens {
		if token.Type == item {
			return true
		}
	}
	return false
}

func (p *parser) ReadTokensUntil(end []hclsyntax.TokenType) hclwrite.Tokens {
	var buffer hclwrite.Tokens

Token:
	for {
		next := p.Read()
		if next.Type != hclsyntax.TokenEOF {
			buffer = append(buffer, next)
		}
		if isEnd(next, end) {
			if next.Type == hclsyntax.TokenEOF {
				p.Read()
			}
			break Token
		}
		if isOpeningToken(next.Type) {
			buffer = append(buffer, p.ReadTokensUntil([]hclsyntax.TokenType{oppositeBracket(next.Type)})...)
			continue
		}
	}

	return buffer
}

func (p *parser) ParseBody() ([]*Section, error) {
	var rawSections []*Section
	var commentBuffer hclwrite.Tokens
	var bodyBuffer hclwrite.Tokens
Token:
	for {
		next := p.Peek()

		if next.Type == hclsyntax.TokenEOF {
			p.Read()
			if bodyBuffer != nil {
				bodyBuffer = append(bodyBuffer, &tokenNewLine)
			}
			break Token
		}

		switch next.Type {
		case hclsyntax.TokenNewline:
			p.Read()
		case hclsyntax.TokenComment:
			commentBuffer = append(commentBuffer, p.Read())
		case hclsyntax.TokenOBrace, hclsyntax.TokenEqual:
			headline, err := getHeadlines(bodyBuffer)
			if err != nil {
				return nil, err
			}
			name := headline[0]
			var labels []string
			if len(headline) > 1 {
				labels = headline[1:]
				for index := range labels {
					labels[index] = strings.Trim(labels[index], "\"")
				}
			}
			section := &Section{
				Type:     name,
				Labels:   labels,
				Comments: commentBuffer,
				Value:    []hclwrite.Tokens{{}},
			}

			if next.Type == hclsyntax.TokenEqual {
				section.Value[0] = append(section.Value[0], p.Read())
			}

			if p.Peek().Type == hclsyntax.TokenOBrack {
				p.Read()
				section.Value[0] = append(section.Value[0], &tokenOBrack)
				var lastToken hclsyntax.TokenType
				for lastToken != hclsyntax.TokenCBrack {
					s := p.ReadTokensUntil([]hclsyntax.TokenType{hclsyntax.TokenCBrack, hclsyntax.TokenComma})
					lastToken = s[len(s)-1].Type
					start := 0
					for start < len(s) && s[start].Type == hclsyntax.TokenNewline {
						start++
					}

					end := len(s) - 2
					for end >= 0 && s[end].Type == hclsyntax.TokenNewline {
						end--
					}

					if start <= end {
						section.Value = append(section.Value, s[start:end+1])
					}
				}
				section.Value = append(section.Value, hclwrite.Tokens{&tokenCBrack})
			} else {
				section.Value[0] = append(section.Value[0], p.ReadTokensUntil(newLineTokens)...)
			}

			rawSections = append(rawSections, section)
			bodyBuffer = nil
			commentBuffer = nil
		case hclsyntax.TokenIdent, hclsyntax.TokenDot:
			bodyBuffer = append(bodyBuffer, p.Read())
		case hclsyntax.TokenOQuote:
			bodyBuffer = append(bodyBuffer, p.ReadTokensUntil([]hclsyntax.TokenType{hclsyntax.TokenCQuote})...)
		default:
			bodyBuffer = append(bodyBuffer, p.ReadTokensUntil([]hclsyntax.TokenType{hclsyntax.TokenNewline})...)
		}
	}

	lastSection := &Section{
		Type:     "",
		Labels:   nil,
		Value:    []hclwrite.Tokens{bodyBuffer},
		Comments: commentBuffer,
	}
	if lastSection.IsEmpty() {
		rawSections = append(rawSections, lastSection)
	}
	return rawSections, nil
}

func tokenIsNewline(tok *hclwrite.Token) bool {
	if tok.Type == hclsyntax.TokenNewline {
		return true
	} else if tok.Type == hclsyntax.TokenComment {
		// Single line tokens (# and //) consume their terminating newline,
		// so we need to treat them as newline tokens as well.
		if len(tok.Bytes) > 0 && tok.Bytes[len(tok.Bytes)-1] == '\n' {
			return true
		}
	}
	return false
}

func GetSectionBody(tokens hclwrite.Tokens) hclwrite.Tokens {
	if tokens == nil {
		return nil
	}

	index := 0
	blockStart := 0
	for index < len(tokens) && !tokenIsNewline(tokens[index]) {
		if tokens[index].Type == hclsyntax.TokenEqual || tokens[index].Type == hclsyntax.TokenOBrace {
			blockStart = index + 1
		} else {
			break
		}
		index++
	}

	if blockStart == 0 {
		return tokens
	}

	result := tokens[blockStart:]
	if tokens[blockStart-1].Type == hclsyntax.TokenOBrace {
		opening := tokens[blockStart-1].Type
		p := &parser{
			&peeker{
				Tokens: tokens[blockStart:],
			},
		}
		result = p.ReadTokensUntil([]hclsyntax.TokenType{oppositeBracket(opening)})
		result = result[:len(result)-1]
	}

	for index = 0; index < len(result); index++ {
		if result[index].Type != hclsyntax.TokenNewline {
			result = result[index:]
			break
		}
	}

	return result
}

func tokenBracketChange(tok *hclwrite.Token) int {
	switch tok.Type {
	case hclsyntax.TokenOBrace, hclsyntax.TokenOBrack, hclsyntax.TokenOParen, hclsyntax.TokenTemplateControl, hclsyntax.TokenTemplateInterp:
		return 1
	case hclsyntax.TokenCBrace, hclsyntax.TokenCBrack, hclsyntax.TokenCParen, hclsyntax.TokenTemplateSeqEnd:
		return -1
	default:
		return 0
	}
}

// The following logic is based on the fact that we ignore attributes that start with open parenthesis and treat them
// as one whole blob. If those attributes were to be considered, the following could fail in producing the correct output.
func getHeadlines(tokens hclwrite.Tokens) ([]string, error) {
	var result []string
	if tokens == nil || len(tokens) == 0 {
		return nil, nil
	}

	tokens[0].SpacesBefore = 0
	buffer := hclwrite.Tokens{
		tokens[0],
	}

	index := 1
	for index < len(tokens) {
		if tokens[index].SpacesBefore == 0 {
			buffer = append(buffer, tokens[index])
		} else {
			result = append(result, string(buffer.Bytes()))
			tokens[index].SpacesBefore = 0
			buffer = hclwrite.Tokens{
				tokens[index],
			}
		}
		index++
	}
	if buffer != nil && len(buffer) > 0 {
		result = append(result, string(buffer.Bytes()))
	}
	return result, nil
}
