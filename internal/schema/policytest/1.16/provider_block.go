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
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("Defines a 'mock' provider configuration used to test provider_policy rules. It allows you to simulate provider metadata and configuration settings"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"expect_failure": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Expect test to fail"),
				},
				"attrs": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsOptional:  true,
					Description: lang.Markdown("Specify the values that should be returned for specific attributes"),
				},
				"meta": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{
							"name": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Local identifier for the provider"),
							},
							"alias": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Local alias of the provider"),
							},
							"type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The official, short name of the provider. This is the simple identifier used to declare a provider block or resource type"),
							},
							"namespace": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("In the context of the registry, this is the organization or user who publishes the provider. It is the first segment of the provider's source address"),
							},
							"source": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The full, canonical registry address used to locate and download the provider plugin. It combines the namespace and the type"),
							},
							"module_path": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Root module path"),
							},
							"version": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Version of the provider"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Meta attributes of the provider"),
				},
			},
		},
	}
}
