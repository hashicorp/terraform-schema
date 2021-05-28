package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/module"
)

func LoadModule(path string, files map[string]*hcl.File) (*module.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	mod := newDecodedModule()
	for _, f := range files {
		fDiags := loadModuleFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

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
			src = tfaddr.NewLegacyProvider(name)
		} else {
			var err error
			src, err = tfaddr.ParseRawProviderSourceString(req.Source)
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
		localRef := module.ProviderRef{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			src := tfaddr.NewLegacyProvider(providerName)
			if _, exists := providerRequirements[src]; !exists {
				providerRequirements[src] = version.Constraints{}
			}

			refs[localRef] = src
		}
	}

	for _, dataSource := range mod.DataSources {
		providerName := dataSource.Provider.LocalName
		localRef := module.ProviderRef{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			src := tfaddr.NewLegacyProvider(providerName)
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
	return &module.Meta{
		Path:                 path,
		ProviderReferences:   refs,
		ProviderRequirements: providerRequirements,
		CoreRequirements:     coreRequirements,
		Variables:            variables,
	}, diags
}
