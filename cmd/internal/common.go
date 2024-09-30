package internal

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"math"
	"regexp"
	"sort"
	"strings"
	ignorantparser "terralint/cmd/internal/ignorant-parser"
	"terralint/cmd/utilities"
)

type PriorityLists struct {
	PrependedAttributes [][]string
	AppendedAttributes  [][]string
	PrependedBlocks     [][]string
}

type internalSection struct {
	Comments    hclwrite.Tokens
	Sections    []*Section
	Name        []string
	IsMultiLine bool
	Rules       []rule
}

type rule interface {
	Name() string
	Apply(section *ignorantparser.Section) *ignorantparser.Section
}

type NamingRule struct{}

func (namingRule NamingRule) Name() string {
	return "ignore_naming"
}
func (namingRule NamingRule) Apply(section *ignorantparser.Section) *ignorantparser.Section {
	section.Type = strings.ReplaceAll(section.Type, "-", "_")
	return section
}

type Section struct {
	Section *ignorantparser.Section
	Rules   []rule
}

var defaultPriorities = PriorityLists{
	PrependedAttributes: [][]string{
		{"count"},
		{"for_each"},
		{"source", "version"},
		{"provider"},
		{"providers"},
	},
	// This order puts anything with lower index number, closer to the end
	// of the block. For example, the following order means a block like this after apply:
	//
	//	block {
	//	  [...]
	//
	//	  lifecycle {}
	//    tags = {}
	//	  depends_on = []
	//	}
	AppendedAttributes: [][]string{
		{"depends_on"},
		{"tags"},
		{"lifecycle"},
	},
	PrependedBlocks: [][]string{
		{"terraform"},
		{"locals"},
	},
}

var priorities = map[string]*PriorityLists{
	"root":   &defaultPriorities,
	"module": &defaultPriorities,
	"resource": {
		PrependedAttributes: [][]string{
			{"count"},
			{"for_each"},
			{"provider"},
		},
		AppendedAttributes: defaultPriorities.AppendedAttributes,
		PrependedBlocks:    nil,
	},
	"variable": {
		PrependedAttributes: [][]string{{"type"}},
	},
	"output":    {},
	"data":      {},
	"terraform": {},
	"locals":    {},
	"":          {},
}

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

var tokenComma = hclwrite.Token{
	Type:         hclsyntax.TokenComma,
	Bytes:        []byte(","),
	SpacesBefore: 0,
}

const sectionLabel = "terralint"

var internalSectionStartRegex = regexp.MustCompile(fmt.Sprintf("^#\\s*%s(?:\\s+[^\\{]+)?\\s*\\{?$", sectionLabel))
var internalSectionEndingRegex = regexp.MustCompile("^\\s*}$")
var internalSectionNameRegex = regexp.MustCompile(fmt.Sprintf("^#\\s*%s(?:((\\s+[^\\{]+)*))?\\s*\\{?$", sectionLabel))
var sectionEndType = utilities.RandString(10)

var rootInternalSectionName = []string{"root"}

func defaultRules() ([]rule, map[string]rule) {
	ruleArray := []rule{NamingRule{}}
	ruleMap := make(map[string]rule)

	for _, value := range ruleArray {
		ruleMap[value.Name()] = value
	}

	return ruleArray, ruleMap
}

func getPriorities(key string) *PriorityLists {
	if _, found := priorities[key]; found {
		return priorities[key]
	}
	return &PriorityLists{}
}

func countChar(s string, c rune) int {
	count := 0
	for _, r := range s {
		if r == c {
			count++
		}
	}
	return count
}

