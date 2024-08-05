// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func mockProviderBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
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
		Description: lang.PlainText(""), // TODO!
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown(""), // TODO!
				},
				"source": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown(""), // TODO!
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"mock_resource":     mockResourceBlockSchema(),
				"mock_data":         mockDataBlockSchema(),
				"override_resource": overrideResourceBlockSchema(),
				"override_data":     overrideDataBlockSchema(),
			},
		},
	}
}
