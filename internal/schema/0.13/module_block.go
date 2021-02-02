package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var moduleBlockSchema = &schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Reference Name"),
		},
	},
	Description: lang.PlainText("Module block to call a locally or remotely stored module"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"source": {
				Expr: schema.LiteralTypeOnly(cty.String),
				Description: lang.Markdown("Source where to load the module from, " +
					"a local directory (e.g. `./module`) or a remote address - e.g. " +
					"`hashicorp/consul/aws` (Terraform Registry address) or " +
					"`github.com/hashicorp/example` (GitHub)"),
				IsRequired: true,
				IsDepKey:   true,
			},
			"version": {
				Expr:       schema.LiteralTypeOnly(cty.String),
				IsOptional: true,
				Description: lang.Markdown("Constraint to set the version of the module, e.g. `~> 1.0`." +
					" Only applicable to modules in a module registry."),
			},
			"providers": {
				Expr: schema.ExprConstraints{
					schema.MapExpr{
						Name: "map of provider references",
					},
				},
				IsOptional:  true,
				Description: lang.Markdown("Explicit mapping of providers which the module uses"),
			},
			"count": {
				Expr:        schema.LiteralTypeOnly(cty.Number),
				IsOptional:  true,
				Description: lang.Markdown("Number of instances of this module, e.g. `3`"),
			},
			"for_each": {
				Expr: schema.ExprConstraints{
					schema.LiteralTypeExpr{Type: cty.Set(cty.DynamicPseudoType)},
					schema.LiteralTypeExpr{Type: cty.Map(cty.DynamicPseudoType)},
				},
				IsOptional:  true,
				Description: lang.Markdown("A set or a map where each item represents an instance of this module"),
			},
			"depends_on": {
				Expr: schema.ExprConstraints{
					schema.TupleConsExpr{
						Name: "set of references",
					},
				},
				IsOptional:  true,
				Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
			},
		},
	},
}
