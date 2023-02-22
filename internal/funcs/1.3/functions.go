// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	funcs_v0_15 "github.com/hashicorp/terraform-schema/internal/funcs/0.15"
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	f := funcs_v0_15.Functions(v)

	f["endswith"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "str",
				Type: cty.String,
			},
			{
				Name: "suffix",
				Type: cty.String,
			},
		},
		ReturnType:  cty.Bool,
		Description: "`endswith` takes two values: a string to check and a suffix string. The function returns true if the first string ends with that exact suffix.",
	}
	f["startswith"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "str",
				Type: cty.String,
			},
			{
				Name: "prefix",
				Type: cty.String,
			},
		},
		ReturnType:  cty.Bool,
		Description: "`startswith` takes two values: a string to check and a prefix string. The function returns true if the string begins with that exact prefix.",
	}
	f["timecmp"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "timestamp_a",
				Type: cty.String,
			},
			{
				Name: "timestamp_b",
				Type: cty.String,
			},
		},
		ReturnType:  cty.Number,
		Description: "`timecmp` compares two timestamps and returns a number that represents the ordering of the instants those timestamps represent.",
	}

	return f
}
