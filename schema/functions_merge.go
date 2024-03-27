// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfmod "github.com/hashicorp/terraform-schema/module"
)

type FunctionsMerger struct {
	coreFunctions    map[string]schema.FunctionSignature
	terraformVersion *version.Version
	schemaReader     SchemaReader
}

func NewFunctionsMerger(coreFunctions map[string]schema.FunctionSignature) *FunctionsMerger {
	return &FunctionsMerger{
		coreFunctions: coreFunctions,
	}
}

func (m *FunctionsMerger) SetSchemaReader(sr SchemaReader) {
	m.schemaReader = sr
}

func (m *FunctionsMerger) SetTerraformVersion(v *version.Version) {
	m.terraformVersion = v
}

func (m *FunctionsMerger) FunctionsForModule(meta *tfmod.Meta) (map[string]schema.FunctionSignature, error) {
	if m.coreFunctions == nil {
		return nil, coreFunctionsRequiredErr{}
	}

	if meta == nil {
		return m.coreFunctions, nil
	}

	mergedFunctions := make(map[string]schema.FunctionSignature, len(m.coreFunctions))
	for fName, fSig := range m.coreFunctions {
		mergedFunctions[fName] = *fSig.Copy()
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	if m.schemaReader != nil {
		for pAddr, pVersionCons := range meta.ProviderRequirements {
			pSchema, err := m.schemaReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
			if err != nil {
				continue
			}

			refs := providerRefs.ReferencesOfProvider(pAddr)

			for _, localRef := range refs {
				for fName, fSig := range pSchema.Functions {
					mergedFunctions[fmt.Sprintf("provider::%s::%s", localRef.LocalName, fName)] = *fSig.Copy()
				}
			}
		}
	}

	return mergedFunctions, nil
}
