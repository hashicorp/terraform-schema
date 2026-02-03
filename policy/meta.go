// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"github.com/hashicorp/go-version"
)

type Meta struct {
	Path      string
	Filenames []string

	CoreRequirements version.Constraints

	Variables        map[string]Variable
	ResourcePolicies map[string]ResourcePolicy
	ProviderPolicies map[string]ProviderPolicy
	ModulePolicies   map[string]ModulePolicy
}