func breakStartingSection(section *ignorantparser.Section) *internalSection {
	for commentIndex, comment := range section.Comments {
		commentString := strings.TrimSpace(string(comment.Bytes))

		if internalSectionStartRegex.MatchString(commentString) {
			var result internalSection

			result.IsMultiLine = strings.HasSuffix(commentString, "{")

			if result.IsMultiLine {
				result.Comments = section.Comments[:commentIndex]
				section.Comments = section.Comments[commentIndex+1:]
			}

			match := internalSectionNameRegex.FindStringSubmatch(commentString)
			var name []string
			if len(match) > 1 {
				name = strings.Fields(strings.TrimSpace(match[1]))
				sort.SliceStable(name, func(i, j int) bool {
					return name[i] < name[j]
				})
			}
			result.Name = name

			_, ruleArray := defaultRules()
			for _, r := range ruleArray {
				if !utilities.Exists(r.Name(), result.Name) {
					result.Rules = append(result.Rules, r)
				}
			}

			result.Sections = append(result.Sections, &Section{
				Section: section,
				Rules:   result.Rules,
			})
			return &result
		}
	}
	return nil
}

func endingBraceIndex(section *ignorantparser.Section) int {
	openingCount := 0
	for index, comment := range section.Comments {
		openingCount += countChar(string(comment.Bytes), '{')
		openingCount -= countChar(string(comment.Bytes), '}')
		if openingCount < 0 {
			return index
		}
	}
	return -1
}

func getInternalSections(sections []*ignorantparser.Section, parentRules []rule) []*internalSection {
	var groups []*internalSection
	root := &internalSection{
		Comments: nil,
		Sections: nil,
		Name:     rootInternalSectionName,
		Rules:    parentRules,
	}

	index := 0
	inSection := false

	for index < len(sections) {
		if !sections[index].HasValue() {
			index++
			continue
		}
		startingSection := breakStartingSection(sections[index])
		if startingSection == nil {
			var rules []rule
			if groups == nil {
				rules = root.Rules
			} else {
				rules = groups[len(groups)-1].Rules
			}

			endingBrace := endingBraceIndex(sections[index])

			if endingBrace != -1 {
				endingComments := &ignorantparser.Section{
					Comments: sections[index].Comments[:endingBrace],
				}

				if endingComments.HasValue() {
					groups[len(groups)-1].Sections = append(groups[len(groups)-1].Sections, &Section{
						Section: endingComments,
						Rules:   rules,
					})
				}

				inSection = false
				sections[index].Comments = sections[index].Comments[endingBrace+1:]
				continue
			} else {
				switch inSection {
				case true:
					groups[len(groups)-1].Sections = append(groups[len(groups)-1].Sections, &Section{
						Section: sections[index],
						Rules:   rules,
					})
				case false:
					root.Sections = append(root.Sections, &Section{
						Section: sections[index],
						Rules:   root.Rules,
					})
				}
			}
		} else if !startingSection.IsMultiLine {
			root.Sections = append(root.Sections, startingSection.Sections[0])
		} else {
			groups = append(groups, startingSection)
			inSection = startingSection.IsMultiLine
		}

		index++
	}
	if root.Sections != nil && len(root.Sections) > 0 {
		groups = append([]*internalSection{root}, groups...)
	}
	return groups
}

func getFormattedContent(filePath string) ([]byte, error) {
	sections, err := ignorantparser.ParseConfigFromFile(filePath)
	if err != nil {
		return nil, err
	}

	for index := range sections {
		standardizeCommentsSection, err := standardizeComments(&Section{
			Section: sections[index],
			Rules:   nil,
		})
		if err != nil {
			return nil, err
		}
		sections[index] = standardizeCommentsSection.Section
	}

	sections = mergeLocals(sections)

	defaultRulesArray, _ := defaultRules()
	sections, err = applyRules(sections, "root", defaultRulesArray)
	if err != nil {
		return nil, err
	}

	formattedBytes := hclwrite.Format(getToken(sections).Bytes())

	return formattedBytes, nil
}

