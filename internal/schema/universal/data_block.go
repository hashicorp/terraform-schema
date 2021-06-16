package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
)

var datasourceBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.StaticStep{Name: "data"},
			schema.LabelStep{Index: 0},
			schema.LabelStep{Index: 1},
		},
		FriendlyName:        "datasource",
		ScopeId:             refscope.DataScope,
		AsReference:         true,
		DependentBodyAsData: true,
		InferDependentBody:  true,
	},
	Labels: []*schema.LabelSchema{
		{
			Name:        "type",
			Description: lang.PlainText("Data Source Type"),
			IsDepKey:    true,
			Completable: true,
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
