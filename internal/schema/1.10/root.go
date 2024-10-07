// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_9_mod "github.com/hashicorp/terraform-schema/internal/schema/1.9"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_9_mod.ModuleSchema(v)

	bs.Blocks["variable"].Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		IsOptional: true,
		Constraint: schema.LiteralType{Type: cty.Bool},
	}
	bs.Blocks["output"].Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		IsOptional: true,
		Constraint: schema.LiteralType{Type: cty.Bool},
	}

	return bs
}
