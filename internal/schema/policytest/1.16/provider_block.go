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
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "test_case_name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Test Case Name"),
			},
		},
		Description: lang.PlainText("Used to validate `provider_policy` blocks by mocking provider configuration and metadata"),
		Body: &schema.BodySchema{
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
							"name": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Local identifier for the provider"),
							},
							"source": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The full, canonical registry address used to locate and download the provider plugin. It combines the namespace and the type"),
							},
							"version": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Version of the provider"),
							},
						},
					},
					IsRequired:  true,
					Description: lang.Markdown("Mocks the `required_providers` information"),
				},
			},
		},
	}
}
