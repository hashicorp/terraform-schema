package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_1_mod "github.com/hashicorp/terraform-schema/internal/schema/1.1"
)

var v1_2 = version.Must(version.NewVersion("1.2.0"))

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_1_mod.ModuleSchema(v)
	if v.GreaterThanOrEqual(v1_2) {
		bs.Blocks["data"].Body.Blocks = map[string]*schema.BlockSchema{
			"lifecycle": datasourceLifecycleBlock,
		}
		bs.Blocks["resource"].Body.Blocks["lifecycle"] = resourceLifecycleBlock
		bs.Blocks["output"].Body.Blocks = map[string]*schema.BlockSchema{
			"lifecycle": outputLifecycleBlock,
		}
	}

	return bs
}

var datasourceLifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations to set validity conditions of the datasource"),
	Body: &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"precondition": {
				Body: conditionBody,
			},
			"postcondition": {
				Body: conditionBody,
			},
		},
	},
}

var outputLifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations, to set a validity condition of the output"),
	Body: &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"precondition": {
				Body: conditionBody,
			},
		},
	},
}

var resourceLifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations to change default resource behaviours during plan or apply"),
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
		Blocks: map[string]*schema.BlockSchema{
			"precondition": {
				Body: conditionBody,
			},
			"postcondition": {
				Body: conditionBody,
			},
		},
	},
}

var conditionBody = &schema.BodySchema{
	Attributes: map[string]*schema.AttributeSchema{
		"condition": {
			Expr: schema.ExprConstraints{
				schema.TraversalExpr{OfType: cty.Bool},
				schema.LiteralTypeExpr{Type: cty.Bool},
			},
			IsRequired: true,
			Description: lang.Markdown("Condition, a boolean expression that should return `true` " +
				"if the intended assumption or guarantee is fulfilled or `false` if it is not."),
		},
		"error_message": {
			Expr: schema.ExprConstraints{
				schema.TraversalExpr{OfType: cty.String},
				schema.LiteralTypeExpr{Type: cty.String},
			},
			IsRequired: true,
			Description: lang.Markdown("Error message to return if the `condition` isn't met " +
				"(evaluates to `false`)."),
		},
	},
}
