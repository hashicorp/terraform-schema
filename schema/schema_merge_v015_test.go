package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

var depKeysModule = schema.DependencyKeys{
	Attributes: []schema.AttributeDependent{
		{
			Name: "source",
			Expr: schema.ExpressionValue{
				Static: cty.StringVal("source"),
			},
		},
	},
}

var data = schema.BlockSchema{
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
		`{"labels":[{"index":0,"value":"hashicup_test"}],"attrs":[{"name":"provider","expr":{"addr":"hashicup"}}]}`: {
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
		`{"labels":[{"index":0,"value":"hashicup_test"}]}`: {
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
}

var provider = schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{Name: "name"},
	},
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"alias": {Expr: schema.LiteralTypeOnly(cty.String), IsOptional: true},
		},
	},
	DependentBody: map[schema.SchemaKey]*schema.BodySchema{
		`{"labels":[{"index":0,"value":"hashicup"}]}`: {
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
}

var resource = schema.BlockSchema{
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
}

var moduleWithoutDependency = schema.BlockSchema{
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
}

var moduleWithDependency = schema.BlockSchema{
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
	DependentBody: map[schema.SchemaKey]*schema.BodySchema{
		schema.NewSchemaKey(depKeysModule): {
			TargetableAs: []*schema.Targetable{
				{
					Address: lang.Address{
						lang.RootStep{Name: "module"},
						lang.AttrStep{Name: "example"},
					},
					ScopeId:           refscope.ModuleScope,
					AsType:            cty.Object(map[string]cty.Type{}),
					NestedTargetables: []*schema.Targetable{},
				},
			},
			Attributes: map[string]*schema.AttributeSchema{
				"test": {
					Description: lang.PlainText("test var"),
					Expr: schema.ExprConstraints{
						schema.TraversalExpr{OfType: cty.String},
						schema.LiteralTypeExpr{Type: cty.String},
					},
					IsRequired: true,
					OriginForTarget: &schema.PathTarget{
						Address: schema.Address{
							schema.StaticStep{Name: "var"},
							schema.AttrNameStep{},
						},
						Path: lang.Path{
							Path:       "path",
							LanguageID: "terraform",
						},
						Constraints: schema.Constraints{
							ScopeId: "variable",
							Type:    cty.String,
						},
					},
				},
			},
		},
	},
}

var expectedMergedSchema_v015 = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"data":     &data,
		"provider": &provider,
		"resource": &resource,
		"module":   &moduleWithoutDependency,
	},
}

var expectedMergedSchemaWithModule_v015 = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"data":     &data,
		"provider": &provider,
		"resource": &resource,
		"module":   &moduleWithDependency,
	},
}
