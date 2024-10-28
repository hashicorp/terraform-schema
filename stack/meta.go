// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

type Meta struct {
	Path      string
	Filenames []string

	Components           map[string]Component
	Variables            map[string]Variable
	Outputs              map[string]Output
	ProviderRequirements map[string]ProviderRequirement
}
