package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema, v *version.Version) *schema.BlockSchema {
	bs.Body.Blocks["required_providers"].Body = &schema.BodySchema{
		AnyAttribute: &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.ObjectExpr{
					Attributes: schema.ObjectExprAttributes{
						"source": &schema.AttributeSchema{
							Expr: schema.LiteralTypeOnly(cty.String),
							Description: lang.Markdown("The global source address for the provider " +
								"you intend to use, such as `hashicorp/aws`"),
						},
						"version": &schema.AttributeSchema{
							Expr: schema.LiteralTypeOnly(cty.String),
							Description: lang.Markdown("Version constraint specifying which subset of " +
								"available provider versions the module is compatible with, e.g. `~> 1.0`"),
						},
						"configuration_aliases": &schema.AttributeSchema{
							Expr: schema.ExprConstraints{
								schema.SetExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{
											Address: &schema.TraversalAddrSchema{
												ScopeId: refscope.ProviderScope,
											},
											Name: "provider",
										},
									},
								},
							},
							Description: lang.Markdown("Aliases under which to make the provider available, " +
								"such as `[ aws.eu-west, aws.us-east ]`"),
						},
					},
				},
				schema.LiteralTypeExpr{Type: cty.String},
			},
			Address: &schema.AttributeAddrSchema{
				Steps: []schema.AddrStep{
					schema.AttrNameStep{},
				},
				FriendlyName: "provider",
				AsReference:  true,
				ScopeId:      refscope.ProviderScope,
			},
			Description: lang.Markdown("Provider source, version constraint and its aliases"),
		},
	}
	bs.Body.Attributes["language"] = &schema.AttributeSchema{
		IsOptional: true,
		Expr: schema.ExprConstraints{
			schema.KeywordExpr{
				Keyword: "TF2021",
			},
		},
	}
	return bs
}
