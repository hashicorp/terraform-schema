// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

// SearchSchema returns the static schema for a search
// configuration (*.tfsearch.hcl) file.
func SearchSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"list":     listBlockSchema(),
			"locals":   localsBlockSchema(),
			"provider": providerBlockSchema(),
			"variable": variableBlockSchema(),
		},
	}
}
