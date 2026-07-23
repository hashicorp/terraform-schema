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

func providerBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "provider",
			ScopeId:              refscope.ProviderScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "provider_type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Type"),
			},
			{
				Name:                   "test_case_name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Test Case Name"),
			},
		},
		Description: lang.PlainText("Used to validate `provider_policy` blocks by mocking provider configuration and metadata"),
		Body: &schema.BodySchema{
			Blocks: map[string]*schema.BlockSchema{
				"inputs": inputsNestedBlockSchema(),
			},
			Attributes: map[string]*schema.AttributeSchema{
				"expect_failure": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Anticipates a policy failure based on provider source, version or config"),
				},
				"attrs": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsOptional:  true,
					Description: lang.Markdown("Mocks the actual provider configuration (e.g., `region`, `alias`)"),
				},
				"meta": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{
							"source": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The full source of the provider"),
							},
							"version": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The resolved version of the provider"),
							},
							"alias": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Alias given to the provider"),
							},
							"name": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The local name of the provider"),
							},
							"namespace": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The provider's registry namespace"),
							},

							"type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The official, short name of the provider. This is the simple identifier used to declare a provider block or resource type"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Mocks the `required_providers` information"),
				},
			},
		},
	}
}
