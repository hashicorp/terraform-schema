package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func moduleBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "module"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "module",
			ScopeId:      refscope.ModuleScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("Module block to call a locally or remotely stored module"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"source": {
					Expr: schema.LiteralTypeOnly(cty.String),
					Description: lang.Markdown("Source where to load the module from, " +
						"a local directory (e.g. `./module`) or a remote address - e.g. " +
						"`hashicorp/consul/aws` (Terraform Registry address) or " +
						"`github.com/hashicorp/example` (GitHub)"),
					IsRequired:             true,
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					CompletionHooks: lang.CompletionHooks{
						{
							Name: "CompleteLocalModuleSources",
						},
						{
							Name: "CompleteRegistryModuleSources",
						},
					},
				},
				"version": {
					Expr:       schema.LiteralTypeOnly(cty.String),
					IsOptional: true,
					Description: lang.Markdown("Constraint to set the version of the module, e.g. `~> 1.0`." +
						" Only applicable to modules in a module registry."),
					CompletionHooks: lang.CompletionHooks{
						{
							Name: "CompleteRegistryModuleVersions",
						},
					},
				},
				"providers": {
					Expr: schema.ExprConstraints{
						schema.MapExpr{
							Name: "map of provider references",
							Elem: schema.ExprConstraints{
								schema.TraversalExpr{OfScopeId: refscope.ProviderScope},
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Explicit mapping of providers which the module uses"),
				},
				"count": {
					Expr: schema.ExprConstraints{
						schema.TraversalExpr{OfType: cty.Number},
						schema.LiteralTypeExpr{Type: cty.Number},
					},
					IsOptional:  true,
					Description: lang.Markdown("Number of instances of this module, e.g. `3`"),
				},
				"for_each": {
					Expr: schema.ExprConstraints{
						schema.TraversalExpr{OfType: cty.Set(cty.DynamicPseudoType)},
						schema.TraversalExpr{OfType: cty.Map(cty.DynamicPseudoType)},
						schema.LiteralTypeExpr{Type: cty.Set(cty.DynamicPseudoType)},
						schema.LiteralTypeExpr{Type: cty.Map(cty.DynamicPseudoType)},
					},
					IsOptional:  true,
					Description: lang.Markdown("A set or a map where each item represents an instance of this module"),
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
}
