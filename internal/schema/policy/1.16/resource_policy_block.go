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

func resourcePolicyBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "resource_policy"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "resource_policy",
			ScopeId:              refscope.ResourcePolicyScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "resource_type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Resource Type"),
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Policy Name"),
			},
		},
		Description: lang.Markdown("Defines a policy against resource `name` of type `resource_type`"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"filter": {
					Constraint:  schema.AnyExpression{OfType: cty.Bool},
					IsOptional:  true,
					Description: lang.Markdown("An expression that determines if the policy should be applied to a resource instance. If it evaluates to `false`, the policy is not applied"),
				},
				"enforcement_level": {
					IsOptional:  true,
					Description: lang.Markdown("Defines the strictness of this policy. Determines if a violation allows the run to proceed, requires a manual override, or blocks it entirely."),
					Constraint: schema.OneOf{
						schema.LiteralValue{
							Value:       cty.StringVal("advisory"),
							Description: lang.Markdown("Provides warnings and best practices during the run without blocking progress")},
						schema.LiteralValue{
							Value:       cty.StringVal("mandatory_overridable"),
							Description: lang.Markdown("Blocks the apply stage on failure unless an authorized user manually overrides the requirement")},
						schema.LiteralValue{
							Value:       cty.StringVal("mandatory"),
							Description: lang.Markdown("Immediately halts the run on failure. Requires a configuration fix to proceed; cannot be bypassed")},
					},
				},
				"operations": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.LiteralValue{
								Value:       cty.StringVal("create"),
								Description: lang.Markdown("Apply policy on resource creation"),
							},
							schema.LiteralValue{
								Value:       cty.StringVal("update"),
								Description: lang.Markdown("Apply policy on resource updates"),
							},
							schema.LiteralValue{
								Value:       cty.StringVal("delete"),
								Description: lang.Markdown("Apply policy on resource deletion"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Scopes the policy to a specific subset of planned resource actions. Defaults to `create` and `update` if omitted."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"enforce": enforceBlockNestedSchema(),
				"locals":  localsBlockNestedSchema(),
			},
		},
	}
}
