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
			},
		},
	}
}
