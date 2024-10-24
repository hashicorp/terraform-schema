// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
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

	var providerRequirements = make(map[string]stack.ProviderRequirement, 0)
	for name, req := range mod.ProviderRequirements {
		var src tfaddr.Provider

		var err error
		src, err = tfaddr.ParseProviderSource(req.Source)
		if err != nil {
			// This is handled in decodeRequiredProvidersBlock now
			continue
		}

		constraints, err := version.NewConstraint(req.VersionConstraints)
		if err != nil {
			// This is handled in decodeRequiredProvidersBlock now
			continue
		}

		providerRequirements[name] = stack.ProviderRequirement{
			Source:             src,
			VersionConstraints: constraints,
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
