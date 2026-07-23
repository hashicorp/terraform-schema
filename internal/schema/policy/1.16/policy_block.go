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
		MaxItems:               1,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Policy},
		Description:            lang.Markdown("The policy block contains high-level configuration for how tfpolicy evaluates a policy, and the conditions Terraform needs to meet to evaluate the policy."),
		Body: &schema.BodySchema{
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
				"required_providers": {
					Description: lang.Markdown("Declares the Terraform providers used by the policy"),
					MaxItems:    1,
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							Description: lang.Markdown("Local label of a Terraform provider, mapped to its source and optional version constraint"),
							Constraint: schema.Object{
								Attributes: schema.ObjectAttributes{
									"source": &schema.AttributeSchema{
										Constraint:  schema.LiteralType{Type: cty.String},
										IsRequired:  true,
										Description: lang.Markdown("The global source address for the provider, e.g. `\"hashicorp/aws\"`"),
									},
									"version": &schema.AttributeSchema{
										Constraint:  schema.LiteralType{Type: cty.String},
										IsOptional:  true,
										Description: lang.Markdown("Version constraint specifying which provider versions this policy is compatible with, e.g. `\">= 5.0.0, < 6.0.0\"`"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
