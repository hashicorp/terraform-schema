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

	constr, err := constraintForDeploymentInputs(*meta)
	if err != nil {
		return mergedSchema, err
	}

	// TODO: isOptional should be set to true if at least one input is required
	mergedSchema.Blocks["deployment"].Body.Attributes["inputs"].Constraint = constr

	// We add specific Targetables with a dynamic type for each store as the AnyAttribute defined in the store schema
	// is not picked up when reference targets are collected because we don't know all the available attributes that
	// are eventually available in the varset and reference targets need a specific name to match against
	// However, it is also possible to target the parent if its type is dynamic. This is why we add this target.
	// We can't set the body to dynamic in the schema as that would remove the completions for the "id" attribute of the
	// varset store block.
	// If there's a better way to do this, we should consider it as this feels like a bit of a hack.
	// TODO: once we parse tfvars files, do this differently for the tfvars type store block
	for name, store := range meta.Stores {
		key := schema.NewSchemaKey(schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: store.Type},
				{Index: 1, Value: name},
			}})

		// Copy the body of the store block for the specific store type to keep attributes and constraints
		// as dependent bodies replace the original ones
		newBody := mergedSchema.Blocks["store"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: store.Type},
			}})].Copy()

		if newBody == nil {
			newBody = mergedSchema.Blocks["store"].Body
		}

		newBody.TargetableAs = schema.Targetables{
			&schema.Targetable{
				Address: lang.Address{
					lang.RootStep{Name: "store"},
					lang.AttrStep{Name: store.Type},
					lang.AttrStep{Name: name},
				},
				AsType:       cty.DynamicPseudoType, // we need this type for the target for this store block as it matches every nested value
				ScopeId:      refscope.StoreScope,
				FriendlyName: "store",
			},
		}

		mergedSchema.Blocks["store"].DependentBody[key] = newBody
	}

	return mergedSchema, nil
}

func constraintForDeploymentInputs(stackMeta stack.Meta) (schema.Constraint, error) {
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
	}, nil
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
