// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchVariableBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["nullable"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.False},
		IsOptional:   true,
		Description:  lang.Markdown("Specifies whether `null` is a valid value for this variable"),
	}

	return bs
}
