// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	tfmod "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty/cty"
)

type StackSchemaMerger struct {
	coreSchema  *schema.BodySchema
	stateReader StateReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error)

	// InstalledModulePath checks if there is an installed module available for
	// the given normalized source address.
	InstalledModulePath(rootPath string, normalizedSource string) (string, bool)

	// LocalModuleMeta returns the module meta data for a local module. This is the result
	// of the [earlydecoder] when processing module files
	LocalModuleMeta(modPath string) (*tfmod.Meta, error)
}

func NewStackSchemaMerger(coreSchema *schema.BodySchema) *StackSchemaMerger {
	return &StackSchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *StackSchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *StackSchemaMerger) SchemaForStack(meta *stack.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, tfschema.CoreSchemaRequiredErr{}
	}

	if meta == nil {
		return m.coreSchema, nil
	}

	if m.stateReader == nil {
		return m.coreSchema, nil
	}

	mergedSchema := m.coreSchema.Copy()

	if mergedSchema.Blocks["provider"].DependentBody == nil {
		mergedSchema.Blocks["provider"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["component"].DependentBody == nil {
		mergedSchema.Blocks["component"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}

	for localName, pReq := range meta.ProviderRequirements {
		pSchema, err := m.stateReader.ProviderSchema(meta.Path, pReq.Source, pReq.VersionConstraints)
		if err != nil {
			continue
		}

		if pSchema.Provider != nil {
			mergedSchema.Blocks["provider"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: localName},
				},
			})] = &schema.BodySchema{
				HoverURL:     pSchema.Provider.HoverURL,
				DocsLink:     pSchema.Provider.DocsLink,
				Detail:       pSchema.Provider.Detail,
				Description:  pSchema.Provider.Description,
				IsDeprecated: pSchema.Provider.IsDeprecated,
				Blocks: map[string]*schema.BlockSchema{
					"config": {
						Body: pSchema.Provider,
					},
				},
			}
		}
	}

	if _, ok := mergedSchema.Blocks["variable"]; ok {
		mergedSchema.Blocks["variable"].Labels = []*schema.LabelSchema{
			{
				Name:        "name",
				IsDepKey:    true,
				Description: lang.PlainText("Variable name"),
			},
		}
		mergedSchema.Blocks["variable"].DependentBody = variableDependentBody(meta.Variables)
	}

	for name, comp := range meta.Components {
		depKeys := schema.DependencyKeys{
			Attributes: []schema.AttributeDependent{
				{
					Name: "source",
					Expr: schema.ExpressionValue{
						Static: cty.StringVal(comp.Source),
					},
				},
			},
		}

		switch sourceAddr := comp.SourceAddr.(type) {
		case tfmod.LocalSourceAddr:
			path := filepath.Join(meta.Path, sourceAddr.String())

			modMeta, err := m.stateReader.LocalModuleMeta(path)
			if err == nil {
				depSchema, err := schemaForDependentComponentBlock(modMeta, comp, name)
				if err == nil {
					mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
				}
			}
		case tfaddr.Module:
			// make sure there always is a dependent body schema, even if we've errors and can't get the module schema
			mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)] = &schema.BodySchema{}

			installedDir, ok := m.stateReader.InstalledModulePath(meta.Path, sourceAddr.String())
			if ok {
				path := filepath.Join(meta.Path, installedDir)

				// TODO: how to ensure this dir is parsed and available?
				modMeta, err := m.stateReader.LocalModuleMeta(path)

				if err == nil {
					depSchema, err := schemaForDependentComponentBlock(modMeta, comp, name)
					if err == nil {
						mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
					}
				}
			}

			// components with a source pointing to a registry module require a version constraint
			if mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)].Attributes == nil {
				mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)].Attributes = make(map[string]*schema.AttributeSchema)
			}

			mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)].Attributes["version"] = &schema.AttributeSchema{
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Accepts a comma-separated list of version constraints for registry modules. Required for registry modules"),
				IsRequired:  true,
			}

			// TODO: support API based schema for registry modules (would require GetModuleDataFromRegistry() job in stacks feature as well)

		case tfmod.RemoteSourceAddr:
			installedDir, ok := m.stateReader.InstalledModulePath(meta.Path, sourceAddr.String())
			if !ok {
				continue
			}
			path := filepath.Join(meta.Path, installedDir)

			// TODO: how to ensure this dir is parsed and available?
			modMeta, err := m.stateReader.LocalModuleMeta(path)
			if err == nil {
				depSchema, err := schemaForDependentComponentBlock(modMeta, comp, name)
				if err == nil {
					mergedSchema.Blocks["component"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
				}
			}
		}
	}

	return mergedSchema, nil
}

