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

func publishOutputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("A publish_output block defines an output value that is published from a deployment for consumption by other stacks"),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Published Output Name"),
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"description": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsOptional:  true,
					Description: lang.PlainText("Human-readable description of the published output"),
				},
				"value": {
					Constraint:  schema.Reference{OfScopeId: refscope.DeploymentScope},
					IsRequired:  true,
					Description: lang.PlainText("The value to publish, typically a reference to a deployment output"),
				},
			},
		},
	}
}
