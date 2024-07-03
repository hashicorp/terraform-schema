// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func outputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Output},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Output Name"),
			},
		},
		Description: lang.PlainText("Output value for consumption by another module or a human interacting via the UI"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.PlainText("Human-readable description of the output (for documentation and UI)"),
				},
				"value": {
					Constraint:  schema.AnyExpression{OfType: cty.DynamicPseudoType},
					IsRequired:  true,
					Description: lang.PlainText("Value, typically a reference to an attribute of a resource or a data source"),
				},
				"type": {
					Constraint:  schema.AnyExpression{OfType: cty.DynamicPseudoType},
					IsRequired:  true,
					Description: lang.PlainText(""), // TODO
				},
			},
		},
	}
}
