// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func providerBlockSchema() *schema.BlockSchema {
	/*
		Reference: https://github.com/hashicorp/terraform/blob/main/internal/stacks/stackconfig/provider_config.go
		TODO:
			- Source better descriptions
			- config should autocomplete from the specified provider
			- Verify all attributes are added here
			- for_each
	*/

	return &schema.BlockSchema{
		Description:            lang.PlainText("A Stack provider block is used to specify a provider configuration"),
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Provider Name"),
				// TODO: this is the index, so is it a depkey?
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"config": {
					Constraint: schema.Map{
						Name: "map of configuration",
						Elem: schema.Reference{OfScopeId: refscope.ProviderScope},
					},
					IsOptional:  true,
					Description: lang.Markdown("Explicit mapping of configuration for the provider"),
				},
			},
		},
	}
}
