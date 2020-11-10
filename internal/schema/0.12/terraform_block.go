package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func terraformBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Description: lang.Markdown("Terraform block used to configure some high-level behaviors of Terraform"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"required_version": {
					ValueType:   cty.String,
					Description: lang.Markdown("Constraint to specify which versions of Terraform can be used "+
						"with this configuration, e.g. `~> 0.12`"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"backend": {
					Description: lang.Markdown("Backend configuration which defines exactly where and how "+
						"operations are performed, where state snapshots are stored, etc."),
					Labels: []*schema.LabelSchema{
						{
							Name:        "type",
							Description: lang.Markdown("Backend Type"),
							IsDepKey:    true,
						},
					},
					MaxItems: 1,
				},
				"required_providers": {
					Description: lang.Markdown("What provider version to use within this configuration"),
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							ValueType:   cty.String,
							Description: lang.Markdown("Version constraint"),
						},
					},
					MaxItems: 1,
				},
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_18) {
		bs.Body.Attributes["experiments"] = &schema.AttributeSchema{
			ValueType:   cty.Set(cty.DynamicPseudoType),
			Description: lang.Markdown("A set of experimental language features to enable"),
		}
	}

	if v.GreaterThanOrEqual(v0_12_20) {
		bs.Body.Blocks["required_providers"].Body = &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				ValueTypes: schema.ValueTypes{
					cty.Object(map[string]cty.Type{
						"version": cty.String,
					}),
					cty.String,
				},
				Description: lang.Markdown("Version constraint"),
			},
		}
	}

	return bs
}
