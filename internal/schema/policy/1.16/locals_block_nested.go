// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func localsBlockNestedSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Locals},
		Description:            lang.Markdown("Local values to be used in the scope"),
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Address: &schema.AttributeAddrSchema{
					Steps: []schema.AddrStep{
						schema.StaticStep{Name: "local"},
						schema.AttrNameStep{},
					},
					ScopeId:     refscope.LocalScope,
					AsExprType:  true,
					AsReference: true,
				},
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
			},
			TargetableFromCurrentBlock: true,
		},
		MaxItems: 1,
	}
}