func mergeLocals(input []*ignorantparser.Section) []*ignorantparser.Section {
	var merged []*ignorantparser.Section
	locals := &ignorantparser.Section{
		Type:     "locals",
		Labels:   nil,
		Value:    []hclwrite.Tokens{{}},
		Comments: nil,
	}

	for _, section := range input {
		if section.Type == "locals" {
			locals.Value[0] = append(locals.Value[0], ignorantparser.GetSectionBody(section.FlattenValue())...)
			locals.Comments = append(locals.Comments, section.Comments...)
		} else {
			merged = append(merged, section)
		}
	}

	if locals.Value != nil {
		locals.Value[0] = append(
			append(hclwrite.Tokens{
				&tokenOBrace,
				&tokenNewLine,
			}, locals.FlattenValue()...),
			hclwrite.Tokens{
				&tokenCBrace,
				&tokenNewLine,
			}...)
	}
	if locals.Value != nil || locals.Comments != nil {
		merged = append(merged, locals)
	}
	return merged
}

func getToken(sections []*ignorantparser.Section) hclwrite.Tokens {
	file := hclwrite.NewEmptyFile()
	for index, section := range sections {
		if !section.HasValue() {
			continue
		}
		file.Body().AppendUnstructuredTokens(section.Tokens())

		if index == len(sections)-1 || (index < len(sections)-1 && sections[index+1].Type == sectionEndType) {
			continue
		}
		file.Body().AppendNewline()
	}
	tokens := file.BuildTokens(nil)

	if tokens == nil || len(tokens) == 0 {
		return tokens
	}

	return tokens
}

