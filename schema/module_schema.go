package schema

import (
	"sort"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func SchemaForDependentModuleBlock(localName string, modMeta *module.Meta) (*schema.BodySchema, error) {
	bodySchema := &schema.BodySchema{
		Attributes: variablesToAttrSchemas(modMeta.Variables),
	}

	if localName == "" {
		// avoid creating output refs if we don't have reference name
		return bodySchema, nil
	}

	modOutputTypes := make(map[string]cty.Type, 0)
	modOutputVals := make(map[string]cty.Value, 0)
	targetableOutputs := make(schema.Targetables, 0)

	for name, output := range modMeta.Outputs {
		addr := lang.Address{
			lang.RootStep{Name: "module"},
			lang.AttrStep{Name: localName},
			lang.AttrStep{Name: name},
		}

		typ := cty.DynamicPseudoType
		if !output.Value.IsNull() {
			typ = output.Value.Type()
		}

		targetable := &schema.Targetable{
			Address:           addr,
			ScopeId:           refscope.ModuleScope,
			AsType:            typ,
			IsSensitive:       output.IsSensitive,
			NestedTargetables: schema.NestedTargetablesForValue(addr, refscope.ModuleScope, output.Value),
		}
		if output.Description != "" {
			targetable.Description = lang.PlainText(output.Description)
		}

		targetableOutputs = append(targetableOutputs, targetable)

		modOutputTypes[name] = typ
		modOutputVals[name] = output.Value
	}

	sort.Sort(targetableOutputs)

	addr := lang.Address{
		lang.RootStep{Name: "module"},
		lang.AttrStep{Name: localName},
	}
	bodySchema.TargetableAs = append(bodySchema.TargetableAs, &schema.Targetable{
		Address:           addr,
		ScopeId:           refscope.ModuleScope,
		AsType:            cty.Object(modOutputTypes),
		NestedTargetables: targetableOutputs,
	})

	return bodySchema, nil
}
