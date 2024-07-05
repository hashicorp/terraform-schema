// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func orchestrateBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Defines an orchestration rule, such as a rule for when to auto-approve one or more deployments in the stack to be evaluated after a plan or apply operation. These rules allow you to define the behavior of various aspects of the stack in code, and make managing large numbers of deployments more manageable. The block labels include the rule type and the rule name, which together must be unique within the stack"),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Rule Type"),
				IsDepKey:               true,
				Completable:            true,
				// TODO: auto_approve is the only one supported now, but converged, replan, rollout are possible values for the first label on the block
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Rule Name"),
			},
		},
		Body: &schema.BodySchema{
			// TODO proper constraints for these
			Blocks: map[string]*schema.BlockSchema{
				"check": {
					Description: lang.Markdown("Each rule has one or more check blocks, which must all pass in order for the rule to execute its action. The check block follows Terraform’s custom conditions concept pattern, and includes expressions for condition and error_message. These are evaluated against the orchestration context below"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Description: lang.Markdown("The condition must evaluate to true or false"),
								Constraint:  schema.LiteralType{Type: cty.Bool},
							},
							"error_message": {
								Description: lang.Markdown("The error message must be a string"),
								Constraint:  schema.LiteralType{Type: cty.String},
							},
						},
					},
				},
			},
		},
	}
}
