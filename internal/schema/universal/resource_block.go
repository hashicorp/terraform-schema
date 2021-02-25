package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

var resourceBlockSchema = &schema.BlockSchema{
	Labels: []*schema.LabelSchema{
		{
			Name:        "type",
			Description: lang.PlainText("Resource Type"),
			IsDepKey:    true,
		},
		{
			Name:        "name",
			Description: lang.PlainText("Reference Name"),
		},
	},
	Description: lang.PlainText("A resource block declares a resource of a given type with a given local name. The name is " +
		"used to refer to this resource from elsewhere in the same Terraform module, but has no significance " +
		"outside of the scope of a module."),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"provider": {
				Expr:        schema.ExprConstraints{},
				IsOptional:  true,
				Description: lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
				IsDepKey:    true,
			},
		},
	},
}
