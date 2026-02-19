// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func policytestBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		MaxItems:               1,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Policy},
		Description:            lang.Markdown("A top-level block that specifies the policies against which the resources in the test file must be evaluated. If omitted, the framework executes all policies against the file as an integration test."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"targets": {
					Constraint: schema.Set{
						Elem: schema.LiteralType{Type: cty.String},
					},
					IsOptional:  true,
					Description: lang.Markdown("Defines a list of policy files against which the tests will be evaluated."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"plugins": {
					Description: lang.Markdown("Defines the location of custom functions that can be used in the context of that policytest file."),
					MaxItems:    1,
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							Constraint: schema.OneOf{
								schema.Object{
									Attributes: schema.ObjectAttributes{
										"source": &schema.AttributeSchema{
											Constraint:  schema.LiteralType{Type: cty.String},
											IsRequired:  true,
											Description: lang.Markdown("Source where to load the plugin from"),
										},
									},
								},
								schema.LiteralType{Type: cty.String},
							},
						},
					},
				},
			},
		},
	}
}
