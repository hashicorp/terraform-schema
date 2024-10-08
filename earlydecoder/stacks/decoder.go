// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/stack"
)

func LoadStack(path string, files map[string]*hcl.File) (*stack.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedStack()
	for filename, f := range files {
		filenames = append(filenames, filename)

		if isStackFilename(filename) {
			fDiags := loadStackFromFile(f, mod)
			diags = append(diags, fDiags...)
		} else if isDeployFilename(filename) {
			fDiags := loadDeployFromFile(f, mod)
			diags = append(diags, fDiags...)
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
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Unable to parse provider source for %q", name),
				Detail:   fmt.Sprintf("%q provider source (%q) is not a valid source string", name, req.Source),
			})
			continue
		}

		constraints, err := version.NewConstraint(req.VersionConstraints)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Unable to parse %q provider requirements", name),
				Detail:   fmt.Sprintf("Constraint %q is not a valid constraint: %s", req.VersionConstraints, err),
			})
			continue
		}

		providerRequirements[name] = stack.ProviderRequirement{
			Source:             src,
			VersionConstraints: constraints,
		}
	}

	deployments := make(map[string]stack.Deployment)
	for key, deployment := range mod.Deployments {
		deployments[key] = *deployment
	}

	stores := make(map[string]stack.Store)
	for key, store := range mod.Stores {
		stores[key] = *store
	}

	orchestrationRules := make(map[string]stack.OrchestrationRule)
	for key, rule := range mod.OrchestrationRules {
		orchestrationRules[key] = *rule
	}

	return &stack.Meta{
		Path:                 path,
		Filenames:            filenames,
		Components:           components,
		Variables:            variables,
		Outputs:              outputs,
		ProviderRequirements: providerRequirements,
		Deployments:          deployments,
		Stores:               stores,
		OrchestrationRules:   orchestrationRules,
	}, diags
}

func isStackFilename(name string) bool {
	return strings.HasSuffix(name, ".tfstack.hcl") ||
		strings.HasSuffix(name, ".tfstack.json")
}

func isDeployFilename(name string) bool {
	return strings.HasSuffix(name, ".tfdeploy.hcl") ||
		strings.HasSuffix(name, ".tfdeploy.json")
}
