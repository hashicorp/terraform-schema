package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func SchemaForVariables(vars map[string]module.Variable) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema)

	for name, modVar := range vars {
		aSchema := moduleVarToAttribute(modVar)
		aSchema.Expr = schema.LiteralTypeOnly(typeOfModuleVar(modVar))
		attributes[name] = aSchema
	}

	return &schema.BodySchema{
		Attributes: attributes,
	}, nil
}

type exprFunc func(cty.Type) schema.ExprConstraints

func moduleVarToAttribute(modVar module.Variable) *schema.AttributeSchema {
	aSchema := &schema.AttributeSchema{
		IsSensitive: modVar.IsSensitive,
	}

	if modVar.Description != "" {
		aSchema.Description = lang.PlainText(modVar.Description)
	}

	if modVar.DefaultValue == cty.NilVal {
		aSchema.IsRequired = true
	} else {
		aSchema.IsOptional = true
	}

	return aSchema
}

func typeOfModuleVar(modVar module.Variable) cty.Type {
	if (modVar.Type == cty.DynamicPseudoType || modVar.Type == cty.NilType) &&
		modVar.DefaultValue != cty.NilVal {
		// infer type from default value if one is not specified
		// or when it's "any"
		return modVar.DefaultValue.Type()
	}

	return modVar.Type
}
