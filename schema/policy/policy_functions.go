// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	pol_v1_16 "github.com/hashicorp/terraform-schema/internal/funcs/policy/1.16"
)

// FunctionsForVersion returns the complete set of functions available
// for policy files at the given Terraform version.
func FunctionsForVersion(v *version.Version) (map[string]schema.FunctionSignature, error) {
	return pol_v1_16.Functions(v), nil
}
