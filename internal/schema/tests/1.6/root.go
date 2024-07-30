// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

// TestSchema returns the static schema for a test
// configuration (*.tftest.hcl) file.
func TestSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"run":       runBlockSchema(),
			"provider":  providerBlockSchema(),
			"variables": variablesBlockSchema(),
		},
	}
}
