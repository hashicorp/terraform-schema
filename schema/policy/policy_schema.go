// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	pol_v1_16 "github.com/hashicorp/terraform-schema/internal/schema/policy/1.16"
)

// CorePolicySchemaForVersion finds a policy schema which is relevant
// for the given Terraform version.
// It will return error if such schema cannot be found.
func CorePolicySchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()
	return pol_v1_16.PolicySchema(ver), nil
}
