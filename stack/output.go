// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"github.com/zclconf/go-cty/cty"
)

type Output struct {
	Description string
	IsSensitive bool
	Value       cty.Value
}
