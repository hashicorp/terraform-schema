// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Blocks["cloud"] = &schema.BlockSchema{
		Description: lang.PlainText("Terraform Cloud configuration"),
		MaxItems:    1,
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"hostname": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The Terraform Enterprise hostname to connect to. " +
						"This optional argument defaults to `app.terraform.io` for use with Terraform Cloud."),
				},
				"organization": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.PlainText("The name of the organization containing the targeted workspace(s)."),
				},
				"token": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The token used to authenticate with Terraform Cloud/Enterprise. " +
						"Typically this argument should not be set, and `terraform login` used instead; " +
						"your credentials will then be fetched from your CLI configuration file " +
						"or configured credential helper."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"workspaces": {
					Description: lang.Markdown("Workspace mapping strategy, either workspace `tags` or `name` is required."),
					MaxItems:    1,
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"name": {
								Constraint: schema.LiteralType{Type: cty.String},
								IsOptional: true,
								Description: lang.Markdown("The name of a single Terraform Cloud workspace " +
									"to be used with this configuration When configured only the specified workspace " +
									"can be used. This option conflicts with `tags`."),
							},
							"tags": {
								Constraint: schema.Set{
									Elem: schema.LiteralType{Type: cty.String},
								},
								IsOptional: true,
								Description: lang.Markdown("A set of tags used to select remote Terraform Cloud workspaces" +
									" to be used for this single configuration.  New workspaces will automatically be tagged " +
									"with these tag values.  Generally, this is the primary and recommended strategy to use. " +
									"This option conflicts with `name`."),
							},
						},
					},
				},
			},
		},
	}

	return bs
}
