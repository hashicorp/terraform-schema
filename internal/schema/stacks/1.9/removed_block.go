// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func removedBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Declaration to specify what component to remove from the stack"),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"from": {
					Description: lang.Markdown("Address of the component to be removed"),
					IsRequired:  true,
					Constraint:  schema.Reference{OfScopeId: refscope.ComponentScope}, // TODO: this component would not exist in the config anymore, only in state (i.e. the component block has to be removed from the config)
				},
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
				"version": {
					Description: lang.Markdown("Accepts a comma-separated list of version constraints for registry modules"),
					IsOptional:  true,
					Constraint: schema.List{
						Elem: schema.AnyExpression{OfType: cty.String}, // TODO: comma separated list
					},
				},
				"providers": {
					Description: lang.Markdown(" A mapping of provider names to providers declared in the stack configuration. Providers must be declared in the top level of the stack and passed into each component in the stack. Components cannot configure their own providers"),
					IsOptional:  true,
					Constraint: schema.Map{
						Name: "map of provider references",
						Elem: schema.Reference{OfScopeId: refscope.ProviderScope},
					},
				},
			},
		},
	}
}
