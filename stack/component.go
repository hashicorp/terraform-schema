// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-schema/module"
)

type Component struct {
	Source     string
	SourceAddr module.ModuleSourceAddr
	Version    version.Constraints
}
