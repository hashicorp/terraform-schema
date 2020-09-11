package refdecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/terraform-schema/internal/addrs"
)

func DecodeProviderReferences(m map[string]*hcl.File) (addrs.ProviderReferences, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	mod := tfconfig.NewModule("")
	for _, f := range m {
		fDiags := tfconfig.LoadModuleFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	refs := make(addrs.ProviderReferences, 0)

	for name, req := range mod.RequiredProviders {
		var src addrs.Provider

		if req.Source == "" {
			if name == "" {
				continue
			}
			src = addrs.ImpliedProviderForUnqualifiedType(name)
		} else {
			var err error
			src, err = addrs.ParseProviderSourceString(req.Source)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Unable to parse provider source for %q", name),
					Detail:   fmt.Sprintf("%q provider source (%q) is not a valid source string", name, req.Source),
				})
				continue
			}
		}

		refs[addrs.LocalProviderConfig{
			LocalName: name,
		}] = src
	}

	for name, aliases := range mod.ProviderAliases {
		src := refs[addrs.LocalProviderConfig{
			LocalName: name,
		}]
		for _, alias := range aliases {
			refs[addrs.LocalProviderConfig{
				LocalName: name,
				Alias:     alias,
			}] = src
		}
	}

	for _, resource := range mod.ManagedResources {
		providerName := resource.Provider.Name
		localRef := addrs.LocalProviderConfig{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			refs[localRef] = addrs.ImpliedProviderForUnqualifiedType(providerName)
		}
	}

	for _, dataSource := range mod.DataResources {
		providerName := dataSource.Provider.Name
		localRef := addrs.LocalProviderConfig{
			LocalName: providerName,
		}
		if _, exists := refs[localRef]; !exists && providerName != "" {
			refs[localRef] = addrs.ImpliedProviderForUnqualifiedType(providerName)
		}
	}

	return refs, diags
}
