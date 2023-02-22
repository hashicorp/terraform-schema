// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	funcs_v0_12 "github.com/hashicorp/terraform-schema/internal/funcs/0.12"
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	f := funcs_v0_12.Functions(v)

	f["sum"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.DynamicPseudoType,
			},
		},
		ReturnType:  cty.DynamicPseudoType,
		Description: "`sum` takes a list or set of numbers and returns the sum of those numbers.",
	}

	return f
}
