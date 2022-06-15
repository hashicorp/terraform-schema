package module

import (
	"github.com/hashicorp/go-version"
	"github.com/zclconf/go-cty/cty"
)

type RegistryModuleMetadataSchema struct {
	Version *version.Version
	Inputs  []RegistryInput
	Outputs []RegistryOutput
}

type RegistryInput struct {
	Name        string    `json:"name"`
	Type        cty.Type  `json:"type"`
	Description string    `json:"description"`
	Default     cty.Value `json:"default"`
	Required    bool      `json:"required"`
}

type RegistryOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
