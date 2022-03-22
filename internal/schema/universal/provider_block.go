package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
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
	SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
	Labels: []*schema.LabelSchema{
		{
			Name:                   "name",
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			Description:            lang.PlainText("Provider Name"),
			IsDepKey:               true,
			Completable:            true,
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
		},
	},
}
