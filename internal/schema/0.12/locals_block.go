package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var localsBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.StaticStep{Name: "local"},
		},
		FriendlyName: "local",
		BodyAsData:   true,
	},
	Description: lang.Markdown("Local values assigning names to expressions, so you can use these multiple times without repetition\n" +
		"e.g. `service_name = \"forum\"`"),
	Body: &schema.BodySchema{
		AnyAttribute: &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.TraversalExpr{OfType: cty.DynamicPseudoType},
				schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
			},
		},
	},
}
