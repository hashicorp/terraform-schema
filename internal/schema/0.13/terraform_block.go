package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var terraformBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Terraform block used to configure some high-level behaviors of Terraform"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"required_version": {
				Expr:       schema.LiteralTypeOnly(cty.String),
				IsOptional: true,
				Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
					"with this configuration, e.g. `~> 0.12`"),
			},
			"experiments": {
				Expr: schema.ExprConstraints{
					schema.TupleConsExpr{
						Name: "set of features",
					},
				},
				IsOptional:  true,
				Description: lang.Markdown("A set of experimental language features to enable"),
			},
		},
		Blocks: map[string]*schema.BlockSchema{
			"backend": {
				Description: lang.Markdown("Backend configuration which defines exactly where and how " +
					"operations are performed, where state snapshots are stored, etc."),
				Labels: []*schema.LabelSchema{
					{
						Name:        "type",
						Description: lang.Markdown("Backend Type"),
						IsDepKey:    true,
					},
				},
			},
			"provider_meta": {
				Description: lang.Markdown("Metadata to pass into a provider which supports this"),
				Labels: []*schema.LabelSchema{
					{
						Name:        "name",
						Description: lang.Markdown("Provider Name"),
						IsDepKey:    true,
					},
				},
			},
			"required_providers": {
				Description: lang.Markdown("What provider version to use within this configuration " +
					"and where to source it from"),
				Body: &schema.BodySchema{
					AnyAttribute: &schema.AttributeSchema{
						Expr: schema.ExprConstraints{
							schema.ObjectExpr{
								Attributes: schema.ObjectExprAttributes{
									"source": schema.ObjectAttribute{
										Expr: schema.LiteralTypeOnly(cty.String),
										Description: lang.Markdown("The global source address for the provider " +
											"you intend to use, such as `hashicorp/aws`"),
									},
									"version": schema.ObjectAttribute{
										Expr: schema.LiteralTypeOnly(cty.String),
										Description: lang.Markdown("Version constraint specifying which subset of " +
											"available provider versions the module is compatible with, e.g. `~> 1.0`"),
									},
								},
							},
							schema.LiteralTypeExpr{Type: cty.String},
						},
						Description: lang.Markdown("Provider source and version constraint"),
					},
				},
			},
		},
	},
}
