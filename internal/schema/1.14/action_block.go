// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func actionBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "action"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:        "action",
			ScopeId:             refscope.ActionScope,
			AsReference:         true,
			DependentBodyAsData: true,
			InferDependentBody:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Action Type"),
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
		Description: lang.PlainText("An action block declares an action of a given type with a given local name. " +
			"The name is used to refer to this action elsewhere in the same Terraform module."),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:   true,
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:             schema.Reference{OfScopeId: refscope.ProviderScope},
					IsOptional:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"config": {
					Description: lang.Markdown("Provider specific action configuration"),
					MaxItems:    1,
				},
			},
		},
	}
}
