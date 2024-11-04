// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/stack"
)

func LoadStack(path string, files map[string]*hcl.File) (*stack.Meta, map[string]hcl.Diagnostics) {
	filenames := make([]string, 0)
	sdiags := make(map[string]hcl.Diagnostics, 0)

	mod := newDecodedStack()
	for filename, f := range files {
		filenames = append(filenames, filename)

		if isStackFilename(filename) {
			sdiags[filename] = loadStackFromFile(f, mod)
		}
	}

	sort.Strings(filenames)

	components := make(map[string]stack.Component)
	for key, component := range mod.Components {
		components[key] = *component
	}

	variables := make(map[string]stack.Variable)
	for key, variable := range mod.Variables {
		variables[key] = *variable
	}

	outputs := make(map[string]stack.Output)
	for key, output := range mod.Outputs {
		outputs[key] = *output
	}

	providerRequirements := make(map[string]stack.ProviderRequirement, len(mod.ProviderRequirements))
	for name, req := range mod.ProviderRequirements {
		providerRequirements[name] = stack.ProviderRequirement{
			Source:             *req.Source,
			VersionConstraints: *req.VersionConstraints,
		}
	}

	return &stack.Meta{
		Path:                 path,
		Filenames:            filenames,
		Components:           components,
		Variables:            variables,
		Outputs:              outputs,
		ProviderRequirements: providerRequirements,
	}, sdiags
}

func isStackFilename(name string) bool {
	return strings.HasSuffix(name, ".tfstack.hcl") ||
		strings.HasSuffix(name, ".tfstack.json")
}
