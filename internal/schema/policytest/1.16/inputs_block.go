// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func inputsBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		MaxItems:               1,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Inputs},
		Description:            lang.Markdown("Supplies values for inputs declared in the linked policy file. Each attribute key must match the name of an `input` block defined in the target `*.policy.hcl` file."),
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Constraint:  schema.AnyExpression{OfType: cty.DynamicPseudoType},
				Description: lang.Markdown("Value for the input with the matching name declared in the linked policy"),
			},
		},
	}
}
