// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-schema/internal/addr"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func TestProviderSchemaFromJson_empty(t *testing.T) {
	jsonSchema := &tfjson.ProviderSchema{}
	providerAddr := addr.NewDefaultProvider("aws")

	ps := ProviderSchemaFromJson(jsonSchema, providerAddr)
	expectedPs := &ProviderSchema{
		Resources:          map[string]*schema.BodySchema{},
		EphemeralResources: map[string]*schema.BodySchema{},
		DataSources:        map[string]*schema.BodySchema{},
		Functions:          map[string]*schema.FunctionSignature{},
		ListResources:      map[string]*schema.BodySchema{},
		ActionResources:    map[string]*schema.BodySchema{},
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
	providerAddr := addr.NewDefaultProvider("aws")

	ps := ProviderSchemaFromJson(jsonSchema, providerAddr)
	expectedPs := &ProviderSchema{
		Resources: map[string]*schema.BodySchema{
			"aws_security_group": {
				Attributes: map[string]*schema.AttributeSchema{
					"egress": {
						IsOptional: true,
						IsComputed: true,
						Constraint: schema.OneOf{
							schema.AnyExpression{
								OfType: cty.List(cty.Object(map[string]cty.Type{
									"cidr_blocks": cty.List(cty.String),
									"description": cty.String,
									"from_port":   cty.Number,
									"self":        cty.Bool,
								})),
								SkipLiteralComplexTypes: true,
							},
							schema.List{
								Elem: schema.OneOf{
									schema.AnyExpression{
										OfType: cty.Object(map[string]cty.Type{
											"cidr_blocks": cty.List(cty.String),
											"description": cty.String,
											"from_port":   cty.Number,
											"self":        cty.Bool,
										}),
										SkipLiteralComplexTypes: true,
									},
									schema.Object{
										Attributes: schema.ObjectAttributes{
											"cidr_blocks": {
												IsRequired: true,
												Constraint: schema.OneOf{
													schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
													schema.List{
														Elem: schema.AnyExpression{OfType: cty.String},
													},
												},
											},
											"description": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.String},
											},
											"from_port": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.Number},
											},
											"self": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.Bool},
											},
										},
										AllowInterpolatedKeys: true,
									},
								},
							},
						},
					},
					"ingress": {
						IsOptional: true,
						IsComputed: true,
						Constraint: schema.OneOf{
							schema.AnyExpression{
								OfType: cty.Set(cty.Object(map[string]cty.Type{
									"cidr_blocks": cty.List(cty.String),
									"description": cty.String,
									"from_port":   cty.Number,
									"self":        cty.Bool,
								})),
								SkipLiteralComplexTypes: true,
							},
							schema.Set{
								Elem: schema.OneOf{
									schema.AnyExpression{
										OfType: cty.Object(map[string]cty.Type{
											"cidr_blocks": cty.List(cty.String),
											"description": cty.String,
											"from_port":   cty.Number,
											"self":        cty.Bool,
										}),
										SkipLiteralComplexTypes: true,
									},
									schema.Object{
										Attributes: schema.ObjectAttributes{
											"cidr_blocks": {
												IsRequired: true,
												Constraint: schema.OneOf{
													schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
													schema.List{
														Elem: schema.AnyExpression{OfType: cty.String},
													},
												},
											},
											"description": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.String},
											},
											"from_port": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.Number},
											},
											"self": {
												IsRequired: true,
												Constraint: schema.AnyExpression{OfType: cty.Bool},
											},
										},
										AllowInterpolatedKeys: true,
									},
								},
							},
						},
					},
					"simple_list": {
						IsOptional: true,
						Constraint: schema.OneOf{
							schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
							schema.List{
								Elem: schema.AnyExpression{OfType: cty.String},
							},
						},
					},
					"textfield": {
						IsOptional: true,
						Constraint: schema.AnyExpression{OfType: cty.String},
					},
				},
				Blocks: map[string]*schema.BlockSchema{
					"egress": {
						Type: schema.BlockTypeList,
						Body: &schema.BodySchema{
							Attributes: map[string]*schema.AttributeSchema{
								"cidr_blocks": {
									IsOptional: true,
									Constraint: schema.OneOf{
										schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
										schema.List{
											Elem: schema.AnyExpression{OfType: cty.String},
										},
									},
								},
								"description": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.String},
								},
								"from_port": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.Number},
								},
								"self": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.Bool},
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
									Constraint: schema.OneOf{
										schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
										schema.List{
											Elem: schema.AnyExpression{OfType: cty.String},
										},
									},
								},
								"description": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.String},
								},
								"from_port": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.Number},
								},
								"self": {
									IsOptional: true,
									Constraint: schema.AnyExpression{OfType: cty.Bool},
								},
							},
						},
					},
				},
				Detail: "hashicorp/aws",
			},
		},
		EphemeralResources: map[string]*schema.BodySchema{},
		DataSources:        map[string]*schema.BodySchema{},
		Functions:          map[string]*schema.FunctionSignature{},
		ListResources:      map[string]*schema.BodySchema{},
		ActionResources:    map[string]*schema.BodySchema{},
	}

	if diff := cmp.Diff(expectedPs, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("provider schema mismatch: %s", diff)
	}
}

