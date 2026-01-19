// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/detect"
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
	SourceAddr ModuleSourceAddr
	Version    *version.Version
	Path       string
}

type DeclaredModuleCall struct {
	LocalName     string
	RawSourceAddr string
	SourceAddr    ModuleSourceAddr
	Version       version.Constraints
	InputNames    []string
	RangePtr      *hcl.Range
}

func (mc DeclaredModuleCall) Copy() DeclaredModuleCall {
	inputNames := make([]string, len(mc.InputNames))
	copy(inputNames, mc.InputNames)

	newModuleCall := DeclaredModuleCall{
		LocalName:     mc.LocalName,
		RawSourceAddr: mc.RawSourceAddr,
		SourceAddr:    mc.SourceAddr,
		Version:       mc.Version,
		InputNames:    inputNames,
	}

	if mc.RangePtr != nil {
		rangeCpy := *mc.RangePtr
		newModuleCall.RangePtr = &rangeCpy
	}

	return newModuleCall
}

type ModuleSourceAddr interface {
	ForDisplay() string
	String() string
}

// LocalSourceAddr represents a local module source address containing a path
// to a local directory
type LocalSourceAddr string

func (lsa LocalSourceAddr) ForDisplay() string {
	return string(lsa)
}
func (lsa LocalSourceAddr) String() string {
	return string(lsa)
}

// RemoteSourceAddr represents a remote module source address containing
// a source string normalized by go-getter's detection logic
type RemoteSourceAddr string

func (rsa RemoteSourceAddr) ForDisplay() string {
	return string(rsa)
}
func (rsa RemoteSourceAddr) String() string {
	return string(rsa)
}

// UnknownSourceAddr represents an unknown module source address as a fallback
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
	} else if remoteAddr, err := parseRemoteModuleSource(source); err == nil {
		sourceAddr = RemoteSourceAddr(remoteAddr)
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

func parseRemoteModuleSource(raw string) (string, error) {
	return detect.Detect(raw, detect.RemoteSourceDetectors)
}
