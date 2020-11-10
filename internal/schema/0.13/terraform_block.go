package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var terraformBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Terraform block used to configure some high-level behaviors of Terraform"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"required_version": {
				ValueType:  cty.String,
				IsOptional: true,
				Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
					"with this configuration, e.g. `~> 0.12`"),
			},
			"experiments": {
				ValueType:   cty.Set(cty.DynamicPseudoType),
				IsOptional:  true,
				Description: lang.Markdown("A set of experimental language features to enable"),
			},
		},
		Blocks: map[string]*schema.BlockSchema{
			"backend": {
				Description: lang.Markdown("Backend configuration which defines exactly where and how " +
					"operations are performed, where state snapshots are stored, etc."),
				Labels: []*schema.LabelSchema{
					{
						Name:        "type",
						Description: lang.Markdown("Backend Type"),
						IsDepKey:    true,
					},
				},
			},
			"provider_meta": {
				Description: lang.Markdown("Metadata to pass into a provider which supports this"),
				Labels: []*schema.LabelSchema{
					{
						Name:        "name",
						Description: lang.Markdown("Provider Name"),
						IsDepKey:    true,
					},
				},
			},
			"required_providers": {
				Description: lang.Markdown("What provider version to use within this configuration " +
					"and where to source it from"),
				Body: &schema.BodySchema{
					AnyAttribute: &schema.AttributeSchema{
						ValueTypes: schema.ValueTypes{
							cty.Object(map[string]cty.Type{
								"source":  cty.String,
								"version": cty.String,
							}),
							cty.String,
						},
						Description: lang.Markdown("Provider source and version constraint"),
					},
				},
			},
		},
	},
}
