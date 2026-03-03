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

func upstreamInputBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("An upstream_input block specifies another Stack in the same project to consume outputs from. Declare an upstream_input block for each Stack you want to reference."),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "upstream_input"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName:             "upstream_input",
			ScopeId:                  refscope.UpstreamInputScope,
			AsReference:              true,
			SupportUnknownNestedRefs: true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Upstream Stack Name"),
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"type": {
					Constraint: schema.OneOf{
						schema.LiteralValue{Value: cty.StringVal("stack")},
					},
					IsRequired:  true,
					Description: lang.PlainText("The type of upstream input"),
				},
				"source": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsRequired:  true,
					Description: lang.PlainText("The upstream Stack's URL, in the format: 'app.terraform.io/{organization_name}/{project_name}/{upstream_stack_name}'"),
				},
			},
		},
	}
}
