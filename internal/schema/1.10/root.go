// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_9_mod "github.com/hashicorp/terraform-schema/internal/schema/1.9"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_9_mod.ModuleSchema(v)

	bs.Blocks["variable"].Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		IsOptional:  true,
		Constraint:  schema.LiteralType{Type: cty.Bool},
		Description: lang.PlainText("Whether the value is ephemeral and should not be persisted in the state"),
	}
	bs.Blocks["output"].Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		IsOptional:  true,
		Constraint:  schema.LiteralType{Type: cty.Bool},
		Description: lang.PlainText("Whether the value is ephemeral and should not be persisted in the state"),
	}

	bs.Blocks["ephemeral"] = ephemeralBlockSchema()

	// all the depends_on attributes can refer to ephemeral blocks
	constraint := schema.Set{
		Elem: schema.OneOf{
			schema.Reference{OfScopeId: refscope.DataScope},
			schema.Reference{OfScopeId: refscope.ModuleScope},
			schema.Reference{OfScopeId: refscope.ResourceScope},
			schema.Reference{OfScopeId: refscope.EphemeralScope}, // This one is new, but overriding is easier than adding to each list
			schema.Reference{OfScopeId: refscope.VariableScope},
			schema.Reference{OfScopeId: refscope.LocalScope},
		},
	}
	bs.Blocks["resource"].Body.Attributes["depends_on"].Constraint = constraint
	bs.Blocks["data"].Body.Attributes["depends_on"].Constraint = constraint
	bs.Blocks["output"].Body.Attributes["depends_on"].Constraint = constraint
	bs.Blocks["module"].Body.Attributes["depends_on"].Constraint = constraint
	bs.Blocks["check"].Body.Blocks["data"].Body.Attributes["depends_on"].Constraint = constraint

	return bs
}
