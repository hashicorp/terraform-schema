// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import "github.com/zclconf/go-cty/cty"

type Deployment struct {
	Inputs map[string]cty.Value
}
