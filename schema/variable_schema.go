// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func SchemaForVariables(vars map[string]module.Variable, modPath string) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema)

	for name, modVar := range vars {
		aSchema := moduleVarToAttribute(modVar)
		varType := typeOfModuleVar(modVar)
		aSchema.Constraint = schema.LiteralType{Type: varType}
		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       modPath,
				LanguageID: ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

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
