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

func componentBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Component represents the declaration of a single component within a particular Terraform Stack. Components are the most important object in a stack configuration, just as resources are the most important object in a Terraform module: each one refers to a Terraform module that describes the infrastructure that the component is 'made of'."),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "component"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "component",
			ScopeId:      refscope.ComponentScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Component Name"),
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"source": {
					Description:            lang.Markdown("The Terraform module location to load the Component from: a local directory (e.g. `./modules`), a git repository (e.g. `github.com/acme/infra/core`, `git::https://vcs.acme.com/acme/infra//core`), or a registry module (e.g. `acme-public/coreinfra/aws`, `app.terraform.io/acme/core-infra/aws`)"),
					IsRequired:             true,
					IsDepKey:               true,
					Constraint:             schema.LiteralType{Type: cty.String},
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					CompletionHooks: lang.CompletionHooks{
						{Name: "CompleteLocalModuleSources"},
					},
				},
				"inputs": {
					Description: lang.Markdown("A mapping of module input variable names to values. The keys of this map must correspond to the Terraform variable names in the module defined by source. Can be any Terraform expression, and can refer to anything which is in scope, including input variables, component outputs, the `each` object, and provider configurations"),
					IsOptional:  true,
					Constraint: schema.Map{
						Name: "map of input references",
						Elem: schema.AnyExpression{OfType: cty.DynamicPseudoType},
					},
				},
				"version": {
					Description: lang.Markdown("Accepts a comma-separated list of version constraints for registry modules"),
					IsOptional:  true,
					Constraint: schema.List{
						Elem: schema.AnyExpression{OfType: cty.String}, // TODO: comma separated list
					},
				},
				"providers": {
					Description: lang.Markdown("A mapping of provider names to providers declared in the stack configuration. Providers must be declared in the top level of the stack and passed into each component in the stack. Components cannot configure their own providers"),
					IsOptional:  true,
					Constraint: schema.Map{
						Name: "map of provider references",
						Elem: schema.Reference{OfScopeId: refscope.ProviderScope},
					},
				},
				"depends_on": {
					Description: lang.Markdown("Optionally specify explicit dependencies for components in a stack configuration, which must also be used when determining an order of operations for components"),
					IsOptional:  true,
					Constraint: schema.List{
						// TODO: This eventually will support embedded stack references
						Elem: schema.Reference{OfScopeId: refscope.ComponentScope},
					},
				},
			},
		},
	}
}
