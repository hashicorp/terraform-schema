// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func providerBlockSchema() *schema.BlockSchema {
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
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				ForEach: true,
			},
			// If we add this here, the dependent body schema won't override the config block and we wont get more
			// specific completions for the provider config block, due to skipping here:
			// https://github.com/hashicorp/hcl-lang/blob/main/decoder/internal/schemahelper/block_schema.go#L52
			// Blocks: map[string]*schema.BlockSchema{
			// 	"config": {
			// 		Body:        &schema.BodySchema{},
			// 		Description: lang.Markdown("Provider configuration"),
			// 		MaxItems:    1,
			// 	},
			// },
		},
	}
}
