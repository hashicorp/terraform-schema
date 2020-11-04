package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var datasourceBlockSchema = &schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{
			Name:        "type",
			Description: lang.PlainText("Data Source Type"),
			IsDepKey:    true,
		},
		{
			Name:        "name",
			Description: lang.PlainText("Reference Name"),
		},
	},
	Description: lang.PlainText("A data block requests that Terraform read from a given data source and export the result " +
		"under the given local name. The name is used to refer to this resource from elsewhere in the same " +
		"Terraform module, but has no significance outside of the scope of a module."),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"provider": {
				ValueType:   cty.DynamicPseudoType,
				Description: lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
				IsDepKey:    true,
			},
			"count": {
				ValueType:   cty.Number,
				Description: lang.Markdown("Number of instances of this data source, e.g. `3`"),
			},
			"depends_on": {
				ValueType:   cty.Set(cty.DynamicPseudoType),
				Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
			},
			"for_each": {
				ValueTypes: schema.ValueTypes{
					cty.Set(cty.DynamicPseudoType),
					cty.Map(cty.DynamicPseudoType),
				},
				Description: lang.Markdown("A set or a map where each item represents an instance of this data source"),
			},
		},
	},
}
