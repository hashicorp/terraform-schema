// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	funcs_v0_13 "github.com/hashicorp/terraform-schema/internal/funcs/0.13"
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	f := funcs_v0_13.Functions(v)

	f["alltrue"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.List(cty.Bool),
			},
		},
		ReturnType:  cty.Bool,
		Description: "`alltrue` returns `true` if all elements in a given collection are `true` or `\"true\"`. It also returns `true` if the collection is empty.",
	}
	f["anytrue"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.List(cty.Bool),
			},
		},
		ReturnType:  cty.Bool,
		Description: "`anytrue` returns `true` if any element in a given collection is `true` or `\"true\"`. It also returns `false` if the collection is empty.",
	}
	f["textdecodebase64"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "source",
				Type: cty.String,
			},
			{
				Name: "encoding",
				Type: cty.String,
			},
		},
		ReturnType:  cty.String,
		Description: "`textdecodebase64` function decodes a string that was previously Base64-encoded, and then interprets the result as characters in a specified character encoding.",
	}
	f["textencodebase64"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "string",
				Type: cty.String,
			},
			{
				Name: "encoding",
				Type: cty.String,
			},
		},
		ReturnType:  cty.String,
		Description: "`textencodebase64` encodes the unicode characters in a given string using a specified character encoding, returning the result base64 encoded because Terraform language strings are always sequences of unicode characters.",
	}

	return f
}
