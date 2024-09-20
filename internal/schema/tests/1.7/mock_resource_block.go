// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func mockResourceBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Resource Type"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Description: lang.PlainText("Allows to specify specific values for targeted resources"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"defaults": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsOptional:  true,
					Description: lang.Markdown("Specify the values that should be returned for specific attributes"),
				},
			},
		},
	}
}
