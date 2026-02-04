// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchOutputBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["deprecated"] = &schema.AttributeSchema{
		Constraint:  schema.LiteralType{Type: cty.String},
		IsOptional:  true,
		Description: lang.Markdown("Setting this value marks the output as deprecated. The string value provided should describe the reason for deprecation and suggest an alternative. Any usage of a deprecated output will result in a warning being emitted to the user."),
	}

	return bs
}
