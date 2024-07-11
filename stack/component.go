// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import "github.com/hashicorp/go-version"

type Component struct {
	Source  string
	Version version.Constraints
}
