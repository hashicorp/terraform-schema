// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func removedBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Declaration to specify what address to remove from the state"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"from": {
					Constraint: schema.OneOf{
						schema.Reference{OfScopeId: refscope.ModuleScope},
						schema.Reference{OfScopeId: refscope.ResourceScope},
					},
					IsRequired:  true,
					Description: lang.Markdown("Address of the module or resource to be removed"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle": {
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"destroy": {
								Constraint: schema.LiteralType{Type: cty.Bool},
								IsRequired: true,
							},
						},
					},
					MinItems: 1,
					MaxItems: 1,
				},
			},
		},
	}
}