func applyRules(sections []*ignorantparser.Section, parentType string, parentRules []rule) ([]*ignorantparser.Section, error) {
	internalSections := getInternalSections(sections, parentRules)
	var result []*ignorantparser.Section
	for index := range internalSections {
		sort.SliceStable(internalSections[index].Sections, func(i, j int) bool {
			return sortLogic(internalSections[index].Sections, i, j, getPriorities(parentType))
		})

		for subIndex := range internalSections[index].Sections {
			subsection := internalSections[index].Sections[subIndex]
			if !subsection.Section.HasValue() {
				continue
			}

			subsection, err := standardizeComments(subsection)
			if err != nil {
				return nil, err
			}

			for _, r := range subsection.Rules {
				subsection.Section = r.Apply(subsection.Section)
			}

			if subsection.Section.LineCounts() == 1 || (!subsection.Section.IsBlock() && !subsection.Section.IsList()) {
				continue
			}

			var tokens hclwrite.Tokens
			if subsection.Section.IsList() {
				sub := subsection.Section.Value[1 : len(subsection.Section.Value)-1]

				var ts []hclwrite.Tokens
				for _, v := range sub {
					var tsItem hclwrite.Tokens
					innerSections, err := ignorantparser.ParseSectionConfig(v)
					if err != nil {
						return nil, err
					}

					innerSections, err = applyRules(innerSections, subsection.Section.Type, internalSections[index].Rules)
					if len(innerSections) == 0 {
						ts = append(ts, hclwrite.Tokens{
							&tokenOBrace,
							&tokenCBrace,
						})
						continue
					}

					obj := innerSections[0].Type != ""
					if obj {
						tsItem = hclwrite.Tokens{
							&tokenOBrace,
							&tokenNewLine,
						}
					}
					for _, si := range innerSections {
						end := len(si.Value[0]) - 1
						if obj {
							tsItem = append(tsItem, hclwrite.TokensForIdentifier(si.Type)...)
							end++
						}
						tsItem = append(tsItem, si.Value[0][:end]...)
					}

					if obj {
						tsItem = append(tsItem, &tokenCBrace)
					}
					ts = append(ts, tsItem)
				}
				var r hclwrite.Tokens
				for _, t := range ts {
					r = append(r, t...)
				}

				tokens = hclwrite.Tokens{&tokenEqual}
				if subsection.Section.ListCount() <= 1 {
					closing := tokenCBrack
					closing.SpacesBefore = 1
					tokens = append(tokens, &tokenOBrack)
					if subsection.Section.ListCount() == 1 {
						tokens = append(tokens, ts[0]...)
					}
					tokens = append(tokens, &closing, &tokenNewLine)
				} else {
					tokens = append(tokens, hclwrite.Tokens{
						&tokenOBrack,
						&tokenNewLine,
					}...)
					for _, val := range ts {
						tokens = append(tokens, val...)
						tokens = append(tokens, &tokenComma, &tokenNewLine)
					}
					tokens = append(tokens, &tokenCBrack)
				}
			} else {
				innerSections, err := ignorantparser.ParseSectionConfig(subsection.Section.Value[0])
				if err != nil {
					return nil, err
				}

				if subsection.Section.IsAttribute() {
					tokens = append(tokens, &tokenEqual)
				}
				if len(innerSections) == 0 {
					tokens = append(tokens, hclwrite.Tokens{&tokenOBrace, &tokenCBrace}...)
					subsection.Section.Value = []hclwrite.Tokens{tokens}
					continue
				}

				if subsection.Section.IsBlock() {
					tokens = append(tokens, hclwrite.Tokens{&tokenOBrace, &tokenNewLine}...)
				}

				innerSections, err = applyRules(innerSections, subsection.Section.Type, internalSections[index].Rules)
				if err != nil {
					return nil, err
				}

				isPreviousBlock := false
				innerSectionsIndex := 0
				previousOuterIndex := math.MaxInt
				for innerSectionsIndex < len(innerSections) {
					currentOuterIndex, _ := getLocation(innerSections[innerSectionsIndex].Type, getPriorities(subsection.Section.Type).PrependedAttributes)
					isPreviousPrependedAttribute := previousOuterIndex != math.MaxInt && currentOuterIndex == math.MaxInt

					isCurrentBlock := innerSections[innerSectionsIndex].LineCounts() > 1
					if (isCurrentBlock || isPreviousBlock || isPreviousPrependedAttribute) && innerSectionsIndex > 0 {
						tokens = append(tokens, &tokenNewLine)
					}
					tokens = append(tokens, innerSections[innerSectionsIndex].Tokens()...)

					previousOuterIndex = currentOuterIndex

					isPreviousBlock = isCurrentBlock
					if isPreviousBlock && innerSectionsIndex < len(innerSections)-1 && innerSections[innerSectionsIndex+1].Type == sectionEndType {
						innerSectionsIndex++
						tokens = append(tokens, innerSections[innerSectionsIndex].Tokens()...)
					}
					innerSectionsIndex++
				}

				if subsection.Section.IsBlock() {
					tokens = append(tokens, hclwrite.Tokens{&tokenCBrace, &tokenNewLine}...)
				}
			}

			subsection.Section.Value = []hclwrite.Tokens{tokens}
		}

		if !utilities.Exists("root", internalSections[index].Name) {
			sectionName := strings.Join(internalSections[index].Name, " ")
			if sectionName == "" {
				sectionName = "{"
			} else {
				sectionName += " {"
			}

			sectionComments := fmt.Sprintf("# %s %s\n", sectionLabel, sectionName)
			if internalSections[index].Comments != nil && len(internalSections[index].Comments) > 0 {
				sectionComments = fmt.Sprintf("%s%s", string(internalSections[index].Comments.Bytes()), sectionComments)
			}

			internalSections[index].Sections[0].Section.Comments = append(hclwrite.Tokens{
				&hclwrite.Token{
					Type:         hclsyntax.TokenComment,
					Bytes:        []byte(sectionComments),
					SpacesBefore: 0,
				},
			}, internalSections[index].Sections[0].Section.Comments...)
			internalSections[index].Sections = append(internalSections[index].Sections, &Section{
				Section: &ignorantparser.Section{

					Type:   sectionEndType,
					Labels: nil,
					Value:  nil,
					Comments: hclwrite.Tokens{
						&hclwrite.Token{
							Type:         hclsyntax.TokenComment,
							Bytes:        []byte("# }\n"),
							SpacesBefore: 0,
						},
					},
				},
				Rules: nil,
			})
		}

		for _, section := range internalSections[index].Sections {
			result = append(result, section.Section)
		}
	}
	return result, nil
}

func getLocation(name string, array [][]string) (int, int) {
	for outer, items := range array {
		for inner, item := range items {
			if name == item {
				return outer, inner
			}
		}
	}
	return math.MaxInt, math.MaxInt
}

