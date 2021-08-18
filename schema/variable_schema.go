package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func SchemaForVariables(vars map[string]module.Variable) (*schema.BodySchema, error) {
	return &schema.BodySchema{
		Attributes: variablesToAttrSchemas(vars, schema.LiteralTypeOnly),
	}, nil
}

type exprFunc func(cty.Type) schema.ExprConstraints

func variablesToAttrSchemas(vars map[string]module.Variable, exprFunc exprFunc) map[string]*schema.AttributeSchema {
	varSchemas := make(map[string]*schema.AttributeSchema)

	for name, v := range vars {
		varType := v.Type
		if (varType == cty.DynamicPseudoType || varType == cty.NilType) &&
			v.DefaultValue != cty.NilVal {
			// infer type from default value if one is not specified
			// or when it's "any"
			varType = v.DefaultValue.Type()
		}

		aSchema := &schema.AttributeSchema{
			Expr:        exprFunc(varType),
			IsSensitive: v.IsSensitive,
		}
		if v.Description != "" {
			aSchema.Description = lang.PlainText(v.Description)
		}

		varSchemas[name] = aSchema

		if v.DefaultValue == cty.NilVal {
			varSchemas[name].IsRequired = true
		} else {
			varSchemas[name].IsOptional = true
		}
	}

	return varSchemas
}
