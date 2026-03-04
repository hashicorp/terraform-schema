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
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "resource",
			ScopeId:              refscope.ResourceScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "resource_type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Resource Type"),
			},
			{
				Name:                   "test_case_name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Test Case Name"),
			},
		},
		Description: lang.PlainText("Defines a specific infrastructure resource to be evaluated as a test case. It consists of a resource type and a unique name used to organize passing or failing scenarios."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"expect_failure": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("If `true`, the test passes only if the policy engine rejects the resource"),
				},
				"skip": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("If `true`, this resource is used only as a dependency/reference for other resources and is not evaluated as a standalone test case. Cannot be used with `expect_failure`"),
				},
				"attrs": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					Description: lang.Markdown("A map of arguments that simulate the resource configuration"),
					IsOptional:  true,
				},
				"meta": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{
							"resource_type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Type of resource (“aws_s3_bucket”, “azurerm_managed_disk”)"),
							},
							"provider_type": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Provider of the resource"),
							},
							"tfe_workspace": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Information on the workspace. The only key available within this is ‘tags’. Please see this section for how to use this meta attribute"),
							},
							"address": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Address of the resource within Terraform"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Mocks the `required_providers` information"),
				},
			},
		},
	}
}