func sortLogic(sections []*Section, first, second int, priorities *PriorityLists) bool {
	firstPrependedLocationOuter, firstPrependedLocationInner := getLocation(sections[first].Section.Type, priorities.PrependedBlocks)
	secondPrependedLocationOuter, secondPrependedLocationInner := getLocation(sections[second].Section.Type, priorities.PrependedBlocks)

	if firstPrependedLocationOuter != secondPrependedLocationOuter {
		return firstPrependedLocationOuter < secondPrependedLocationOuter
	} else if firstPrependedLocationOuter == secondPrependedLocationOuter && firstPrependedLocationOuter != math.MaxInt {
		return firstPrependedLocationInner < secondPrependedLocationInner
	}

	firstPrependedLocationOuter, firstPrependedLocationInner = getLocation(sections[first].Section.Type, priorities.PrependedAttributes)
	secondPrependedLocationOuter, secondPrependedLocationInner = getLocation(sections[second].Section.Type, priorities.PrependedAttributes)

	if firstPrependedLocationOuter != secondPrependedLocationOuter {
		return firstPrependedLocationOuter < secondPrependedLocationOuter
	} else if firstPrependedLocationOuter == secondPrependedLocationOuter && firstPrependedLocationOuter != math.MaxInt {
		return firstPrependedLocationInner < secondPrependedLocationInner
	}

	firstAppendedLocationOuter, firstAppendedLocationInner := getLocation(sections[first].Section.Type, priorities.AppendedAttributes)
	secondAppendedLocationOuter, secondAppendedLocationInner := getLocation(sections[second].Section.Type, priorities.AppendedAttributes)

	if firstAppendedLocationOuter != secondAppendedLocationOuter {
		return firstAppendedLocationOuter > secondAppendedLocationOuter
	} else if firstAppendedLocationOuter == secondAppendedLocationOuter && firstAppendedLocationOuter != math.MaxInt {
		return firstAppendedLocationInner > secondAppendedLocationInner
	}

	if sections[first].Section.IsAttribute() != sections[second].Section.IsAttribute() {
		return sections[first].Section.IsAttribute()
	}

	if !sections[first].Section.IsAttribute() && !sections[second].Section.IsAttribute() {
		if sections[first].Section.Type == "dynamic" && sections[second].Section.Type != "dynamic" {
			return false
		}

		if sections[second].Section.Type == "dynamic" && sections[first].Section.Type != "dynamic" {
			return true
		}
	}

	// Handling comment only sections
	if sections[first].Section.Value == nil {
		return false
	} else if sections[second].Section.Value == nil {
		return true
	}

	lowerFirstType := strings.ToLower(sections[first].Section.Type)
	lowerSecondType := strings.ToLower(sections[second].Section.Type)
	if lowerFirstType != lowerSecondType {
		return lowerFirstType < lowerSecondType
	}

	firstHead := strings.Join(sections[first].Section.Labels, " ")
	secondHead := strings.Join(sections[second].Section.Labels, " ")
	return strings.ToLower(firstHead) < strings.ToLower(secondHead)
}

func standardizeComments(section *Section) (*Section, error) {
	if section.Section.Comments == nil || len(section.Section.Comments) == 0 {
		return section, nil
	}

	unconventionalComment := regexp.MustCompile("^\\/{0,1}\\s*\\**\\/{0,1}\\s*")
	var updatedComments hclwrite.Tokens
	splitted := strings.Split(string(section.Section.Comments.Bytes()), "\n")
	for index, commentLine := range splitted {
		commentLine = strings.TrimSpace(commentLine)
		if unconventionalComment.FindString(commentLine) != "" {
			commentLine = unconventionalComment.ReplaceAllString(commentLine, "# ")
			if commentLine == "# " {
				continue
			}
		}

		if strings.TrimSpace(commentLine) == "" {
			if index == len(splitted)-1 {
				continue
			}
			commentLine = "#"
		}

		updatedComments = append(updatedComments, &hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(fmt.Sprintf("%s\n", commentLine)),
			SpacesBefore: 0,
		})
	}

	section.Section.Comments = updatedComments
	return section, nil
}
