package module

import (
	"github.com/hashicorp/go-version"
)

var ModuleSourceLocalPrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

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
	SourceAddr ModuleSourceAddr
	Version    version.Constraints
}

type ModuleSourceAddr interface {
	ForDisplay() string
	String() string
}

type LocalSourceAddr string

func (lsa LocalSourceAddr) ForDisplay() string {
	return string(lsa)
}
func (lsa LocalSourceAddr) String() string {
	return string(lsa)
}
