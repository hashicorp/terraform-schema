package module

import (
	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/backend"
)

type Meta struct {
	Path string

	Backend              *Backend
	ProviderReferences   map[ProviderRef]tfaddr.Provider
	ProviderRequirements map[tfaddr.Provider]version.Constraints
	CoreRequirements     version.Constraints
	Variables            map[string]Variable
	Outputs              map[string]Output
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

type ModuleCall struct {
	LocalName  string
	SourceAddr string
	Path       string
}
