// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func enforceBlockNestedSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Specifies conditions that must be true for the policy to pass"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"condition": {
					Constraint:  schema.AnyExpression{OfType: cty.Bool},
					IsRequired:  true,
					Description: lang.Markdown("An expression that must evaluate to `true` for the policy to pass"),
				},
				"error_message": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("A custom string describing why the policy failed"),
				},
			},
		},
	}
}
