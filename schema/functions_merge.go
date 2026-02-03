// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	tfmod "github.com/hashicorp/terraform-schema/module"
)

// FunctionsStateReader exposes a set of methods to read data from the internal language server state
// for function merging
type FunctionsStateReader interface {
	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*ProviderSchema, error)
}

type FunctionsMerger struct {
	coreFunctions    map[string]schema.FunctionSignature
	terraformVersion *version.Version
	stateReader      FunctionsStateReader
}

func NewFunctionsMerger(coreFunctions map[string]schema.FunctionSignature) *FunctionsMerger {
	return &FunctionsMerger{
		coreFunctions: coreFunctions,
	}
}

func (m *FunctionsMerger) SetStateReader(mr FunctionsStateReader) {
	m.stateReader = mr
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

	if m.stateReader == nil {
		return m.coreFunctions, nil
	}

	if m.terraformVersion.LessThan(v1_8) {
		return m.coreFunctions, nil
	}

	mergedFunctions := make(map[string]schema.FunctionSignature, len(m.coreFunctions))
	for fName, fSig := range m.coreFunctions {
		mergedFunctions[fName] = *fSig.Copy()
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	for pAddr, pVersionCons := range meta.ProviderRequirements {
		pSchema, err := m.stateReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
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

	return mergedFunctions, nil
}
