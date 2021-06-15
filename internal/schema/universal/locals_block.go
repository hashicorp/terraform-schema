package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

var localsBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Local values assigning names to expressions, so you can use these multiple times without repetition\n" +
		"e.g. `service_name = \"forum\"`"),
	Body: &schema.BodySchema{
		AnyAttribute: &schema.AttributeSchema{
			Address: &schema.AttributeAddrSchema{
				Steps: []schema.AddrStep{
					schema.StaticStep{Name: "local"},
					schema.AttrNameStep{},
				},
				ScopeId:    refscope.LocalScope,
				AsExprType: true,
			},
			Expr: schema.ExprConstraints{
				schema.TraversalExpr{OfType: cty.DynamicPseudoType},
				schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
			},
		},
	},
}
