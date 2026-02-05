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

func providerPolicyBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "provider_policy"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "provider_policy",
			ScopeId:              refscope.ProviderPolicyScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Type"),
				IsDepKey:               true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Policy Name"),
			},
		},
		Description: lang.Markdown("Defines a policy against provider `name` of type `type`"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"filter": {
					Constraint:  schema.AnyExpression{OfType: cty.Bool},
					IsOptional:  true,
					Description: lang.Markdown("An expression that determines if the policy should be applied to a provider. If it evaluates to `false`, the policy is not applied"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"locals": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Locals},
					Description:            lang.Markdown("Local values to be used in the scope"),
					Body:                   &schema.BodySchema{
						// AnyAttribute: &schema.AttributeSchema{
						// 	Address: &schema.AttributeAddrSchema{
						// 		Steps: []schema.AddrStep{
						// 			schema.StaticStep{Name: "local"},
						// 			schema.AttrNameStep{},
						// 		},
						// 		ScopeId:     refscope.ProviderPolicyScope,
						// 		AsExprType:  true,
						// 		AsReference: true,
						// 	},
						// 	Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
						// 	// Constraint: schema.OneOf{
						// 	// 	schema.Reference{OfType: cty.DynamicPseudoType, OfScopeId: refscope.ProviderPolicyScope},
						// 	// },
						// },
					},
					MaxItems: 1,
				},
				"enforce": enforceBlockNestedSchema(),
			},
		},
	}
}
