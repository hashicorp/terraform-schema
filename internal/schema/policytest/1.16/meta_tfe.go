// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func tfeStackMetaAttribute() *schema.AttributeSchema {
	return &schema.AttributeSchema{
		Constraint: schema.Object{
			Attributes: schema.ObjectAttributes{
				"stack_name": &schema.AttributeSchema{
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("The stack name for Terraform Stacks "),
				},
				"deployment_name": &schema.AttributeSchema{
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("The deployment name for Terraform Stacks"),
				},
				"deployment_group": &schema.AttributeSchema{
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("The deployment group for Terraform Stacks"),
				},
			},
		},
		IsOptional:  true,
		Description: lang.Markdown("Terraform Stack the policy set is evaluating for. Empty strings when evaluating workspaces"),
	}
}

func tfeWorkspaceMetaAttribute() *schema.AttributeSchema {
	return &schema.AttributeSchema{
		Constraint: schema.Object{
			Attributes: schema.ObjectAttributes{
				"tags": &schema.AttributeSchema{
					Constraint:  schema.Map{Elem: schema.AnyExpression{OfType: cty.String}},
					IsOptional:  true,
					Description: lang.Markdown("Tags for the current workspace. Empty map when evaluating Stacks or when the workspace has no tags"),
				},
			},
		},
		IsOptional:  true,
		Description: lang.Markdown("Terraform workspace the policy set is evaluating for"),
	}
}
