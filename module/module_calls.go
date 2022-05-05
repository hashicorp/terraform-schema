package module

import "github.com/hashicorp/go-version"

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
	SourceAddr string
	Version    version.Constraints
}
