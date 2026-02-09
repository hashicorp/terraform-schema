// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func policyBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Policy},
		Description:            lang.Markdown("The policy block contains high-level configuration for how tfpolicy evaluates a policy, and the conditions Terraform needs to meet to evaluate the policy."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"enforcement_level": {
					IsOptional:  true,
					Description: lang.Markdown("Describes how 'strict' the policy is. It determines whether a failure merely warns the user or strictly halts the infrastructure run"),
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "advisory",
							Description: lang.Markdown("Informational only, no enforcement"),
						},
						schema.Keyword{
							Keyword:     "soft-mandatory",
							Description: lang.Markdown("Enforced with acknowledgment or justification"),
						},
						schema.Keyword{
							Keyword:     "hard-mandatory",
							Description: lang.Markdown("Strictly enforced, blocking"),
						},
					},
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"terraform_config": {
					Description: lang.Markdown("Defines a configuration which is specific to Terraform"),
					MaxItems:    1,
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"required_version": {
								Constraint: schema.LiteralType{Type: cty.String},
								IsOptional: true,
								Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
									"with this configuration, e.g. `~> 1.16`"),
							},
						},
					},
				},
				"plugins": {
					Description: lang.Markdown("Defines the location of custom functions that can be used in the context of that policy file."),
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
