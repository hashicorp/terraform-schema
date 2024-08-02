// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func inputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Variable},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Input Name"),
			},
		},
		Description: lang.Markdown("Stack inputs allow users to customize aspects of the stack that differ between deployments" +
			"(e.g. different instance sizes, regions, etc.)"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Description to document the purpose of the input and what value is expected"),
				},
				"type": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. `string` or `list(string)`"),
				},
				"default": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("A literal expression of an appropriate type for the input"),
				},
			},
		},
	}
}
