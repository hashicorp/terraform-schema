// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func dataBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "data"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "data",
			ScopeId:              refscope.DataScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Data Source Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				Description:            lang.PlainText("Reference Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.PlainText("Introduces a mock data source containing the specific values returned by a lookup. It uses two labels to refer to the data source type and name."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"attrs": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsRequired:  true,
					Description: lang.Markdown("Specify the values that should be returned for specific attributes"),
				},
			},
		},
	}
}
