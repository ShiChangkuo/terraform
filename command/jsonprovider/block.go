package jsonprovider

import (
	"github.com/hashicorp/terraform/configs/configschema"
)

type block struct {
	Attributes      map[string]*attribute `json:"attributes,omitempty"`
	BlockTypes      map[string]*blockType `json:"block_types,omitempty"`
	Description     string                `json:"description,omitempty"`
	DescriptionKind string                `json:"description_kind,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
}

type blockType struct {
	NestingMode string `json:"nesting_mode,omitempty"`
	Block       *block `json:"block,omitempty"`
	MinItems    uint64 `json:"min_items,omitempty"`
	MaxItems    uint64 `json:"max_items,omitempty"`
}

func marshalBlockTypes(nestedBlock *configschema.NestedBlock) *blockType {
	if nestedBlock == nil {
		return &blockType{}
	}
	block := marshalBlock(&nestedBlock.Block)
	if block == nil {
		return nil
	}

	ret := &blockType{
		Block:    block,
		MinItems: uint64(nestedBlock.MinItems),
		MaxItems: uint64(nestedBlock.MaxItems),
	}

	switch nestedBlock.Nesting {
	case configschema.NestingSingle:
		ret.NestingMode = "single"
	case configschema.NestingGroup:
		ret.NestingMode = "group"
	case configschema.NestingList:
		ret.NestingMode = "list"
	case configschema.NestingSet:
		ret.NestingMode = "set"
	case configschema.NestingMap:
		ret.NestingMode = "map"
	default:
		ret.NestingMode = "invalid"
	}
	return ret
}

func marshalBlock(configBlock *configschema.Block) *block {
	if configBlock == nil {
		return &block{}
	}
	if configBlock.Deprecated {
		return nil
	}

	ret := block{
		Description: configBlock.Description,
	}

	if len(configBlock.Attributes) > 0 {
		attrs := make(map[string]*attribute, len(configBlock.Attributes))
		for k, attr := range configBlock.Attributes {
			if !attr.Deprecated {
				if k == "id" || k == "region" || k == "tenant_id" {
					attr.Optional = false
				}
				attrs[k] = marshalAttribute(attr)
			}
		}
		ret.Attributes = attrs
	}

	// discard "timeouts" block
	delete(configBlock.BlockTypes, "timeouts")
	if len(configBlock.BlockTypes) > 0 {
		blockTypes := make(map[string]*blockType, len(configBlock.BlockTypes))
		for k, bt := range configBlock.BlockTypes {
			if blockType := marshalBlockTypes(bt); blockType != nil {
				blockTypes[k] = blockType
			}
		}
		ret.BlockTypes = blockTypes
	}

	return &ret
}
