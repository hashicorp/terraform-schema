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
					Description: lang.Markdown("The Terraform command to be run in this test. By default, each run block executes with `command = apply` instructing Terraform to execute a complete apply operation against your configuration. Replacing the command value with `command = plan` instructs Terraform to _not_ create new infrastructure for this run block. This allows test authors to validate logical operations and custom conditions within their infrastructure in a process analogous to unit testing."),
					IsOptional:  true,
				},
				"providers": {
					Constraint: schema.Object{
						Attributes: schema.ObjectAttributes{},
					},
					Description: lang.Markdown("Explicit mapping of providers to use for this run block. If not set, each provider you specify is directly available within each run block."),
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
					Description: lang.Markdown("Takes a list of checkable objects (resources, data sources, check blocks, input variables, and outputs) that should fail their custom conditions. The test passes if the checkable objects you specify report an issue, and the test fails overall if they do not. Read more on [expect_failures](https://developer.hashicorp.com/terraform/language/tests#expecting-failures)"),
					IsOptional:  true,
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"plan_options": {
					Description: lang.Markdown("Allows to customize the planning mode and options typically needed to edit via command-line flags and options."),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"mode": {
								Constraint: schema.OneOf{
									schema.Keyword{
										Keyword:     "normal",
										Description: lang.Markdown("Default planning mode which will run an implicit in-memory refresh as part of the plan"),
									},
									schema.Keyword{
										Keyword:     "refresh-only",
										Description: lang.Markdown("Allows to safely refresh Terraform state without making any modifications to your infrastructure"),
									},
								},
								Description: lang.Markdown("Whether to run Terraform normal or in refresh-only mode"),
								IsOptional:  true,
							},
							"refresh": {
								Constraint:   schema.LiteralType{Type: cty.Bool},
								Description:  lang.Markdown("When set to false, it disables the default behavior of synchronizing the Terraform state with remote objects before checking for configuration changes."),
								DefaultValue: schema.DefaultValue{Value: cty.BoolVal(true)},
								IsOptional:   true,
							},
							"replace": {
								Constraint: schema.List{
									Elem: schema.Reference{OfScopeId: refscope.ResourceScope},
								},
								Description: lang.Markdown("Instructs Terraform to plan to replace the resource instance with the given address."),
								IsOptional:  true,
							},
							"target": {
								Constraint: schema.List{
									Elem: schema.Reference{OfScopeId: refscope.ResourceScope},
								},
								Description: lang.Markdown("Instructs Terraform to focus its planning efforts only on resource instances which match the given address and on any objects that those instances depend on."),
								IsOptional:  true,
							},
						},
					},
				},
				"assert": {
					Description: lang.Markdown("At the conclusion of a Terraform test command execution, Terraform presents any failed assertions as part of a tests passed or failed status."),
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
					Description: lang.Markdown("Used to modify the module that a given run block executes"),
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
