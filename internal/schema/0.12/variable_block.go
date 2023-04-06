// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func variableBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "var"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "variable",
			ScopeId:      refscope.VariableScope,
			AsReference:  true,
			AsTypeOf: &schema.BlockAsTypeOf{
				AttributeExpr:  "type",
				AttributeValue: "default",
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
		Description: lang.Markdown("Input variable allowing users to customizate aspects of the configuration when used directly " +
			"(e.g. via CLI, `tfvars` file or via environment variables), or as a module (via `module` arguments)"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Description to document the purpose of the variable and what value is expected"),
				},
				"type": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. `string` or `list(string)`"),
				},
			},
			Blocks: make(map[string]*schema.BlockSchema, 0),
		},
	}

	if v.GreaterThanOrEqual(v0_12_20) {
		bs.Body.Blocks["validation"] = &schema.BlockSchema{
			Description: lang.Markdown("Custom validation rule to restrict what value is expected for the variable"),
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"condition": {
						Constraint: schema.LiteralType{Type: cty.Bool},
						IsRequired: true,
						Description: lang.Markdown("Condition under which a variable value is valid, " +
							"e.g. `length(var.example) >= 4` enforces minimum of 4 characters"),
					},
					"error_message": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsRequired: true,
						Description: lang.Markdown("Error message to present when the variable is considered invalid, " +
							"i.e. when `condition` evaluates to `false`"),
					},
				},
			},
		}
	}

	return bs
}
