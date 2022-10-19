package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func datasourceBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
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
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Data Source Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				Description:            lang.PlainText("Reference Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.PlainText("A data block requests that Terraform read from a given data source and export the result " +
			"under the given local name. The name is used to refer to this resource from elsewhere in the same " +
			"Terraform module, but has no significance outside of the scope of a module."),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:   true,
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Expr: schema.ExprConstraints{
						schema.TraversalExpr{OfScopeId: refscope.ProviderScope},
					},
					IsOptional:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
				"depends_on": {
					Expr: schema.ExprConstraints{
						schema.TupleConsExpr{
							Name: "set of references",
							AnyElem: schema.ExprConstraints{
								schema.TraversalExpr{OfScopeId: refscope.DataScope},
								schema.TraversalExpr{OfScopeId: refscope.ModuleScope},
								schema.TraversalExpr{OfScopeId: refscope.ResourceScope},
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
				},
			},
		},
	}

	return bs
}
