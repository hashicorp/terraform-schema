package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func SchemaForVariables(vars map[string]module.Variable) (*schema.BodySchema, error) {
	varSchemas := make(map[string]*schema.AttributeSchema)

	for name, v := range vars {
		varType := v.Type
		if (varType == cty.DynamicPseudoType || varType == cty.NilType) &&
			v.DefaultValue != cty.NilVal {
			// infer type from default value if one is not specified
			// or when it's "any"
			varType = v.DefaultValue.Type()
		}

		varSchemas[name] = &schema.AttributeSchema{
			Description: lang.MarkupContent{
				Value: v.Description,
				Kind:  lang.PlainTextKind,
			},
			Expr:        schema.ExprConstraints{schema.LiteralTypeExpr{Type: varType}},
			IsSensitive: v.IsSensitive,
		}

		if v.DefaultValue == cty.NilVal {
			varSchemas[name].IsRequired = true
		} else {
			varSchemas[name].IsOptional = true
		}
	}

	return &schema.BodySchema{
		Attributes: varSchemas,
	}, nil
}
