// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func inputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "input"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "input",
			ScopeId:      refscope.InputScope,
			AsReference:  true,
			AsTypeOf: &schema.BlockAsTypeOf{
				AttributeExpr: "type",
			},
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Variable},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Input Name"),
			},
		},
		Description: lang.Markdown("Input allowing users to customize aspects of the policy configuration"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Description to document the purpose of the input and what value is expected"),
				},
				"type": {
					Constraint:  schema.TypeDeclaration{},
					IsRequired:  true,
					Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. `string` or `list(string)`"),
				},
				"nullable": {
					Constraint:   schema.LiteralType{Type: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Specifies whether `null` is a valid value for this input"),
				},
				"sensitive": {
					Constraint:   schema.LiteralType{Type: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Whether the input contains sensitive material and should be hidden in the UI"),
				},
			},
		},
	}
}
