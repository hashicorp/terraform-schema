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
									"with this configuration, e.g. `~> 0.12`"),
							},
						},
					},
				},
			},
		},
	}
}
