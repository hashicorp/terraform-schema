// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
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
}

func NewStackSchemaMerger(coreSchema *schema.BodySchema) *StackSchemaMerger {
	return &StackSchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *StackSchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
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

	return mergedSchema, nil
}
