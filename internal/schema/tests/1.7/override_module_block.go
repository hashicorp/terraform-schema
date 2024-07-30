// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

func overrideModuleBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText(""), // TODO!
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"target": {
					Constraint: schema.Reference{
						OfScopeId: refscope.ModuleScope,
					},
					IsRequired:  true,
					Description: lang.Markdown(""), // TODO!
				},
				"outputs": {
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
