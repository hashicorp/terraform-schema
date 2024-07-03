// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func variableBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
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
				"default": {
					Constraint:  schema.TypeDeclaration{},
					IsOptional:  true,
					Description: lang.Markdown("A literal expression of an appropriate type for the variable"),
				},
			},
		},
	}
}
