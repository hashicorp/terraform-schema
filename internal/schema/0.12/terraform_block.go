// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func terraformBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Terraform},
		Description:            lang.Markdown("Terraform block used to configure some high-level behaviors of Terraform"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"required_version": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
						"with this configuration, e.g. `~> 0.12`"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"backend": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Backend},
					Description: lang.Markdown("Backend configuration which defines exactly where and how " +
						"operations are performed, where state snapshots are stored, etc."),
					Labels: []*schema.LabelSchema{
						{
							Name:                   "backend type",
							SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
							Description:            lang.Markdown("Backend type"),
							IsDepKey:               true,
							Completable:            true,
						},
					},
					MaxItems:      1,
					DependentBody: backends.ConfigsAsDependentBodies(v),
				},
				"required_providers": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.RequiredProviders},
					Description:            lang.Markdown("What provider version to use within this configuration"),
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Version constraint"),
							Address: &schema.AttributeAddrSchema{
								Steps: []schema.AddrStep{
									schema.AttrNameStep{},
								},
								AsReference:  true,
								FriendlyName: "provider",
								ScopeId:      refscope.ProviderScope,
							},
						},
					},
					MaxItems: 1,
				},
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_18) {
		experiments := schema.OneOf{}
		if v.GreaterThanOrEqual(v0_12_20) {
			experiments = append(experiments, schema.Keyword{
				Keyword: "variable_validation",
				Name:    "feature",
			})
		}
		bs.Body.Attributes["experiments"] = &schema.AttributeSchema{
			Constraint: schema.Set{
				Elem: experiments,
			},
			IsOptional:  true,
			Description: lang.Markdown("A set of experimental language features to enable"),
		}
	}

	if v.GreaterThanOrEqual(v0_12_20) {
		bs.Body.Blocks["required_providers"].Body = &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Constraint: schema.OneOf{
					schema.Object{
						Attributes: schema.ObjectAttributes{
							"version": &schema.AttributeSchema{
								Constraint: schema.LiteralType{Type: cty.String},
								Description: lang.Markdown("Version constraint specifying which subset of " +
									"available provider versions the module is compatible with, e.g. `~> 1.0`"),
							},
						},
					},
					schema.LiteralType{Type: cty.String},
				},
				Description: lang.Markdown("Version constraint"),
			},
		}
	}

	return bs
}
