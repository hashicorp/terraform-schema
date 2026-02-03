// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package policy

import "github.com/hashicorp/hcl/v2"

type ResourcePolicy struct {
	// Type is the resource type (first label)
	Type string

	// Name is the policy name (second label)
	Name string

	// Range is the range of the block declaration
	Range hcl.Range
}
