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

func listBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "list"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "list",
			ScopeId:              refscope.ListScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("List Type"),
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
		Description: lang.PlainText("A list block defines a mechanism to retrieve collections of resources. " +
			"It specifies the type of resource to be listed and is uniquely identified by the name."),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:   true,
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:             schema.Reference{OfScopeId: refscope.ProviderScope},
					IsRequired:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
				"include_resource": {
					Constraint:   schema.AnyExpression{OfType: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.False},
					IsOptional:   true,
					Description: lang.Markdown("By default, the results of a list resource only include the identities of the discovered resources. " +
						"If it is marked true then the provider should include the resource data in the result."),
				},
				"limit": {
					Constraint: schema.AnyExpression{OfType: cty.Number},
					IsOptional: true,
					Description: lang.Markdown("Limit is an optional value that can be used to limit the " +
						"number of results returned by the list resource."),
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.ListScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other list"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"config": {
					Description: lang.Markdown("Filters specific to the list type"),
					MaxItems:    1,
				},
			},
		},
	}
}
