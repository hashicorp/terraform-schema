package module

import (
	"github.com/hashicorp/go-version"
)

type RegistryModuleMetadataSchema struct {
	Version *version.Version
	Inputs  []RegistryInput
	Outputs []RegistryOutput
}

type RegistryInput struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

type RegistryOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
