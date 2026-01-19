// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	// removes module_variable_optional_attrs experiment (defined in 0.14)
	bs.Body.Attributes["experiments"] = &schema.AttributeSchema{
		Constraint: schema.Set{
			Elem: schema.OneOf{
				schema.Keyword{
					Keyword: "provider_sensitive_attrs",
					Name:    "feature",
				},
			},
		},
		IsOptional:  true,
		Description: lang.Markdown("A set of experimental language features to enable"),
	}

	return bs
}
