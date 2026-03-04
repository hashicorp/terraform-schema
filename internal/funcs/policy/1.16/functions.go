// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	tfschema "github.com/hashicorp/terraform-schema/schema"
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	coreFunctions, err := tfschema.FunctionsForVersion(v)
	if err != nil {
		// this should never happen
		panic(err)
	}

	policyFns := policyFunctions()

	functions := make(map[string]schema.FunctionSignature, len(coreFunctions)+len(policyFns))

	for name, signature := range coreFunctions {
		functions["core::"+name] = signature
	}

	// Policy-specific functions
	for name, signature := range policyFns {
		functions["core::"+name] = signature
	}

	return functions
}

func policyFunctions() map[string]schema.FunctionSignature {
	return map[string]schema.FunctionSignature{
		"getresources": {
			Params: []function.Parameter{
				{
					Name:        "type",
					Type:        cty.String,
					Description: "The resource type to retrieve.",
				},
			},
			VarParam: &function.Parameter{
				Name:        "config",
				Type:        cty.DynamicPseudoType,
				Description: "Optional configuration for filtering resources.",
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`getresources` retrieves resources of the given type from the current policy context.",
		},
		"getdatasource": {
			Params: []function.Parameter{
				{
					Name:        "type",
					Type:        cty.String,
					Description: "The data source type to retrieve.",
				},
			},
			VarParam: &function.Parameter{
				Name:        "config",
				Type:        cty.DynamicPseudoType,
				Description: "Optional configuration for filtering data sources.",
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`getdatasource` retrieves data sources of the given type from the current policy context.",
		},
	}
}
