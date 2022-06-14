package module

import (
	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type ModuleCalls struct {
	Installed map[string]InstalledModuleCall
	Declared  map[string]DeclaredModuleCall
}

type InstalledModuleCall struct {
	LocalName  string
	SourceAddr string
	Version    *version.Version
	Path       string
}

type DeclaredModuleCall struct {
	LocalName  string
	SourceAddr tfaddr.ModuleSourceRegistry
	Version    version.Constraints
}
