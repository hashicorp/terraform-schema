package module

import (
	"strings"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

var moduleSourceLocalPrefixes = []string{
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

type UnknownSourceAddr string

func (usa UnknownSourceAddr) ForDisplay() string {
	return string(usa)
}
func (usa UnknownSourceAddr) String() string {
	return string(usa)
}

// Parses the raw module source string from a module block
func ParseModuleSourceAddr(source string) ModuleSourceAddr {
	var sourceAddr ModuleSourceAddr
	registryAddr, err := tfaddr.ParseModuleSource(source)
	if err == nil {
		sourceAddr = registryAddr
	} else if isModuleSourceLocal(source) {
		sourceAddr = LocalSourceAddr(source)
	} else if source != "" {
		sourceAddr = UnknownSourceAddr(source)
	}

	return sourceAddr
}

func isModuleSourceLocal(raw string) bool {
	for _, prefix := range moduleSourceLocalPrefixes {
		if strings.HasPrefix(raw, prefix) {
			return true
		}
	}
	return false
}
