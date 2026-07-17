// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	policytest_v1_16 "github.com/hashicorp/terraform-schema/internal/schema/policytest/1.16"
)

// CorePolicyTestSchemaForVersion finds a policytest schema which is relevant
// for the given Terraform version.
// It will return error if such schema cannot be found.
func CorePolicyTestSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()
	return policytest_v1_16.PolicyTestSchema(ver), nil
}
