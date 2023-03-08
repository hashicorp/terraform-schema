// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/backend"
)

type Meta struct {
	Path      string
	Filenames []string

	Backend              *Backend
	Cloud                *backend.Cloud
	ProviderReferences   map[ProviderRef]tfaddr.Provider
	ProviderRequirements ProviderRequirements
	CoreRequirements     version.Constraints
	Variables            map[string]Variable
	Outputs              map[string]Output
	ModuleCalls          map[string]DeclaredModuleCall
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

type Backend struct {
	Type string
	Data backend.BackendData
}

func (be *Backend) Equals(b *Backend) bool {
	if be == nil && b == nil {
		return true
	}

	if be == nil || b == nil {
		return false
	}

	if be.Type != b.Type {
		return false
	}

	return be.Data.Equals(b.Data)
}

type ProviderRef struct {
	LocalName string

	// If not empty, Alias identifies which non-default (aliased) provider
	// configuration this address refers to.
	Alias string
}
