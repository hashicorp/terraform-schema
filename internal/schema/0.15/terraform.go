// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Blocks["required_providers"].Body = &schema.BodySchema{
		AnyAttribute: &schema.AttributeSchema{
			Constraint: schema.OneOf{
				schema.Object{
					Attributes: schema.ObjectAttributes{
						"source": &schema.AttributeSchema{
							Constraint: schema.LiteralType{Type: cty.String},
							IsRequired: true,
							Description: lang.Markdown("The global source address for the provider " +
								"you intend to use, such as `hashicorp/aws`"),
						},
						"version": &schema.AttributeSchema{
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
							Description: lang.Markdown("Version constraint specifying which subset of " +
								"available provider versions the module is compatible with, e.g. `~> 1.0`"),
						},
						"configuration_aliases": &schema.AttributeSchema{
							IsOptional: true,
							Constraint: schema.Set{
								Elem: schema.Reference{
									Address: &schema.ReferenceAddrSchema{
										ScopeId: refscope.ProviderScope,
									},
									Name: "provider",
								},
							},
							Description: lang.Markdown("Aliases under which to make the provider available, " +
								"such as `[ aws.eu-west, aws.us-east ]`"),
						},
					},
				},
				schema.LiteralType{Type: cty.String},
			},
			Address: &schema.AttributeAddrSchema{
				Steps: []schema.AddrStep{
					schema.AttrNameStep{},
				},
				FriendlyName: "provider",
				AsReference:  true,
				ScopeId:      refscope.ProviderScope,
			},
			Description: lang.Markdown("Provider source, version constraint and its aliases"),
		},
	}
	bs.Body.Attributes["language"] = &schema.AttributeSchema{
		IsOptional: true,
		Constraint: schema.Keyword{
			Keyword: "TF2021",
		},
	}
	return bs
}
