package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func terraformBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Description: lang.Markdown("Terraform block used to configure some high-level behaviors of Terraform"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"required_version": {
					Expr:       schema.LiteralTypeOnly(cty.String),
					IsOptional: true,
					Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
						"with this configuration, e.g. `~> 0.12`"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"backend": {
					Description: lang.Markdown("Backend configuration which defines exactly where and how " +
						"operations are performed, where state snapshots are stored, etc."),
					Labels: []*schema.LabelSchema{
						{
							Name:        "backend type",
							Description: lang.Markdown("Backend type"),
							IsDepKey:    true,
							Completable: true,
						},
					},
					MaxItems:      1,
					DependentBody: backends.ConfigsAsDependentBodies(v),
				},
				"required_providers": {
					Description: lang.Markdown("What provider version to use within this configuration"),
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							Expr:        schema.LiteralTypeOnly(cty.String),
							Description: lang.Markdown("Version constraint"),
							Address: &schema.AttributeAddrSchema{
								Steps: []schema.AddrStep{
									schema.AttrNameStep{},
								},
								AsReference:  true,
								FriendlyName: "provider",
								ScopeId:      refscope.ProviderScope,
							},
						},
					},
					MaxItems: 1,
				},
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_18) {
		experiments := schema.TupleConsExpr{
			AnyElem: schema.ExprConstraints{},
			Name:    "set of features",
		}
		if v.GreaterThanOrEqual(v0_12_20) {
			experiments.AnyElem = append(experiments.AnyElem, schema.KeywordExpr{
				Keyword: "variable_validation",
				Name:    "feature",
			})
		}
		bs.Body.Attributes["experiments"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				experiments,
			},
			IsOptional:  true,
			Description: lang.Markdown("A set of experimental language features to enable"),
		}
	}

	if v.GreaterThanOrEqual(v0_12_20) {
		bs.Body.Blocks["required_providers"].Body = &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Expr: schema.ExprConstraints{
					schema.ObjectExpr{
						Attributes: schema.ObjectExprAttributes{
							"version": &schema.AttributeSchema{
								Expr: schema.LiteralTypeOnly(cty.String),
								Description: lang.Markdown("Version constraint specifying which subset of " +
									"available provider versions the module is compatible with, e.g. `~> 1.0`"),
							},
						},
					},
					schema.LiteralTypeExpr{Type: cty.String},
				},
				Description: lang.Markdown("Version constraint"),
			},
		}
	}

	return bs
}
