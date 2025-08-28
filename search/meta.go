// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package search

import (
	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type Meta struct {
	Path      string
	Filenames []string

	Variables          map[string]Variable
	Lists              map[string]List
	ProviderReferences map[ProviderRef]tfaddr.Provider
}

type ProviderRef struct {
	LocalName string

	// If not empty, Alias identifies which non-default (aliased) provider
	// configuration this address refers to.
	Alias string
}

type ProviderRequirements map[tfaddr.Provider]version.Constraints

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
