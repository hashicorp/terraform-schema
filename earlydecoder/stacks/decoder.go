// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/stack"
)

func LoadStack(path string, files map[string]*hcl.File) (*stack.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedStack()
	for filename, f := range files {
		filenames = append(filenames, filename)
		// TODO need more stack metas
		fDiags := loadStackFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	sort.Strings(filenames)

	components := make(map[string]stack.Component)
	for key, variable := range mod.Components {
		components[key] = *variable
	}

	variables := make(map[string]stack.Variable)
	for key, variable := range mod.Variables {
		variables[key] = *variable
	}

	outputs := make(map[string]stack.Output)
	for key, output := range mod.Outputs {
		outputs[key] = *output
	}

	providerRequirements := make(map[string]stack.ProviderRequirement)
	for key, providerRequirement := range mod.ProviderRequirements {
		providerRequirements[key] = stack.ProviderRequirement{
			Source:             providerRequirement.Source,
			VersionConstraints: providerRequirement.VersionConstraints,
		}
	}

	return &stack.Meta{
		Path:                 path,
		Filenames:            filenames,
		Components:           components,
		Variables:            variables,
		Outputs:              outputs,
		ProviderRequirements: providerRequirements,
	}, diags
}
