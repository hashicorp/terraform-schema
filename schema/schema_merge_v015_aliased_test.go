// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

var expectedMergedSchema_v015_aliased = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"data": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "type",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				},
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Constraint: schema.AnyExpression{OfType: cty.Number},
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
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"config1": {
							IsOptional: true,
							Constraint: schema.Object{
								Attributes: schema.ObjectAttributes{
									"first": {
										IsOptional: true,
										Constraint: schema.AnyExpression{OfType: cty.String},
									},
									"second": {
										IsOptional: true,
										Constraint: schema.AnyExpression{OfType: cty.Number},
									},
									"third": {
										IsOptional: true,
										Constraint: schema.Object{
											Attributes: schema.ObjectAttributes{
												"nested": {
													IsOptional: true,
													Constraint: schema.AnyExpression{OfType: cty.String},
												},
											},
										},
									},
								},
							},
						},
						"config2": {
							IsOptional: true,
							Constraint: schema.List{
								Elem: schema.Object{
									Attributes: schema.ObjectAttributes{
										"first": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.String},
										},
										"second": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Number},
										},
										"third": {
											IsOptional: true,
											Constraint: schema.Object{
												Attributes: schema.ObjectAttributes{
													"nested": {
														IsOptional: true,
														Constraint: schema.AnyExpression{OfType: cty.String},
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
						"config3": {
							IsOptional: true,
							Constraint: schema.Set{
								Elem: schema.Object{
									Attributes: schema.ObjectAttributes{
										"first": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.String},
										},
										"second": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Number},
										},
										"third": {
											IsOptional: true,
											Constraint: schema.Object{
												Attributes: schema.ObjectAttributes{
													"nested": {
														IsOptional: true,
														Constraint: schema.AnyExpression{OfType: cty.String},
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
						"config4": {
							IsOptional: true,
							Constraint: schema.Map{
								Elem: schema.Object{
									Attributes: schema.ObjectAttributes{
										"first": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.String},
										},
										"second": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Number},
										},
										"third": {
											IsOptional: true,
											Constraint: schema.Object{
												Attributes: schema.ObjectAttributes{
													"nested": {
														IsOptional: true,
														Constraint: schema.AnyExpression{OfType: cty.String},
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
						"workspace": {
							IsOptional: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
					Detail: "hashicorp/hashicup",
				},
			},
		},
		"provider": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {
						Constraint: schema.LiteralType{Type: cty.String},
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
				{
					Name:                   "type",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				},
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Constraint: schema.AnyExpression{OfType: cty.Number},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
		"module": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"source": {
						Constraint:             schema.LiteralType{Type: cty.String},
						IsRequired:             true,
						IsDepKey:               true,
						SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					},
					"version": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
	},
}
