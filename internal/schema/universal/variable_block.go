package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

var variableBlockSchema = &schema.BlockSchema{
	Address: &schema.BlockAddrSchema{
		Steps: []schema.AddrStep{
			schema.StaticStep{Name: "var"},
			schema.LabelStep{Index: 0},
		},
		FriendlyName: "variable",
		ScopeId:      refscope.VariableScope,
		AsReference:  true,
		AsTypeOf: &schema.BlockAsTypeOf{
			AttributeExpr:  "type",
			AttributeValue: "default",
		},
	},
	Labels: []*schema.LabelSchema{
		{
			Name:        "name",
			Description: lang.PlainText("Variable Name"),
		},
	},
	Description: lang.Markdown("Input variable allowing users to customizate aspects of the configuration when used directly " +
		"(e.g. via CLI, `tfvars` file or via environment variables), or as a module (via `module` arguments)"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"description": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Description to document the purpose of the variable and what value is expected"),
			},
			"type": {
				Expr:        schema.ExprConstraints{schema.TypeDeclarationExpr{}},
				IsOptional:  true,
				Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. `string` or `list(string)`"),
			},
		},
		Blocks: make(map[string]*schema.BlockSchema, 0),
	},
}
