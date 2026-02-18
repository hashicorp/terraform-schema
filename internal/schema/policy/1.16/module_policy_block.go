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

func modulePolicyBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "module_policy"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "module_policy",
			ScopeId:              refscope.ModulePolicyScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "source",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Module Source"),
				IsDepKey:               true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Policy Name"),
			},
		},
		Description: lang.Markdown("Defines a policy against module `name` of source `source`"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"filter": {
					Constraint:  schema.AnyExpression{OfType: cty.Bool},
					IsOptional:  true,
					Description: lang.Markdown("An expression that determines if the policy should be applied to a module. If it evaluates to `false`, the policy is not applied"),
				},
				"enforcement_level": {
					IsOptional:  true,
					Description: lang.Markdown("Defines the strictness of this policy. Determines if a violation allows the Run to proceed, requires a manual override, or blocks it entirely."),
					Constraint: schema.OneOf{
						schema.LiteralValue{Value: cty.StringVal("mandatory-overridable")},
						schema.LiteralValue{Value: cty.StringVal("mandatory-ov")},
					},
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"enforce": enforceBlockNestedSchema(),
			},
		},
	}
}
