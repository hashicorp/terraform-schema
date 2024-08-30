// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import "github.com/hashicorp/hcl/v2"

type OrchestrationRule struct {
	Type  string
	Range hcl.Range
}
