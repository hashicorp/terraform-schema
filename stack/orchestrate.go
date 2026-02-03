// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package stack

import "github.com/hashicorp/hcl/v2"

type OrchestrationRule struct {
	Type  string
	Range hcl.Range
}
