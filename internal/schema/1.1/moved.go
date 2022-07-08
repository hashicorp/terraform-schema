package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

var movedBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Refactoring declaration to specify what address to move where"),
	Body: &schema.BodySchema{
		HoverURL: "https://www.terraform.io/language/modules/develop/refactoring#moved-block-syntax",
		Attributes: map[string]*schema.AttributeSchema{
			"from": {
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfScopeId: refscope.ModuleScope},
					schema.TraversalExpr{OfScopeId: refscope.ResourceScope},
				},
				IsRequired:  true,
				Description: lang.Markdown("Source address to move away from"),
			},
			"to": {
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfScopeId: refscope.ModuleScope},
					schema.TraversalExpr{OfScopeId: refscope.ResourceScope},
				},
				IsRequired:  true,
				Description: lang.Markdown("Destination address to move to"),
			},
		},
	},
}
