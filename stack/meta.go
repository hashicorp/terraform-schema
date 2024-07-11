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

type ProviderRequirements map[tfaddr.Provider]version.Constraints // ??

func (pr ProviderRequirements) Equals(reqs ProviderRequirements) bool {
	if len(pr) != len(reqs) {
		return false
	}

	for pAddr, vCons := range pr {
		c, ok := reqs[pAddr]
		if !ok {
			return false
		}
		if !vCons.Equals(c) {
			return false
		}
	}

	return true
}

type ProviderRef struct {
	LocalName string

	// If not empty, Alias identifies which non-default (aliased) provider
	// configuration this address refers to.
	Alias string
}
