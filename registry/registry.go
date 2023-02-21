// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
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
	Description lang.MarkupContent
	Default     cty.Value
	Required    bool
}

type Output struct {
	Name        string
	Description lang.MarkupContent
}
