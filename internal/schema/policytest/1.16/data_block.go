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

func dataBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "data"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName: "data",
			ScopeId:      refscope.DataScope,
			AsReference:  true,
			InferBody:    true,
			BodyAsData:   true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "data_type",
				Description:            lang.PlainText("Data Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
			},
			{
				Name:                   "name",
				Description:            lang.PlainText("Reference Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.PlainText("Unlike `resource`, this block contains the returned attributes of a data source. These can be referenced by other resource blocks using `data.<name>.<attribute>`"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"attrs": {
					Address: &schema.AttributeAddrSchema{
						Skip:    true,
						ScopeId: refscope.DataScope,
					},
					Description: lang.Markdown("Specify the values that should be returned for specific attributes"),
					Constraint:  schema.AnyExpression{OfType: cty.DynamicPseudoType},
					IsRequired:  true,
				},
			},
		},
	}
}
