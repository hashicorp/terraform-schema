// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

// PolicySchema returns the static schema for a policy
// configuration (*.policy.hcl) file.
func PolicySchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"locals":          localsBlockSchema(),
			"module_policy":   modulePolicyBlockSchema(),
			"policy":          policyBlockSchema(),
			"provider_policy": providerPolicyBlockSchema(),
			"resource_policy": resourcePolicyBlockSchema(),
			// "variable":        variableBlockSchema(),
		},
	}
}
