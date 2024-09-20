// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func mockProviderBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Description: lang.PlainText("In Terraform tests, you can mock a provider with the mock_provider block. Mock providers return the same schema as the original provider and you can pass the mocked provider to your tests in place of the matching provider. All resources and data sources retrieved by a mock provider will set the relevant values from the configuration, and generate fake data for any computed attributes."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Alias for using the same provider with different configurations for different resources, e.g. `mock`"),
				},
				"source": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Path to a directory that includes dedicated mock data files (*.tfmock.hcl). These can be used to share mock provider data between tests. You can combine the source attribute with directly nested mock_resource and mock_data blocks. If the source location and a directly nested block describe the same resource or data source then the directly nested block takes precedence."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"mock_resource":     mockResourceBlockSchema(),
				"mock_data":         mockDataBlockSchema(),
				"override_resource": overrideResourceBlockSchema(),
				"override_data":     overrideDataBlockSchema(),
			},
		},
	}
}
