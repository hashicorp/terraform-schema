// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

// PolicyTestSchema returns the static schema for a policytest
// configuration (*.policytest.hcl) file.
func PolicyTestSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"data":       dataBlockSchema(),
			"locals":     localsBlockSchema(),
			"module":     moduleBlockSchema(),
			"policytest": policytestBlockSchema(),
			"provider":   providerBlockSchema(),
			"resource":   resourceBlockSchema(),
		},
	}
}
