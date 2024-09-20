// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

func overrideDataBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Allows overriding the values of a specific data source in the targeted configuration"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"target": {
					Constraint: schema.Reference{
						OfScopeId: refscope.DataScope,
					},
					IsRequired:  true,
					Description: lang.Markdown("Reference to the data source to override"),
				},
				"values": {
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
