// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	funcs_v0_14 "github.com/hashicorp/terraform-schema/internal/funcs/0.14"
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	f := funcs_v0_14.Functions(v)

	f["nonsensitive"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.DynamicPseudoType,
			},
		},
		ReturnType:  cty.DynamicPseudoType,
		Description: "`nonsensitive` takes a sensitive value and returns a copy of that value with the sensitive marking removed, thereby exposing the sensitive value.",
	}
	f["one"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.DynamicPseudoType,
			},
		},
		ReturnType:  cty.DynamicPseudoType,
		Description: "`one` takes a list, set, or tuple value with either zero or one elements. If the collection is empty, `one` returns `null`. Otherwise, `one` returns the first element. If there are two or more elements then `one` will return an error.",
	}
	f["sensitive"] = schema.FunctionSignature{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.DynamicPseudoType,
			},
		},
		ReturnType:  cty.DynamicPseudoType,
		Description: "`sensitive` takes any value and returns a copy of it marked so that Terraform will treat it as sensitive, with the same meaning and behavior as for [sensitive input variables](/language/values/variables#suppressing-values-in-cli-output).",
	}

	delete(f, "list") // list was removed in 0.15
	delete(f, "map")  // map was removed in 0.15

	return f
}