func TestProviderSchemaFromJson_nested_set_list(t *testing.T) {
	rawSchema := `{
	"resource_schemas": {
		"azurerm_site_recovery_replicated_vm": {
			"version": 0,
			"block": {
				"attributes": {
					"managed_disk": {
						"type": [
							"set",
							[
								"object",
								{
									"target_disk_encryption": [
										"list",
										[
											"object",
											{
												"disk_encryption_key": "string",
												"key_encryption_key": "string"
											}
										]
									]
								}
							]
						],
						"description_kind": "plain",
						"optional": true
					}
				},
				"block_types": {},
				"description_kind": "plain"
			}
		}
	}
}`
	jsonSchema := &tfjson.ProviderSchema{}
	err := json.Unmarshal([]byte(rawSchema), jsonSchema)
	if err != nil {
		t.Fatal(err)
	}
	providerAddr := addr.NewDefaultProvider("aws")

	ps := ProviderSchemaFromJson(jsonSchema, providerAddr)

	expectedPs := &ProviderSchema{
		Resources: map[string]*schema.BodySchema{
			"azurerm_site_recovery_replicated_vm": {
				Attributes: map[string]*schema.AttributeSchema{
					"managed_disk": {
						IsOptional: true,
						Constraint: schema.OneOf{
							schema.AnyExpression{
								OfType: cty.Set(cty.Object(
									map[string]cty.Type{
										"target_disk_encryption": cty.List(cty.Object(
											map[string]cty.Type{
												"disk_encryption_key": cty.String,
												"key_encryption_key":  cty.String,
											},
										)),
									},
								)),
								SkipLiteralComplexTypes: true,
							},
							schema.Set{
								Elem: schema.OneOf{
									schema.AnyExpression{
										OfType: cty.Object(
											map[string]cty.Type{
												"target_disk_encryption": cty.List(cty.Object(
													map[string]cty.Type{
														"disk_encryption_key": cty.String,
														"key_encryption_key":  cty.String,
													},
												)),
											},
										),
										SkipLiteralComplexTypes: true,
									},
									schema.Object{
										Attributes: schema.ObjectAttributes{
											"target_disk_encryption": {
												IsRequired: true,
												Constraint: schema.OneOf{
													schema.AnyExpression{
														OfType: cty.List(cty.Object(
															map[string]cty.Type{
																"disk_encryption_key": cty.String,
																"key_encryption_key":  cty.String,
															},
														)),
														SkipLiteralComplexTypes: true,
													},
													schema.List{
														Elem: schema.OneOf{
															schema.AnyExpression{
																OfType: cty.Object(
																	map[string]cty.Type{
																		"disk_encryption_key": cty.String,
																		"key_encryption_key":  cty.String,
																	},
																),
																SkipLiteralComplexTypes: true,
															},
															schema.Object{
																Attributes: schema.ObjectAttributes{
																	"disk_encryption_key": {
																		IsRequired: true,
																		Constraint: schema.AnyExpression{OfType: cty.String},
																	},
																	"key_encryption_key": {
																		IsRequired: true,
																		Constraint: schema.AnyExpression{OfType: cty.String},
																	},
																},
																AllowInterpolatedKeys: true,
															},
														},
													},
												},
											},
										},
										AllowInterpolatedKeys: true,
									},
								},
							},
						},
					},
				},
				Blocks: map[string]*schema.BlockSchema{
					"managed_disk": {
						Type: schema.BlockTypeSet,
						Body: &schema.BodySchema{
							Blocks: map[string]*schema.BlockSchema{
								"target_disk_encryption": {
									Type: schema.BlockTypeList,
									Body: &schema.BodySchema{
										Attributes: map[string]*schema.AttributeSchema{
											"disk_encryption_key": {
												IsOptional: true,
												Constraint: schema.AnyExpression{OfType: cty.String},
											},
											"key_encryption_key": {
												IsOptional: true,
												Constraint: schema.AnyExpression{OfType: cty.String},
											},
										},
									},
								},
							},
							Attributes: map[string]*schema.AttributeSchema{
								"target_disk_encryption": {
									IsOptional: true,
									Constraint: schema.OneOf{
										schema.AnyExpression{
											OfType: cty.List(cty.Object(
												map[string]cty.Type{
													"disk_encryption_key": cty.String,
													"key_encryption_key":  cty.String,
												},
											)),
											SkipLiteralComplexTypes: true,
										},
										schema.List{
											Elem: schema.OneOf{
												schema.AnyExpression{
													OfType: cty.Object(
														map[string]cty.Type{
															"disk_encryption_key": cty.String,
															"key_encryption_key":  cty.String,
														},
													),
													SkipLiteralComplexTypes: true,
												},
												schema.Object{
													Attributes: schema.ObjectAttributes{
														"disk_encryption_key": {
															IsRequired: true,
															Constraint: schema.AnyExpression{OfType: cty.String},
														},
														"key_encryption_key": {
															IsRequired: true,
															Constraint: schema.AnyExpression{OfType: cty.String},
														},
													},
													AllowInterpolatedKeys: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Detail: "hashicorp/aws",
			},
		},
		EphemeralResources: map[string]*schema.BodySchema{},
		DataSources:        map[string]*schema.BodySchema{},
		Functions:          map[string]*schema.FunctionSignature{},
		ListResources:      map[string]*schema.BodySchema{},
		ActionResources:    map[string]*schema.BodySchema{},
	}

	if diff := cmp.Diff(expectedPs, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("provider schema mismatch: %s", diff)
	}
}

func TestProviderSchemaFromJson_function(t *testing.T) {
	testCases := []struct {
		testName       string
		rawSchema      string
		expectedSchema ProviderSchema
	}{
		{
			"basic",
			`{
			"functions": {
				"example": {
				  "description": "Echoes given argument as result",
				  "summary": "Example function",
				  "return_type": "string",
				  "parameters": [
					{
					  "name": "input",
					  "description": "String to echo",
					  "type": "string"
					}
				  ]
				}
			  }
		}`,
			ProviderSchema{
				Resources:          map[string]*schema.BodySchema{},
				EphemeralResources: map[string]*schema.BodySchema{},
				DataSources:        map[string]*schema.BodySchema{},
				Functions: map[string]*schema.FunctionSignature{
					"example": {
						Description: "Echoes given argument as result",
						Detail:      "hashicorp/aws",
						ReturnType:  cty.String,
						Params: []function.Parameter{
							{
								Name:        "input",
								Description: "String to echo",
								Type:        cty.String,
							},
						},
						VarParam: nil,
					},
				},
				ListResources:   map[string]*schema.BodySchema{},
				ActionResources: map[string]*schema.BodySchema{},
			},
		},
		{
			"no parameters",
			`{
			"functions": {
				"example": {
				  "description": "Returns a string",
				  "summary": "Example function",
				  "return_type": "string",
				  "parameters": []
				}
			  }
		}`,
			ProviderSchema{
				Resources:          map[string]*schema.BodySchema{},
				EphemeralResources: map[string]*schema.BodySchema{},
				DataSources:        map[string]*schema.BodySchema{},
				Functions: map[string]*schema.FunctionSignature{
					"example": {
						Description: "Returns a string",
						Detail:      "hashicorp/aws",
						ReturnType:  cty.String,
						Params:      []function.Parameter{},
						VarParam:    nil,
					},
				},
				ListResources:   map[string]*schema.BodySchema{},
				ActionResources: map[string]*schema.BodySchema{},
			},
		},
		{
			"with variadic parameter",
			`{
			"functions": {
				"example": {
				  "description": "Echoes given argument as result",
				  "summary": "Example function",
				  "return_type": "string",
				  "parameters": [
					{
					  "name": "input",
					  "description": "String to echo",
					  "type": "string"
					}
				  ],
				  "variadic_parameter": {
					"name": "vars",
					"description": "Optional additional arguments",
					"type": "string"
				  }
				}
			  }
		}`,
			ProviderSchema{
				Resources:          map[string]*schema.BodySchema{},
				EphemeralResources: map[string]*schema.BodySchema{},
				DataSources:        map[string]*schema.BodySchema{},
				Functions: map[string]*schema.FunctionSignature{
					"example": {
						Description: "Echoes given argument as result",
						Detail:      "hashicorp/aws",
						ReturnType:  cty.String,
						Params: []function.Parameter{
							{
								Name:        "input",
								Description: "String to echo",
								Type:        cty.String,
							},
						},
						VarParam: &function.Parameter{
							Name:        "vars",
							Type:        cty.String,
							Description: "Optional additional arguments",
						},
					},
				},
				ListResources:   map[string]*schema.BodySchema{},
				ActionResources: map[string]*schema.BodySchema{},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%2d-%s", i, tc.testName), func(t *testing.T) {

			jsonSchema := &tfjson.ProviderSchema{}
			err := json.Unmarshal([]byte(tc.rawSchema), jsonSchema)
			if err != nil {
				t.Fatal(err)
			}
			providerAddr := addr.NewDefaultProvider("aws")

			ps := ProviderSchemaFromJson(jsonSchema, providerAddr)
			if diff := cmp.Diff(&tc.expectedSchema, ps, ctydebug.CmpOptions); diff != "" {
				t.Fatalf("provider schema mismatch: %s", diff)
			}
		})
	}

}
