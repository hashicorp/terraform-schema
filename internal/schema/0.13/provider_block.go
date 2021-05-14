package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

var providerBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.LabelStep{Index: 0},
			schema.AttrValueStep{Name: "alias", IsOptional: true},
		},
		FriendlyName: "provider",
		ScopeId:      refscope.ProviderScope,
		AsReference:  true,
	},
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Provider Name"),
			IsDepKey:    true,
		},
	},
	Description: lang.PlainText("A provider block is used to specify a provider configuration"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"alias": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Alias for using the same provider with different configurations for different resources, e.g. `eu-west`"),
			},
			"version": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				IsDeprecated: true,
				Description: lang.Markdown("Specifies a version constraint for the provider. e.g. `~> 1.0`.\n" +
					"**DEPRECATED:** Use `required_providers` block to manage provider version instead."),
			},
		},
	},
}
