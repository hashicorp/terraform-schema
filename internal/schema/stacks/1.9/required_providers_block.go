// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func requiredProvidersBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("What provider version to use within this configuration and where to source it from"),
		Body: &schema.BodySchema{
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
				Description: lang.Markdown("Provider source, version constraint"),
			},
		},
	}
}
