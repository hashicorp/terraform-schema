// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func overrideDuringAttributeSchema() *schema.AttributeSchema {
	return &schema.AttributeSchema{
		Description: lang.Markdown("Allows overriding the point in time where terraform generates data"),
		IsOptional:  true,
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
	}
}
