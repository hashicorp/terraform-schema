package schema

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestProviderSchemaFromJson_empty(t *testing.T) {
	jsonSchema := &tfjson.ProviderSchema{}
	providerAddr := tfaddr.NewDefaultProvider("aws")

	ps := ProviderSchemaFromJson(jsonSchema, providerAddr)
	expectedPs := &ProviderSchema{
		Resources:   map[string]*schema.BodySchema{},
		DataSources: map[string]*schema.BodySchema{},
	}

	if diff := cmp.Diff(expectedPs, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("provider schema mismatch: %s", diff)
	}
}

func TestProviderSchemaFromJson_basic(t *testing.T) {
	rawSchema := `{
	"resource_schemas": {
		"aws_security_group": {
			"version": 1,
			"block": {
				"attributes": {
					"textfield": {
						"type": "string",
						"description_kind": "plain",
						"optional": true
					},
					"simple_list": {
						"type": [
							"list",
							"string"
						],
						"description_kind": "plain",
						"optional": true
					},
					"ingress": {
						"type": [
							"set",
							[
								"object",
								{
									"cidr_blocks": [
										"list",
										"string"
									],
									"description": "string",
									"from_port": "number",
									"self": "bool"
								}
							]
						],
						"description_kind": "plain",
						"optional": true,
						"computed": true
					},
					"egress": {
						"type": [
							"list",
							[
								"object",
								{
									"cidr_blocks": [
										"list",
										"string"
									],
									"description": "string",
									"from_port": "number",
									"self": "bool"
								}
							]
						],
						"description_kind": "plain",
						"optional": true,
						"computed": true
					}
				}
			}
		}
	}
}`
	jsonSchema := &tfjson.ProviderSchema{}
	err := json.Unmarshal([]byte(rawSchema), jsonSchema)
	if err != nil {
		t.Fatal(err)
	}
	providerAddr := tfaddr.NewDefaultProvider("aws")

	ps := ProviderSchemaFromJson(jsonSchema, providerAddr)
	expectedPs := &ProviderSchema{
		Resources: map[string]*schema.BodySchema{
			"aws_security_group": {
				Attributes: map[string]*schema.AttributeSchema{
					"egress": {
						IsOptional: true,
						IsComputed: true,
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.List(cty.Object(map[string]cty.Type{
								"cidr_blocks": cty.List(cty.String),
								"description": cty.String,
								"from_port":   cty.Number,
								"self":        cty.Bool,
							}))},
							schema.LiteralTypeExpr{Type: cty.List(cty.Object(map[string]cty.Type{
								"cidr_blocks": cty.List(cty.String),
								"description": cty.String,
								"from_port":   cty.Number,
								"self":        cty.Bool,
							}))},
							schema.ListExpr{
								Elem: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.Object(map[string]cty.Type{
										"cidr_blocks": cty.List(cty.String),
										"description": cty.String,
										"from_port":   cty.Number,
										"self":        cty.Bool,
									})},
									schema.LiteralTypeExpr{Type: cty.Object(map[string]cty.Type{
										"cidr_blocks": cty.List(cty.String),
										"description": cty.String,
										"from_port":   cty.Number,
										"self":        cty.Bool,
									})},
									schema.ObjectExpr{
										Attributes: schema.ObjectExprAttributes{
											"cidr_blocks": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.List(cty.String)},
													schema.LiteralTypeExpr{Type: cty.List(cty.String)},
													schema.ListExpr{
														Elem: schema.ExprConstraints{
															schema.TraversalExpr{OfType: cty.String},
															schema.LiteralTypeExpr{Type: cty.String},
														},
													},
												},
											},
											"description": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.String},
													schema.LiteralTypeExpr{Type: cty.String},
												},
											},
											"from_port": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.Number},
													schema.LiteralTypeExpr{Type: cty.Number},
												},
											},
											"self": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.Bool},
													schema.LiteralTypeExpr{Type: cty.Bool},
												},
											},
										},
									},
								},
							},
						},
					},
					"ingress": {
						IsOptional: true,
						IsComputed: true,
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.Set(cty.Object(map[string]cty.Type{
								"cidr_blocks": cty.List(cty.String),
								"description": cty.String,
								"from_port":   cty.Number,
								"self":        cty.Bool,
							}))},
							schema.LiteralTypeExpr{Type: cty.Set(cty.Object(map[string]cty.Type{
								"cidr_blocks": cty.List(cty.String),
								"description": cty.String,
								"from_port":   cty.Number,
								"self":        cty.Bool,
							}))},
							schema.SetExpr{
								Elem: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.Object(map[string]cty.Type{
										"cidr_blocks": cty.List(cty.String),
										"description": cty.String,
										"from_port":   cty.Number,
										"self":        cty.Bool,
									})},
									schema.LiteralTypeExpr{Type: cty.Object(map[string]cty.Type{
										"cidr_blocks": cty.List(cty.String),
										"description": cty.String,
										"from_port":   cty.Number,
										"self":        cty.Bool,
									})},
									schema.ObjectExpr{
										Attributes: schema.ObjectExprAttributes{
											"cidr_blocks": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.List(cty.String)},
													schema.LiteralTypeExpr{Type: cty.List(cty.String)},
													schema.ListExpr{
														Elem: schema.ExprConstraints{
															schema.TraversalExpr{OfType: cty.String},
															schema.LiteralTypeExpr{Type: cty.String},
														},
													},
												},
											},
											"description": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.String},
													schema.LiteralTypeExpr{Type: cty.String},
												},
											},
											"from_port": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.Number},
													schema.LiteralTypeExpr{Type: cty.Number},
												},
											},
											"self": {
												IsRequired: true,
												Expr: schema.ExprConstraints{
													schema.TraversalExpr{OfType: cty.Bool},
													schema.LiteralTypeExpr{Type: cty.Bool},
												},
											},
										},
									},
								},
							},
						},
					},
					"simple_list": {
						IsOptional: true,
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.List(cty.String)},
							schema.LiteralTypeExpr{Type: cty.List(cty.String)},
							schema.ListExpr{
								Elem: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.String},
									schema.LiteralTypeExpr{Type: cty.String},
								},
							},
						},
					},
					"textfield": {
						IsOptional: true,
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.String},
							schema.LiteralTypeExpr{Type: cty.String},
						},
					},
				},
				Blocks: map[string]*schema.BlockSchema{
					"egress": {
						Type: schema.BlockTypeList,
						Body: &schema.BodySchema{
							Attributes: map[string]*schema.AttributeSchema{
								"cidr_blocks": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.List(cty.String)},
										schema.LiteralTypeExpr{Type: cty.List(cty.String)},
										schema.ListExpr{
											Elem: schema.ExprConstraints{
												schema.TraversalExpr{OfType: cty.String},
												schema.LiteralTypeExpr{Type: cty.String},
											},
										},
									},
								},
								"description": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
								"from_port": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.Number},
										schema.LiteralTypeExpr{Type: cty.Number},
									},
								},
								"self": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.Bool},
										schema.LiteralTypeExpr{Type: cty.Bool},
									},
								},
							},
						},
					},
					"ingress": {
						Type: schema.BlockTypeSet,
						Body: &schema.BodySchema{
							Attributes: map[string]*schema.AttributeSchema{
								"cidr_blocks": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.List(cty.String)},
										schema.LiteralTypeExpr{Type: cty.List(cty.String)},
										schema.ListExpr{
											Elem: schema.ExprConstraints{
												schema.TraversalExpr{OfType: cty.String},
												schema.LiteralTypeExpr{Type: cty.String},
											},
										},
									},
								},
								"description": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
								"from_port": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.Number},
										schema.LiteralTypeExpr{Type: cty.Number},
									},
								},
								"self": {
									IsOptional: true,
									Expr: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.Bool},
										schema.LiteralTypeExpr{Type: cty.Bool},
									},
								},
							},
						},
					},
				},
				Detail: "hashicorp/aws",
			},
		},
		DataSources: map[string]*schema.BodySchema{},
	}

	if diff := cmp.Diff(expectedPs, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("provider schema mismatch: %s", diff)
	}
}
