// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/policy"
)

func LoadPolicy(path string, files map[string]*hcl.File) (*policy.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedPolicy()
	for filename, f := range files {
		filenames = append(filenames, filename)
		fDiags := loadPolicyFromFile(f, mod)
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

	resourcePolicies := make(map[string]policy.ResourcePolicy)
	for key, rp := range mod.ResourcePolicies {
		resourcePolicies[key] = *rp
	}

	providerPolices := make(map[string]policy.ProviderPolicy)
	for key, pp := range mod.ProviderPolicies {
		providerPolices[key] = *pp
	}

	modulePolices := make(map[string]policy.ModulePolicy)
	for key, mp := range mod.ModulePolicies {
		modulePolices[key] = *mp
	}

	return &policy.Meta{
		Path:             path,
		Filenames:        filenames,
		CoreRequirements: coreRequirements,

		ResourcePolicies: resourcePolicies,
		ProviderPolicies: providerPolices,
		ModulePolicies:   modulePolices,
	}, diags
}
