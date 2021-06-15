package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

var outputBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.StaticStep{Name: "output"},
			schema.LabelStep{Index: 0},
		},
		FriendlyName: "output",
		ScopeId:      refscope.OutputScope,
		AsReference:  true,
	},
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Output Name"),
		},
	},
	Description: lang.PlainText("Output value for consumption by another module or a human interacting via the UI"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"description": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.PlainText("Human-readable description of the output (for documentation and UI)"),
			},
			"value": {
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfType: cty.DynamicPseudoType},
					schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
				},
				IsRequired:  true,
				Description: lang.PlainText("Value, typically a reference to an attribute of a resource or a data source"),
			},
			"depends_on": {
				Expr: schema.ExprConstraints{
					schema.TupleConsExpr{
						Name: "set of references",
						AnyElem: schema.ExprConstraints{
							schema.TraversalExpr{OfScopeId: refscope.DataScope},
							schema.TraversalExpr{OfScopeId: refscope.ModuleScope},
							schema.TraversalExpr{OfScopeId: refscope.ResourceScope},
						},
					},
				},
				IsOptional:  true,
				Description: lang.PlainText("Set of references to hidden dependencies (e.g. resources or data sources)"),
			},
		},
	},
}
