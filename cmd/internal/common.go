package internal

import (
	"math"
)

type PrioritySetting struct {
	Names              []string
	NewLineCountAfter  int
	NewlineCountBefore int
}

type LocationSettings struct {
	InnerIndex         int
	OuterIndex         int
	NewLineCountAfter  int
	NewlineCountBefore int
}

type PriorityLists struct {
	PrependedAttributes []PrioritySetting
	AppendedAttributes  []PrioritySetting
	PrependedBlocks     []PrioritySetting
}

var defaultPriorities = PriorityLists{
	PrependedAttributes: []PrioritySetting{
		{[]string{"source", "version"}, 1, 0},
		{[]string{"count"}, 1, 0},
		{[]string{"for_each"}, 1, 0},
		{[]string{"provider"}, 1, 0},
		{[]string{"providers"}, 1, 0},
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
	AppendedAttributes: []PrioritySetting{
		{[]string{"depends_on"}, 0, 1},
		{[]string{"tags"}, 0, 1},
		{[]string{"lifecycle"}, 0, 1},
	},
	PrependedBlocks: []PrioritySetting{
		{[]string{"terraform"}, 1, 0},
		{[]string{"locals"}, 1, 0},
	},
}

var priorities = map[string]*PriorityLists{
	"root":   &defaultPriorities,
	"module": &defaultPriorities,
	"resource": {
		PrependedAttributes: []PrioritySetting{
			{[]string{"count"}, 1, 0},
			{[]string{"for_each"}, 1, 0},
			{[]string{"provider"}, 1, 0},
		},
		AppendedAttributes: defaultPriorities.AppendedAttributes,
		PrependedBlocks:    nil,
	},
	"variable": {
		PrependedAttributes: []PrioritySetting{{[]string{"type"}, 0, 0}},
	},
	"output": {},
	"data": {
		PrependedAttributes: []PrioritySetting{
			{[]string{"count"}, 1, 0},
			{[]string{"for_each"}, 1, 0},
			{[]string{"provider"}, 1, 0},
		},
	},
	"terraform": {},
	"locals":    {},
	"":          {},
}

const sectionLabel = "terralint"

var rootInternalSectionName = []string{"root"}

func getPriorities(key string) *PriorityLists {
	if _, found := priorities[key]; found {
		return priorities[key]
	}
	return &PriorityLists{}
}

func getLocation(name string, array []PrioritySetting) *LocationSettings {
	for outer, settings := range array {
		for inner, item := range settings.Names {
			if name == item {
				return &LocationSettings{
					InnerIndex:         inner,
					OuterIndex:         outer,
					NewLineCountAfter:  settings.NewLineCountAfter,
					NewlineCountBefore: settings.NewlineCountBefore,
				}
			}
		}
	}
	return &LocationSettings{
		InnerIndex:         math.MaxInt,
		OuterIndex:         math.MaxInt,
		NewLineCountAfter:  0,
		NewlineCountBefore: 0,
	}
}
