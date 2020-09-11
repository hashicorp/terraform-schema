package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var providerBlockSchema = &schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Provider Name"),
			IsDepKey:    true,
		},
	},
	Description: lang.PlainText("A provider block is used to specify a provider configuration. The body of the block (between " +
		"{ and }) contains configuration arguments for the provider itself. Most arguments in this section are " +
		"specified by the provider itself."),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"alias": {
				ValueType:   cty.String,
				Description: lang.PlainText("Alias for using the same provider with different configurations for different resources"),
			},
		},
	},
}
