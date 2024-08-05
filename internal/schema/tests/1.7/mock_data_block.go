// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func mockDataBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Data Source Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Description: lang.PlainText(""), // TODO!
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"defaults": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsOptional:  true,
					Description: lang.Markdown(""), // TODO!
				},
			},
		},
	}
}
