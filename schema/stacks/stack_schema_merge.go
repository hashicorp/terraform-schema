// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	tfmod "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty/cty"
)

type StackSchemaMerger struct {
	coreSchema   *schema.BodySchema
	stateReader  StateReader
	moduleReader ModuleReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error)
}

type ModuleReader interface {
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

func (m *StackSchemaMerger) SetModuleReader(mr ModuleReader) {
	m.moduleReader = mr
}

func (m *StackSchemaMerger) SchemaForModule(meta *stack.Meta) (*schema.BodySchema, error) {
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

	for _, comp := range meta.Components {
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

			modMeta, err := m.moduleReader.LocalModuleMeta(path)
			if err == nil {
				depSchema, err := schemaForDependentComponentBlock(comp, modMeta)
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

// TODO: move to extra file? (there's an extra file in the original code in modules)
func schemaForDependentComponentBlock(component stack.Component, modMeta *tfmod.Meta) (*schema.BodySchema, error) {
	// TODO: implement

	return nil, nil
}
