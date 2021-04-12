package module

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-registry-address"
)

type Meta struct {
	Path string

	ProviderReferences   map[ProviderRef]tfaddr.Provider
	ProviderRequirements map[tfaddr.Provider]version.Constraints
	CoreRequirements     version.Constraints
}

type ProviderRef struct {
	LocalName string

	// If not empty, Alias identifies which non-default (aliased) provider
	// configuration this address refers to.
	Alias string
}
