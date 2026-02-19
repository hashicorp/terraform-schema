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

func moduleBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "module",
			ScopeId:              refscope.ModuleScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Module Source"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("Defines a `mock` module configuration to test policies targeting module usage."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"expect_failure": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description:  lang.Markdown("Expect test to fail"),
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
							"address": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The `address` is the internal, logical path used by Terraform to reference resources within a configuration for commands like state management"),
							},
							"source": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("The `source` is the external location where Terraform physically finds and downloads the module code"),
							},
							"version": &schema.AttributeSchema{
								Constraint:  schema.AnyExpression{OfType: cty.String},
								Description: lang.Markdown("Version of the module"),
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Meta attributes of the module"),
				},
			},
		},
	}
}
