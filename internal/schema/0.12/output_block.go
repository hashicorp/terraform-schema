package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var outputBlockSchema = &schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Output Name"),
		},
	},
	Description: lang.PlainText("Output value for consumption by another module or a human interacting via the UI"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"description": {
				ValueType:   cty.String,
				IsOptional:  true,
				Description: lang.PlainText("Human-readable description of the output (for documentation and UI)"),
			},
			"value": {
				ValueType:   cty.DynamicPseudoType,
				IsRequired:  true,
				Description: lang.PlainText("Value, typically a reference to an attribute of a resource or a data source"),
			},
			"sensitive": {
				ValueType:   cty.Bool,
				IsOptional:  true,
				Description: lang.PlainText("Whether the output contains sensitive material and should be hidden in the UI"),
			},
			"depends_on": {
				ValueType:   cty.Set(cty.DynamicPseudoType),
				IsOptional:  true,
				Description: lang.PlainText("Set of references to hidden dependencies (e.g. resources or data sources)"),
			},
		},
	},
}
