// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/schema"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	tftest "github.com/hashicorp/terraform-schema/test"
)

type MockSchemaMerger struct {
	coreSchema  *schema.BodySchema
	stateReader StateReader
}

func NewMockSchemaMerger(coreSchema *schema.BodySchema) *MockSchemaMerger {
	return &MockSchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *MockSchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *MockSchemaMerger) SchemaForMock(meta *tftest.Meta) (*schema.BodySchema, error) {
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

	// TODO merge mock_resource blocks - use the label as dependency key TFECO-7471
	// TODO merge mock_data blocks - use the label as dependency key TFECO-7472

	return mergedSchema, nil
}
