// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func variablesBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		MaxItems:               1,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Variables},
		Description:            lang.Markdown("Provides values for input variables within your configuration directly from your test files"),
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Address: &schema.AttributeAddrSchema{
					Steps: []schema.AddrStep{
						schema.StaticStep{Name: "var"},
						schema.AttrNameStep{},
					},
					ScopeId:     refscope.VariableScope,
					AsExprType:  true,
					AsReference: true,
				},
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
			},
		},
	}
}
