package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

var resourceBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.LabelStep{Index: 0},
			schema.LabelStep{Index: 1},
		},
		FriendlyName:        "resource",
		ScopeId:             refscope.ResourceScope,
		AsReference:         true,
		DependentBodyAsData: true,
		InferDependentBody:  true,
	},
	SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
	Labels: []*schema.LabelSchema{
		{
			Name:                   "type",
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
			Description:            lang.PlainText("Resource Type"),
			IsDepKey:               true,
			Completable:            true,
		},
		{
			Name:                   "name",
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			Description:            lang.PlainText("Reference Name"),
		},
	},
	Description: lang.PlainText("A resource block declares a resource of a given type with a given local name. The name is " +
		"used to refer to this resource from elsewhere in the same Terraform module, but has no significance " +
		"outside of the scope of a module."),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"provider": {
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfScopeId: refscope.ProviderScope},
				},
				IsOptional:  true,
				Description: lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
				IsDepKey:    true,
			},
		},
	},
}
