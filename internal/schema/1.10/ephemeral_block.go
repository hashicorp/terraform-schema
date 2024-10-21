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

func ephemeralBlockSchema() *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "ephemeral"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "ephemeral",
			ScopeId:              refscope.EphemeralScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Ephemeral},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Ephemeral Resource Type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("An ephemeral block declares an ephemeral resource of a given type with a given local name. The name is " +
			"used to refer to this ephemeral resource from elsewhere in the same Terraform module, but has no significance " +
			"outside of the scope of a module."),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:         true,
				ForEach:       true,
				DynamicBlocks: true,
			},
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
							schema.Reference{OfScopeId: refscope.EphemeralScope},
							schema.Reference{OfScopeId: refscope.VariableScope},
							schema.Reference{OfScopeId: refscope.LocalScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. resources or data sources"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle": ephemeralLifecycleBlock(),
			},
		},
	}

	return bs
}

func ephemeralLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations to change default ephemeral resource behaviors during apply"),
		Body: &schema.BodySchema{
			Blocks: map[string]*schema.BlockSchema{
				"precondition": {
					Body: conditionBody(false),
				},
				"postcondition": {
					Body: conditionBody(true),
				},
			},
		},
	}
}

func conditionBody(enableSelfRefs bool) *schema.BodySchema {
	bs := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"condition": {
				Constraint: schema.AnyExpression{OfType: cty.Bool},
				IsRequired: true,
				Description: lang.Markdown("Condition, a boolean expression that should return `true` " +
					"if the intended assumption or guarantee is fulfilled or `false` if it is not."),
			},
			"error_message": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsRequired: true,
				Description: lang.Markdown("Error message to return if the `condition` isn't met " +
					"(evaluates to `false`)."),
			},
		},
	}

	if enableSelfRefs {
		bs.Extensions = &schema.BodyExtensions{
			SelfRefs: true,
		}
	}

	return bs
}
