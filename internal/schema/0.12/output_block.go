// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func outputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "output"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "output",
			ScopeId:      refscope.OutputScope,
			AsReference:  true,
		},
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
				"sensitive": {
					Constraint:  schema.LiteralType{Type: cty.Bool},
					IsOptional:  true,
					Description: lang.PlainText("Whether the output contains sensitive material and should be hidden in the UI"),
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.ModuleScope},
							schema.Reference{OfScopeId: refscope.ResourceScope},
						},
					},
					IsOptional:  true,
					Description: lang.PlainText("Set of references to hidden dependencies (e.g. resources or data sources)"),
				},
			},
		},
	}
}
