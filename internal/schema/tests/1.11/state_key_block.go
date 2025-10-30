// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func stateKeyBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Allows explicitly setting the identifier for a state file"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"state_key": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsRequired:  false,
					Description: lang.Markdown("Identifier used for the terraform state file. Read more on [module states](https://developer.hashicorp.com/terraform/language/tests#modules-state)"),
				},
			},
		},
	}
}
