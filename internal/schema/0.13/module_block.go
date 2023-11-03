// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func moduleBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "module"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "module",
			ScopeId:      refscope.ModuleScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("Module block to call a locally or remotely stored module"),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:   true,
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"source": {
					Constraint: schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Source where to load the module from, " +
						"a local directory (e.g. `./module`) or a remote address - e.g. " +
						"`hashicorp/consul/aws` (Terraform Registry address) or " +
						"`github.com/hashicorp/example` (GitHub)"),
					IsRequired:             true,
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					CompletionHooks: lang.CompletionHooks{
						{
							Name: "CompleteLocalModuleSources",
						},
						{
							Name: "CompleteRegistryModuleSources",
						},
					},
				},
				"version": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("Constraint to set the version of the module, e.g. `~> 1.0`." +
						" Only applicable to modules in a module registry."),
					CompletionHooks: lang.CompletionHooks{
						{
							Name: "CompleteRegistryModuleVersions",
						},
					},
				},
				"providers": {
					Constraint: schema.Map{
						Name: "map of provider references",
						Elem: schema.Reference{OfScopeId: refscope.ProviderScope},
					},
					IsOptional:  true,
					Description: lang.Markdown("Explicit mapping of providers which the module uses"),
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.ModuleScope},
							schema.Reference{OfScopeId: refscope.ResourceScope},
							schema.Reference{OfScopeId: refscope.VariableScope},
							schema.Reference{OfScopeId: refscope.LocalScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
				},
			},
		},
	}
}
