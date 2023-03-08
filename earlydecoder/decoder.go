// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/addr"
	"github.com/hashicorp/terraform-schema/module"
)

func LoadModule(path string, files map[string]*hcl.File) (*module.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedModule()
	for filename, f := range files {
		filenames = append(filenames, filename)
		fDiags := loadModuleFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	sort.Strings(filenames)

	var coreRequirements version.Constraints
	for _, rc := range mod.RequiredCore {
		c, err := version.NewConstraint(rc)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unable to parse terraform requirements",
				Detail:   fmt.Sprintf("Constraint %q is not a valid constraint: %s", rc, err),
			})
			continue
		}
		coreRequirements = append(coreRequirements, c...)
	}

	var backend *module.Backend
	if len(mod.Backends) == 1 {
		for bType, data := range mod.Backends {
			backend = &module.Backend{
				Type: bType,
				Data: data,
			}
		}
	} else if len(mod.Backends) > 1 {
		backendTypes := make([]string, len(mod.Backends))
		for bType := range mod.Backends {
			backendTypes = append(backendTypes, bType)
		}

		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unable to parse backend configuration",
			Detail:   fmt.Sprintf("Multiple backend definitions: %q", backendTypes),
		})
	}

	var (
		providerRequirements = make(map[tfaddr.Provider]version.Constraints, 0)
		refs                 = make(map[module.ProviderRef]tfaddr.Provider, 0)
	)

	for name, req := range mod.ProviderRequirements {
		var src tfaddr.Provider

		if req.Source == "" {
			if name == "" {
				continue
			}

			_, err := tfaddr.ParseProviderPart(name)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid provider name",
					Detail:   fmt.Sprintf("%q is not a valid provider name: %s", name, err),
				})
				continue
			}

			src = addr.NewLegacyProvider(name)
		} else {
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
		}

		constraints := make(version.Constraints, 0)
		for _, vc := range req.VersionConstraints {
			c, err := version.NewConstraint(vc)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Unable to parse %q provider requirements", name),
					Detail:   fmt.Sprintf("Constraint %q is not a valid constraint: %s", vc, err),
				})
				continue
			}
			constraints = append(constraints, c...)
		}

		providerRequirements[src] = constraints

		refs[module.ProviderRef{
			LocalName: name,
		}] = src

		for _, alias := range req.ConfigurationAliases {
			refs[module.ProviderRef{
				LocalName: alias.LocalName,
				Alias:     alias.Alias,
			}] = src
		}
	}

	for _, cfg := range mod.ProviderConfigs {
		src := refs[module.ProviderRef{
			LocalName: cfg.Name,
		}]
		if cfg.Alias != "" {
			refs[module.ProviderRef{
				LocalName: cfg.Name,
				Alias:     cfg.Alias,
			}] = src
		}
	}

	for _, resource := range mod.Resources {
		providerName := resource.Provider.LocalName

		_, err := tfaddr.ParseProviderPart(providerName)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider name",
				Detail:   fmt.Sprintf("%q is not a valid implied provider name: %s", providerName, err),
			})
			continue
		}

		localRef := module.ProviderRef{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			src := addr.NewLegacyProvider(providerName)
			if _, exists := providerRequirements[src]; !exists {
				providerRequirements[src] = version.Constraints{}
			}

			refs[localRef] = src
		}
	}

	for _, dataSource := range mod.DataSources {
		providerName := dataSource.Provider.LocalName

		_, err := tfaddr.ParseProviderPart(providerName)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider name",
				Detail:   fmt.Sprintf("%q is not a valid implied provider name: %s", providerName, err),
			})
			continue
		}

		localRef := module.ProviderRef{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			src := addr.NewLegacyProvider(providerName)
			if _, exists := providerRequirements[src]; !exists {
				providerRequirements[src] = version.Constraints{}
			}
			refs[localRef] = src
		}
	}

	variables := make(map[string]module.Variable)
	for key, variable := range mod.Variables {
		variables[key] = *variable
	}

	outputs := make(map[string]module.Output)
	for key, output := range mod.Outputs {
		outputs[key] = *output
	}

	modulesCalls := make(map[string]module.DeclaredModuleCall)
	for key, moduleCall := range mod.ModuleCalls {
		modulesCalls[key] = *moduleCall
	}

	return &module.Meta{
		Path:                 path,
		Backend:              backend,
		Cloud:                mod.CloudBackend,
		ProviderReferences:   refs,
		ProviderRequirements: providerRequirements,
		CoreRequirements:     coreRequirements,
		Variables:            variables,
		Outputs:              outputs,
		Filenames:            filenames,
		ModuleCalls:          modulesCalls,
	}, diags
}