func variableDependentBody(vars map[string]stack.Variable) map[schema.SchemaKey]*schema.BodySchema {
	depBodies := make(map[schema.SchemaKey]*schema.BodySchema)

	for name, mVar := range vars {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: name},
			},
		}
		depBodies[schema.NewSchemaKey(depKeys)] = &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"default": {
					Constraint:  schema.LiteralType{Type: mVar.Type},
					IsOptional:  true,
					Description: lang.Markdown("Default value to use when variable is not explicitly set"),
				},
			},
		}
	}

	return depBodies
}

func schemaForDependentComponentBlock(modMeta *tfmod.Meta, component stack.Component, componentName string) (*schema.BodySchema, error) {
	inputs := make(map[string]*schema.AttributeSchema, 0)
	providers := make(map[string]*schema.AttributeSchema, 0)

	for name, modVar := range modMeta.Variables {
		varType := modVar.Type
		if varType == cty.NilType {
			varType = cty.DynamicPseudoType
		}
		aSchema := tfschema.ModuleVarToAttribute(modVar)
		aSchema.Constraint = tfschema.ConvertAttributeTypeToConstraint(varType)
		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: tfschema.ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

		inputs[name] = aSchema
	}

	for pRef := range modMeta.ProviderReferences {
		addr := pRef.LocalName
		if pRef.Alias != "" {
			addr += "." + pRef.Alias
		}

		providers[addr] = &schema.AttributeSchema{
			Constraint: schema.Reference{OfScopeId: refscope.ProviderScope},
		}
	}

	bodySchema := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"inputs": {
				Constraint: schema.Object{
					Attributes: inputs,
				},
			},
			"providers": {
				Constraint: schema.Object{
					Attributes: providers,
				},
			},
		},
	}

	if component.Source == "" {
		// avoid creating output refs if we don't have reference name
		return bodySchema, nil
	}

	modOutputTypes := make(map[string]cty.Type, 0)
	modOutputVals := make(map[string]cty.Value, 0)
	targetableOutputs := make(schema.Targetables, 0)
	impliedOrigins := make(schema.ImpliedOrigins, 0)

	for name, output := range modMeta.Outputs {
		addr := lang.Address{
			lang.RootStep{Name: "component"},
			lang.AttrStep{Name: componentName},
			lang.AttrStep{Name: name},
		}

		typ := cty.DynamicPseudoType
		if !output.Value.IsNull() {
			typ = output.Value.Type()
		}

		targetable := &schema.Targetable{
			Address:           addr,
			ScopeId:           refscope.ComponentScope,
			AsType:            typ,
			IsSensitive:       output.IsSensitive,
			NestedTargetables: schema.NestedTargetablesForValue(addr, refscope.ComponentScope, output.Value),
		}
		if output.Description != "" {
			targetable.Description = lang.PlainText(output.Description)
		}

		targetableOutputs = append(targetableOutputs, targetable)

		modOutputTypes[name] = typ
		modOutputVals[name] = output.Value

		impliedOrigins = append(impliedOrigins, schema.ImpliedOrigin{
			OriginAddress: lang.Address{
				lang.RootStep{Name: "component"},
				lang.AttrStep{Name: componentName},
				lang.AttrStep{Name: name},
			},
			TargetAddress: lang.Address{
				lang.RootStep{Name: "output"},
				lang.AttrStep{Name: name},
			},
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: tfschema.ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.OutputScope,
			},
		})
	}

	bodySchema.ImpliedOrigins = impliedOrigins

	sort.Sort(targetableOutputs)

	addr := lang.Address{
		lang.RootStep{Name: "component"},
		lang.AttrStep{Name: componentName},
	}
	bodySchema.TargetableAs = append(bodySchema.TargetableAs, &schema.Targetable{
		Address:           addr,
		ScopeId:           refscope.ModuleScope,
		AsType:            cty.Object(modOutputTypes),
		NestedTargetables: targetableOutputs,
	})

	if len(modMeta.Filenames) > 0 {
		filename := modMeta.Filenames[0]

		// Prioritize main.tf based on best practices as documented at
		// https://learn.hashicorp.com/tutorials/terraform/module-create
		if sliceContains(modMeta.Filenames, "main.tf") {
			filename = "main.tf"
		}

		bodySchema.Targets = &schema.Target{
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: tfschema.ModuleLanguageID,
			},
			Range: hcl.Range{
				Filename: filename,
				Start:    hcl.InitialPos,
				End:      hcl.InitialPos,
			},
		}
	}

	return bodySchema, nil
}

func sliceContains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}
