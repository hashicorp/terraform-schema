// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func inputsNestedBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		MaxItems:    1,
		Description: lang.Markdown("Overrides input values scoped to this specific test case. Each attribute key must match the name of an `input` block declared in the linked policy file."),
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Constraint:  schema.AnyExpression{OfType: cty.DynamicPseudoType},
				Description: lang.Markdown("Value for the input with the matching name declared in the linked policy"),
			},
		},
	}
}
