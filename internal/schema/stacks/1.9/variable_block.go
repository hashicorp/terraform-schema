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

func variableBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "var"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "variable",
			ScopeId:      refscope.VariableScope,
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
				Description:            lang.PlainText("Variable Name"),
			},
		},
		Description: lang.Markdown("Stack variable allowing users to customize aspects of the stack that differ between deployments" +
			"(e.g. different instance sizes, regions, etc.)"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Description to document the purpose of the variable and what value is expected"),
				},
				"sensitive": {
					Constraint:   schema.LiteralType{Type: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Whether the variable contains sensitive material and should be hidden in the UI"),
				},
				"ephemeral": {
					Constraint:  schema.LiteralType{Type: cty.Bool},
					IsOptional:  true,
					Description: lang.PlainText("Whether the value is ephemeral and should not be persisted in the state"),
				},
				"type": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. `string` or `list(string)`"),
				},
				"default": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("A literal expression of an appropriate type for the variable"),
				},
			},
		},
	}
}
