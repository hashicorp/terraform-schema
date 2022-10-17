package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func resourceBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
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
			Extensions: &schema.BodyExtensions{
				Count:         true,
				DynamicBlocks: true,
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
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle":   lifecycleBlock,
				"connection":  connectionBlock(v),
				"provisioner": provisionerBlock(v),
				"dynamic":     dynamicBlock,
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_6) {
		bs.Body.Attributes["for_each"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.TraversalExpr{OfType: cty.Set(cty.DynamicPseudoType)},
				schema.TraversalExpr{OfType: cty.Map(cty.DynamicPseudoType)},
				schema.LiteralTypeExpr{Type: cty.Set(cty.DynamicPseudoType)},
				schema.LiteralTypeExpr{Type: cty.Map(cty.DynamicPseudoType)},
			},
			IsOptional:  true,
			Description: lang.Markdown("A set or a map where each item represents an instance of this resource"),
		}
	}

	return bs
}

var dynamicBlock = &schema.BlockSchema{
	Description: lang.Markdown("A dynamic block to produce blocks dynamically " +
		"by iterating over a given complex value"),
	Labels: []*schema.LabelSchema{
		{
			Name:        "Block Name",
			IsDepKey:    true,
			Completable: true,
		},
	},
	DependentBody: map[schema.SchemaKey]*schema.BodySchema{
		schema.NewSchemaKey(schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{
					Index: 0,
					Value: "setting",
				},
			},
		}): {
			Attributes: map[string]*schema.AttributeSchema{
				"for_each": {
					Expr: schema.ExprConstraints{
						schema.TraversalExpr{OfType: cty.List(cty.DynamicPseudoType)},
						schema.TraversalExpr{OfType: cty.Set(cty.DynamicPseudoType)},
						schema.TraversalExpr{OfType: cty.Map(cty.DynamicPseudoType)},
						schema.LiteralTypeExpr{Type: cty.List(cty.DynamicPseudoType)},
						schema.LiteralTypeExpr{Type: cty.Set(cty.DynamicPseudoType)},
						schema.LiteralTypeExpr{Type: cty.Map(cty.DynamicPseudoType)},
					},
					IsRequired:  true,
					Description: lang.Markdown("A set or a map where each item represents a block instance"),
				},
				"iterator": {
					Expr:       schema.LiteralTypeOnly(cty.String),
					IsOptional: true,
					Description: lang.Markdown("The name of a temporary variable that represents the current " +
						"element of the complex value. Defaults to the label of the dynamic block."),
				},
				"labels": {
					Expr: schema.ExprConstraints{
						schema.ListExpr{
							Elem: schema.ExprConstraints{
								schema.LiteralTypeExpr{Type: cty.String},
								schema.TraversalExpr{OfType: cty.String},
							},
						},
					},
					IsOptional: true,
					Description: lang.Markdown("A list of strings that specifies the block labels, " +
						"in order, to use for each generated block."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"content": {
					Description: lang.PlainText("The body of each generated block"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"name": {
								IsRequired: true,
								Expr: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.String, OfScopeId: lang.ScopeId("builtin")},
									schema.LiteralTypeExpr{Type: cty.String},
								},
							},
							"namespace": {
								IsRequired: true,
								Expr: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.String},
									schema.LiteralTypeExpr{Type: cty.String},
								},
							},
							"resource": {
								IsOptional: true,
								Expr: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.String},
									schema.LiteralTypeExpr{Type: cty.String},
								},
							},
							"value": {
								IsRequired: true,
								Expr: schema.ExprConstraints{
									schema.TraversalExpr{OfType: cty.String},
									schema.LiteralTypeExpr{Type: cty.String},
								},
							},
						},
					},
				},
			},
		},
	},
}

var lifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations to change default resource behaviours during apply"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"create_before_destroy": {
				Expr:       schema.LiteralTypeOnly(cty.Bool),
				IsOptional: true,
				Description: lang.Markdown("Whether to reverse the default order of operations (destroy -> create) during apply " +
					"when the resource requires replacement (cannot be updated in-place)"),
			},
			"prevent_destroy": {
				Expr:       schema.LiteralTypeOnly(cty.Bool),
				IsOptional: true,
				Description: lang.Markdown("Whether to prevent accidental destruction of the resource and cause Terraform " +
					"to reject with an error any plan that would destroy the resource"),
			},
			"ignore_changes": {
				Expr: schema.ExprConstraints{
					schema.TupleConsExpr{},
					schema.KeywordExpr{
						Keyword: "all",
						Description: lang.Markdown("Ignore all attributes, which means that Terraform can create" +
							" and destroy the remote object but will never propose updates to it"),
					},
				},
				IsOptional:  true,
				Description: lang.Markdown("A set of fields (references) of which to ignore changes to, e.g. `tags`"),
			},
		},
	},
}
