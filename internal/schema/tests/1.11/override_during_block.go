// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func overrideDuringBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Allows overriding the point in time where terraform generates data"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"override_during": {
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "apply",
							Description: lang.Markdown("Default behavior where data is generated during the apply operation and (known after apply) is returned during the plan"),
						},
						schema.Keyword{
							Keyword:     "plan",
							Description: lang.Markdown("Allows to generate data during the plan operation. The same data will be used during the apply"),
						},
					},
				},
			},
		},
	}
}
