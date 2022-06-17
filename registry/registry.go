package registry

import (
	"github.com/hashicorp/go-version"
	"github.com/zclconf/go-cty/cty"
)

type ModuleData struct {
	Version *version.Version
	Inputs  []Input
	Outputs []Output
}

type Input struct {
	Name        string
	Type        cty.Type
	Description string
	Default     cty.Value
	Required    bool
}

type Output struct {
	Name        string
	Description string
}
