// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Description string
	Type        cty.Type

	// DefaultValue represents default value if one is defined
	// and is decodable without errors, else cty.NilVal
	DefaultValue cty.Value

	// TypeDefaults represents any default values for optional object
	// attributes assuming Type is of cty.Object and has defaults.
	//
	// Any relationships between DefaultValue & TypeDefaults are left
	// for downstream to deal with using e.g. TypeDefaults.Apply().
	TypeDefaults *typeexpr.Defaults
}
