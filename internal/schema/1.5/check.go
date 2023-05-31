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

func checkBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				Description:            lang.PlainText("Local Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.Markdown("Check customized infrastructure requirements to provide ongoing and continuous verification."),
		Body: &schema.BodySchema{
			HoverURL: "https://developer.hashicorp.com/terraform/language/checks",
			Blocks: map[string]*schema.BlockSchema{
				"data": scopedDataBlock(),
				"assert": {
					MinItems: 1,
					Description: lang.Markdown(`Assertion, which does not affect Terraform's execution of an operation. ` +
						`A failed assertion reports a warning without halting the ongoing operation.`),
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
			},
		},
	}
}

func scopedDataBlock() *schema.BlockSchema {
	bs := &schema.BlockSchema{
		// TODO: Address: &schema.BlockAddrSchema{},
		// See https://github.com/hashicorp/terraform-schema/issues/234
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Data Source Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				Description:            lang.PlainText("Reference Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.PlainText("A data block requests that Terraform read from a given data source and export the result " +
			"under the given locally scoped name (i.e. only within the `check` block)."),
		MaxItems: 1,
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:             schema.Reference{OfScopeId: refscope.ProviderScope},
					IsOptional:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.ModuleScope},
							schema.Reference{OfScopeId: refscope.ResourceScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. resources or data sources"),
				},
			},
		},
	}

	return bs
}
