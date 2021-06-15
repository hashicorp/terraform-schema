package schema

import (
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var expectedMergedSchema_v015_aliased = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"data": {
			Labels: []*schema.LabelSchema{
				{Name: "type"},
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.Number},
							schema.LiteralTypeExpr{Type: cty.Number},
						},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"hashicup_test"}],"attrs":[{"name":"provider","expr":{"addr":"hcc"}}]}`: {
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"backend": {
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"config1": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.ObjectExpr{
									Attributes: schema.ObjectExprAttributes{
										"first": {
											IsOptional: true,
											Expr: schema.ExprConstraints{
												schema.TraversalExpr{OfType: cty.String},
												schema.LiteralTypeExpr{Type: cty.String},
											},
										},
										"second": {
											IsOptional: true,
											Expr: schema.ExprConstraints{
												schema.TraversalExpr{OfType: cty.Number},
												schema.LiteralTypeExpr{Type: cty.Number},
											},
										},
										"third": {
											IsOptional: true,
											Expr: schema.ExprConstraints{
												schema.ObjectExpr{
													Attributes: schema.ObjectExprAttributes{
														"nested": {
															IsOptional: true,
															Expr: schema.ExprConstraints{
																schema.TraversalExpr{OfType: cty.String},
																schema.LiteralTypeExpr{Type: cty.String},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"config2": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.ListExpr{
									Elem: schema.ExprConstraints{
										schema.ObjectExpr{
											Attributes: schema.ObjectExprAttributes{
												"first": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.String},
														schema.LiteralTypeExpr{Type: cty.String},
													},
												},
												"second": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.Number},
														schema.LiteralTypeExpr{Type: cty.Number},
													},
												},
												"third": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.ObjectExpr{
															Attributes: schema.ObjectExprAttributes{
																"nested": {
																	IsOptional: true,
																	Expr: schema.ExprConstraints{
																		schema.TraversalExpr{OfType: cty.String},
																		schema.LiteralTypeExpr{Type: cty.String},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									MinItems: 2,
									MaxItems: 3,
								},
							},
						},
						"config3": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.SetExpr{
									Elem: schema.ExprConstraints{
										schema.ObjectExpr{
											Attributes: schema.ObjectExprAttributes{
												"first": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.String},
														schema.LiteralTypeExpr{Type: cty.String},
													},
												},
												"second": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.Number},
														schema.LiteralTypeExpr{Type: cty.Number},
													},
												},
												"third": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.ObjectExpr{
															Attributes: schema.ObjectExprAttributes{
																"nested": {
																	IsOptional: true,
																	Expr: schema.ExprConstraints{
																		schema.TraversalExpr{OfType: cty.String},
																		schema.LiteralTypeExpr{Type: cty.String},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									MinItems: 1,
									MaxItems: 5,
								},
							},
						},
						"config4": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.ObjectExpr{
											Attributes: schema.ObjectExprAttributes{
												"first": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.String},
														schema.LiteralTypeExpr{Type: cty.String},
													},
												},
												"second": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.TraversalExpr{OfType: cty.Number},
														schema.LiteralTypeExpr{Type: cty.Number},
													},
												},
												"third": {
													IsOptional: true,
													Expr: schema.ExprConstraints{
														schema.ObjectExpr{
															Attributes: schema.ObjectExprAttributes{
																"nested": {
																	IsOptional: true,
																	Expr: schema.ExprConstraints{
																		schema.TraversalExpr{OfType: cty.String},
																		schema.LiteralTypeExpr{Type: cty.String},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									MinItems: 9,
									MaxItems: 10,
								},
							},
						},
						"workspace": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
					Detail: "hashicorp/hashicup",
				},
			},
		},
		"provider": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {
						Expr: schema.ExprConstraints{
							schema.LiteralTypeExpr{Type: cty.String},
						},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"hcc"}]}`: {
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
					Detail:     "hashicorp/hashicup",
					HoverURL:   "https://registry.terraform.io/providers/hashicorp/hashicup/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/hashicorp/hashicup/latest/docs",
						Tooltip: "hashicorp/hashicup Documentation",
					},
				},
			},
		},
		"resource": {
			Labels: []*schema.LabelSchema{
				{Name: "type"},
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.Number},
							schema.LiteralTypeExpr{Type: cty.Number},
						},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
		"module": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"source": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsRequired: true,
						IsDepKey:   true,
					},
					"version": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
	},
}
