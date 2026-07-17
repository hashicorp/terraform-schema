// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfpolicy "github.com/hashicorp/terraform-schema/policy"
	tfschema "github.com/hashicorp/terraform-schema/schema"
)

type SchemaMerger struct {
	coreSchema       *schema.BodySchema
	terraformVersion *version.Version
	stateReader      StateReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
}

func NewSchemaMerger(coreSchema *schema.BodySchema) *SchemaMerger {
	return &SchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *SchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *SchemaMerger) SetTerraformVersion(v *version.Version) {
	m.terraformVersion = v
}

func (m *SchemaMerger) SchemaForPolicy(meta *tfpolicy.Meta) (*schema.BodySchema, error) {
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

	if _, ok := mergedSchema.Blocks["input"]; ok {
		mergedSchema.Blocks["input"].Labels = []*schema.LabelSchema{
			{
				Name:        "name",
				IsDepKey:    true,
				Description: lang.PlainText("Input name"),
			},
		}
		mergedSchema.Blocks["input"].DependentBody = inputDependentBody(meta.Inputs)
	}

	return mergedSchema, nil
}

func inputDependentBody(inputs map[string]tfpolicy.Input) map[schema.SchemaKey]*schema.BodySchema {
	depBodies := make(map[schema.SchemaKey]*schema.BodySchema)

	for name, mInput := range inputs {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: name},
			},
		}
		depBodies[schema.NewSchemaKey(depKeys)] = &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"default": {
					Constraint:  schema.LiteralType{Type: mInput.Type},
					IsOptional:  true,
					Description: lang.Markdown("Default value to use when input is not explicitly set"),
				},
			},
		}
	}

	return depBodies
}
