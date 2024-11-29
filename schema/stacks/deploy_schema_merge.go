// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty/cty"
)

type DeploySchemaMerger struct {
	coreSchema *schema.BodySchema
}

func NewDeploySchemaMerger(coreSchema *schema.BodySchema) *DeploySchemaMerger {
	return &DeploySchemaMerger{
		coreSchema: coreSchema,
	}
}

// SchemaForDeployment returns the schema for a deployment block
func (m *DeploySchemaMerger) SchemaForDeployment(meta *stack.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, tfschema.CoreSchemaRequiredErr{}
	}

	mergedSchema := m.coreSchema.Copy()

	constr := constraintForDeploymentInputs(*meta)

	// TODO: isOptional should be set to true if at least one input is required
	mergedSchema.Blocks["deployment"].Body.Attributes["inputs"].Constraint = constr

	return mergedSchema, nil
}

func constraintForDeploymentInputs(stackMeta stack.Meta) schema.Constraint {
	inputs := make(map[string]*schema.AttributeSchema, 0)

	for name, variable := range stackMeta.Variables {
		varType := variable.Type
		if varType == cty.NilType {
			varType = cty.DynamicPseudoType
		}
		aSchema := StackVarToAttribute(variable)
		aSchema.Constraint = tfschema.ConvertAttributeTypeToConstraint(varType)

		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       stackMeta.Path,
				LanguageID: tfschema.StackLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

		inputs[name] = aSchema
	}

	return schema.Object{
		Attributes: inputs,
	}
}

func StackVarToAttribute(stackVar stack.Variable) *schema.AttributeSchema {
	aSchema := &schema.AttributeSchema{
		IsSensitive: stackVar.IsSensitive,
	}

	if stackVar.Description != "" {
		aSchema.Description = lang.PlainText(stackVar.Description)
	}

	if stackVar.DefaultValue == cty.NilVal {
		aSchema.IsRequired = true
	} else {
		aSchema.IsOptional = true
	}

	return aSchema
}
