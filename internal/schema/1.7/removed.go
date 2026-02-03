// Copyright IBM Corp. 2020, 2026
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
		Description: lang.Markdown("Declaration to specify what resource or module to remove from the state"),
		Body: &schema.BodySchema{
			HoverURL: "https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources",
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
					Description: lang.Markdown("Lifecycle customizations controlling the removal"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"destroy": {
								Constraint:  schema.LiteralType{Type: cty.Bool},
								IsRequired:  true,
								Description: lang.Markdown("Whether Terraform will attempt to destroy the objects (`true`) or not, i.e. just remove from state (`false`)."),
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
