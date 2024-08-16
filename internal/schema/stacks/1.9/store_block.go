// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func storeBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("A store block allows to retrieve credentials at plan and apply time. These credentials can be used as inputs to deployment blocks."),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Store type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Store name"),
				IsDepKey:               true,
			},
		},
		DependentBody: map[schema.SchemaKey]*schema.BodySchema{
			schema.NewSchemaKey(schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: "tfvars"},
				},
			}): {
				Attributes: map[string]*schema.AttributeSchema{
					"path": {
						IsRequired:  true,
						Constraint:  schema.LiteralType{Type: cty.String},
						Description: lang.Markdown("The path to the tfvars file."),
					},
				},
			},
			schema.NewSchemaKey(schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: "varset"},
				},
			}): {
				Attributes: map[string]*schema.AttributeSchema{
					"id": {
						IsRequired:  true,
						Constraint:  schema.LiteralType{Type: cty.String},
						Description: lang.Markdown("The id of the varset. In the form of 'varset-nnnnnnnnnnnnnnnn'."),
					},
				},
				AnyAttribute: &schema.AttributeSchema{
					IsComputed: true,
					Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
				},
			},
		},
	}
}
