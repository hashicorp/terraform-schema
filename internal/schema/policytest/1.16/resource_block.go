// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func resourceBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "resource"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:        "resource",
			ScopeId:             refscope.ResourceScope,
			AsReference:         true,
			DependentBodyAsData: true,
			InferDependentBody:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Resource Type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("Defines a specific infrastructure resource to be evaluated as a test case. It consists of a resource type and a unique name used to organize passing or failing scenarios."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"expect_failure": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Expect test to fail"),
				},
				"skip": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Whether to skip the test"),
				},
				"attrs": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					IsOptional:  true,
					Description: lang.Markdown("Specify the values that should be returned for specific attributes"),
				},
				"meta": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{
							"resource_type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Resource type"),
							},
							"provider_type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Provider powering the resource"),
							},
							"tfe_workspace": &schema.AttributeSchema{
								Constraint: schema.Object{
									Attributes: schema.ObjectAttributes{},
								},
								Description: lang.Markdown("Workspace config for which the resource belongs to"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Meta attributes of the resource"),
				},
			},
		},
	}
}
