// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func deploymentBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Deployment"),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Deplyoment Name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"variables": {
					Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
					IsOptional:  true,
					Constraint: schema.Map{
						Name: "map of variable references",
						Elem: schema.Reference{OfScopeId: refscope.VariableScope},
					},
				},
			},
		},
	}
}
