// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchVariableBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["deprecated"] = &schema.AttributeSchema{
		Constraint:  schema.LiteralType{Type: cty.String},
		IsOptional:  true,
		Description: lang.Markdown("Setting this value marks the variable as deprecated. The string value provided should describe the reason for deprecation and suggest an alternative. Any usage of a deprecated variable will result in a warning being emitted to the user."),
	}

	return bs
}
