// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	tfmod "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	tfsearch "github.com/hashicorp/terraform-schema/search"
)

type SearchSchemaMerger struct {
	coreSchema  *schema.BodySchema
	stateReader StateReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error)

	// LocalModuleMeta returns the module meta data for a local module. This is the result
	// of the [earlydecoder] when processing module files
	LocalModuleMeta(modPath string) (*tfmod.Meta, error)
}

func NewSearchSchemaMerger(coreSchema *schema.BodySchema) *SearchSchemaMerger {
	return &SearchSchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *SearchSchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *SearchSchemaMerger) SchemaForSearch(meta *tfsearch.Meta) (*schema.BodySchema, error) {
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

	if mergedSchema.Blocks["list"].DependentBody == nil {
		mergedSchema.Blocks["list"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
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

	providerRefs := ProviderReferences(meta.ProviderReferences)

	modMeta, err := m.stateReader.LocalModuleMeta(meta.Path)
	if err == nil && modMeta != nil && modMeta.ProviderRequirements != nil {

		for pAddr, pVersionCons := range modMeta.ProviderRequirements {
			pSchema, err := m.stateReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
			if err != nil {
				continue
			}
			refs := providerRefs.ReferencesOfProvider(pAddr)

			for _, localRef := range refs {
				if pSchema.Provider != nil {
					mergedSchema.Blocks["provider"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: localRef.LocalName},
						},
					})] = pSchema.Provider
				}

				providerAddr := lang.Address{
					lang.RootStep{Name: localRef.LocalName},
				}
				if localRef.Alias != "" {
					providerAddr = append(providerAddr, lang.AttrStep{Name: localRef.Alias})
				}
				for lrName, lrSchema := range pSchema.ListResources {
					// Create a BodySchema that ensures a config block exists
					listBodySchema := &schema.BodySchema{
						HoverURL:     lrSchema.HoverURL,
						DocsLink:     lrSchema.DocsLink,
						Detail:       lrSchema.Detail,
						Description:  lrSchema.Description,
						IsDeprecated: lrSchema.IsDeprecated,
					}

					listBodySchema.Blocks = map[string]*schema.BlockSchema{
						"config": {
							Description: lang.Markdown("Filters specific to the list type"),
							MaxItems:    1,
							Body:        lrSchema,
						},
					}
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: lrName},
						},
						Attributes: []schema.AttributeDependent{
							{
								Name: "provider",
								Expr: schema.ExpressionValue{
									Address: providerAddr,
								},
							},
						},
					}
					mergedSchema.Blocks["list"].DependentBody[schema.NewSchemaKey(depKeys)] = listBodySchema

				}
			}
		}
	}

	return mergedSchema, nil
}

func variableDependentBody(vars map[string]tfsearch.Variable) map[schema.SchemaKey]*schema.BodySchema {
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

type ProviderReferences map[tfsearch.ProviderRef]tfaddr.Provider

func (pr ProviderReferences) ReferencesOfProvider(addr tfaddr.Provider) []tfsearch.ProviderRef {
	refs := make([]tfsearch.ProviderRef, 0)

	for ref, pAddr := range pr {
		if pAddr.Equals(addr) {
			refs = append(refs, ref)
		}
	}
	return refs
}

func TypeBelongsToProvider(typeName string, pRef tfsearch.ProviderRef) bool {
	return typeName == pRef.LocalName || strings.HasPrefix(typeName, pRef.LocalName+"_")
}
