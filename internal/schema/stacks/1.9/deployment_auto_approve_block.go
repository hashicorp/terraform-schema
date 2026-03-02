// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func deploymentAutoApproveBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Defines rules that automatically approve deployment plans when specific conditions are met. A deployment_group block can reference one or more deployment_auto_approve blocks. All checks within the deployment_auto_approve blocks must pass for the plans of a deployment group to automatically apply"),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "deployment_auto_approve"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "deployment_auto_approve",
			ScopeId:      refscope.DeploymentAutoApproveScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Auto-Approve Rule Name"),
			},
		},
		Body: &schema.BodySchema{
			Blocks: map[string]*schema.BlockSchema{
				"check": {
					Description: lang.Markdown("A check block contains conditions that must be met for the deployment_auto_approve block to automatically approve plans"),
					MinItems:    1,
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Description: lang.Markdown("Expression that Terraform evaluates."),
								Constraint:  schema.AnyExpression{OfType: cty.Bool},
								IsRequired:  true,
							},
							"reason": {
								Description: lang.Markdown("Message to display if the condition evaluates to `false`. The error message is displayed in HCP Terraform when manual approval is required"),
								Constraint:  schema.LiteralType{Type: cty.String},
								IsRequired:  true,
							},
						},
					},
				},
			},
		},
	}
}
