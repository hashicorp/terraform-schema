// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/reference"

	refs_v0_12 "github.com/hashicorp/terraform-schema/internal/references/0.12"
	refs_v1_10 "github.com/hashicorp/terraform-schema/internal/references/1.10"
)

// BuiltinReferencesForVersion returns known "built-in" reference targets
// (range-less references available within any module)
func BuiltinReferencesForVersion(v *version.Version, modPath string) reference.Targets {
	ver := v.Core()

	if ver.GreaterThanOrEqual(v1_10) {
		return refs_v1_10.BuiltinReferences(modPath)
	}

	return refs_v0_12.BuiltinReferences(modPath)
}
