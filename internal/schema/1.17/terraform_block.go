// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Blocks["state_store"] = &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.StateStore},
		Labels: []*schema.LabelSchema{
			{
				Name:        "name",
				Description: lang.Markdown("State Store Name"),
				IsDepKey:    true,
				Completable: true,
				SemanticTokenModifiers: lang.SemanticTokenModifiers{
					tokmod.Name,
					lang.TokenModifierDependent,
				},
			},
		},
		Description: lang.PlainText("A state_store block describing where and/or how the state is managed"),
		MaxItems:    1,
		Body: &schema.BodySchema{
			Blocks: map[string]*schema.BlockSchema{
				"provider": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
					Labels: []*schema.LabelSchema{
						{
							Name: "name",
							SemanticTokenModifiers: lang.SemanticTokenModifiers{
								tokmod.Name,
								lang.TokenModifierDependent,
							},
							Description: lang.PlainText("Provider Name"),
							IsDepKey:    true,
							Completable: true,
						},
					},
					Description: lang.PlainText("A provider block is used to specify a provider configuration"),
				},
			},
		},
	}

	return bs
}
