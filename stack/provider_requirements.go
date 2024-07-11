// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

type ProviderRequirement struct {
	Source             string
	VersionConstraints []string
}
