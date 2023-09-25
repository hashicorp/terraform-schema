// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Blocks["cloud"].Body.Blocks["workspaces"].Body = &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"name": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The name of a single Terraform Cloud workspace " +
					"to be used with this configuration. When configured only the specified workspace " +
					"can be used. This option conflicts with `tags`."),
			},
			"project": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.PlainText("The name of a project that resulting workspace(s) will be created in."),
			},
			"tags": {
				Constraint: schema.Set{
					Elem: schema.LiteralType{Type: cty.String},
				},
				IsOptional: true,
				Description: lang.Markdown("A set of tags used to select remote Terraform Cloud workspaces" +
					" to be used for this single configuration. New workspaces will automatically be tagged " +
					"with these tag values. Generally, this is the primary and recommended strategy to use. " +
					"This option conflicts with `name`."),
			},
		},
	}

	return bs
}
