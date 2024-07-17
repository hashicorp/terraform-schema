// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type ProviderRequirement struct {
	Source             tfaddr.Provider
	VersionConstraints version.Constraints
}
