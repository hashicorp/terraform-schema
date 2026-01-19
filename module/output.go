// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"github.com/zclconf/go-cty/cty"
)

type Output struct {
	Description string
	IsSensitive bool
	Value       cty.Value
}
