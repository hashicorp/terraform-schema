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
