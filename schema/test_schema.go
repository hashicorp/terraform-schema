// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func variablesBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Address: &schema.AttributeAddrSchema{
					Steps: []schema.AddrStep{
						schema.StaticStep{Name: "var"},
						schema.AttrNameStep{},
					},
					ScopeId:    refscope.VariableScope,
					AsExprType: true,
				},
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
			},
		},
	}
}

func providerBlockSchema(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.AttrValueStep{Name: "alias", IsOptional: true},
			},
			FriendlyName: "provider",
			ScopeId:      refscope.ProviderScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Description: lang.PlainText("A provider block is used to specify a provider configuration"),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Alias for using the same provider with different configurations for different resources, e.g. `eu-west`"),
				},
				"version": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Specifies a version constraint for the provider, e.g. `~> 1.0`"),
				},
			},
		},
	}
}

func runBlockSchema(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"command": {
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword: "apply",
						},
						schema.Keyword{
							Keyword: "plan",
						},
					},
					IsOptional: true,
				},
				// providers
				// expect_failures
			},
			Blocks: map[string]*schema.BlockSchema{
				"variables": variablesBlockSchema(),
				"assert": {
					MinItems: 1,
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Constraint:  schema.AnyExpression{OfType: cty.Bool},
								IsRequired:  true,
								Description: lang.Markdown("Condition to meet for the check to pass (any expression which evaluates to boolean)"),
							},
							"error_message": {
								Constraint:  schema.AnyExpression{OfType: cty.String},
								IsRequired:  true,
								Description: lang.Markdown("Text that Terraform will include as part of error messages when it detects an unmet condition"),
							},
						},
					},
				},
				// modules
				// plan_options
			},
		},
	}
}

func TestSchema(v *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"run":       runBlockSchema(v),
			"variables": variablesBlockSchema(),
			"provider":  providerBlockSchema(v),
		},
	}
}

func TestSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()
	if ver.GreaterThanOrEqual(v1_6) {
		return TestSchema(v), nil
	}

	return nil, NoCompatibleSchemaErr{Version: ver}
}
