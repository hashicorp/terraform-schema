// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func runBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Run},
		Description:            lang.PlainText("A run block represents a single Terraform command to be executed and a set of validations to run after the command."),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Run Name"),
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"command": {
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "apply",
							Description: lang.Markdown("The operation is an apply operation."),
						},
						schema.Keyword{
							Keyword:     "plan",
							Description: lang.Markdown("The operation is a plan operation."),
						},
					},
					Description: lang.Markdown(""), // TODO!
					IsOptional:  true,
				},
				"providers": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					Description: lang.Markdown(""), // TODO!
					IsOptional:  true,
				},
				"expect_failures": {
					Constraint: schema.List{
						Elem: schema.OneOf{
							// TODO? check blocks
							schema.Reference{OfScopeId: refscope.ResourceScope},
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.VariableScope},
							schema.Reference{OfScopeId: refscope.OutputScope},
						},
					},
					Description: lang.Markdown(""), // TODO!
					IsOptional:  true,
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"plan_options": { // TODO!
					Description: lang.Markdown(""),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"mode": {
								Constraint: schema.OneOf{
									schema.Keyword{
										Keyword:     "normal",
										Description: lang.Markdown(""), // TODO!
									},
									schema.Keyword{
										Keyword:     "refresh-only",
										Description: lang.Markdown(""), // TODO!
									},
								},
								Description: lang.Markdown(""), // TODO!
								IsOptional:  true,
							},
							"refresh": {
								Constraint:   schema.LiteralType{Type: cty.Bool},
								Description:  lang.Markdown(""), // TODO!
								DefaultValue: schema.DefaultValue{Value: cty.BoolVal(true)},
								IsOptional:   true,
							},
							"replace": {
								Constraint: schema.List{
									Elem: schema.Reference{OfScopeId: refscope.ResourceScope},
								},
								Description: lang.Markdown(""), // TODO!
								IsOptional:  true,
							},
							"target": {
								Constraint: schema.List{
									Elem: schema.Reference{OfScopeId: refscope.ResourceScope},
								},
								Description: lang.Markdown(""), // TODO!
								IsOptional:  true,
							},
						},
					},
				},
				"assert": {
					Description: lang.Markdown(""), // TODO!
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Constraint:  schema.AnyExpression{OfType: cty.Bool},
								IsRequired:  true,
								Description: lang.Markdown("Condition to meet for the check to pass (any expression which evaluates to boolean)"),
							},
							"error_message": {
								Constraint:  schema.AnyExpression{OfType: cty.String},
								IsRequired:  true,
								Description: lang.Markdown("Text that Terraform will include as part of error messages when it detects an unmet condition"),
							},
						},
					},
				},
				"variables": variablesBlockSchema(),
				"module": {
					Description: lang.Markdown(""), // TODO!
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"source": {
								Constraint: schema.LiteralType{Type: cty.String},
								Description: lang.Markdown("Source where to load the module from, " +
									"a local directory (e.g. `./module`) or a remote address - e.g. " +
									"`hashicorp/consul/aws` (Terraform Registry address)"),
								IsRequired:             true,
								IsDepKey:               true,
								SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
							},
							"version": {
								Constraint: schema.LiteralType{Type: cty.String},
								IsOptional: true,
								Description: lang.Markdown("Constraint to set the version of the module, e.g. `~> 1.0`." +
									" Only applicable to modules in a module registry."),
							},
						},
					},
				},
			},
		},
	}
}
