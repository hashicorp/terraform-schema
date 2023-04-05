// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import "github.com/hashicorp/hcl-lang/schema"

func objectConstraintFromBodySchema(bs *schema.BodySchema) schema.Object {
	if bs == nil {
		return schema.Object{}
	}

	oe := schema.Object{
		Description: bs.Description,
		Attributes:  bs.Attributes,
	}

	for bType, block := range bs.Blocks {
		oe.Attributes[bType] = &schema.AttributeSchema{
			Description:  block.Description,
			IsDeprecated: block.IsDeprecated,
		}

		if block.MinItems > 0 {
			oe.Attributes[bType].IsRequired = true
		} else {
			oe.Attributes[bType].IsOptional = true
		}

		switch block.Type {
		case schema.BlockTypeObject:
			oe.Attributes[bType].Constraint = objectConstraintFromBodySchema(block.Body)
		case schema.BlockTypeList:
			oe.Attributes[bType].Constraint = schema.List{
				Elem:     objectConstraintFromBodySchema(block.Body),
				MinItems: block.MinItems,
				MaxItems: block.MaxItems,
			}
		case schema.BlockTypeSet:
			oe.Attributes[bType].Constraint = schema.Set{
				Elem:     objectConstraintFromBodySchema(block.Body),
				MinItems: block.MinItems,
				MaxItems: block.MaxItems,
			}
		case schema.BlockTypeMap:
			oe.Attributes[bType].Constraint = schema.Map{
				Elem:     objectConstraintFromBodySchema(block.Body),
				MinItems: block.MinItems,
				MaxItems: block.MaxItems,
			}
		}
	}

	return oe
}
