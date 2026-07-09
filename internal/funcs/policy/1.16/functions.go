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
				Name:        "attrs",
				Type:        cty.DynamicPseudoType,
				Description: "The set of attributes that must match for a resource to be returned.",
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`getresources` gets matching resources from the current Terraform configuration.",
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
				Name:        "attrs",
				Type:        cty.DynamicPseudoType,
				Description: "The set of attributes for the data source.",
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`getdatasource` gets a data source using the current Terraform provider.",
		},
		"semvercmp": {
			Params: []function.Parameter{
				{
					Name:        "version_a",
					Type:        cty.String,
					Description: "First semver version string to compare",
				},
				{
					Name:        "version_b",
					Type:        cty.String,
					Description: "Second semver version string to compare",
				},
			},
			ReturnType:  cty.Number,
			Description: "`semvercmp` compares two semantic versions. Returns -1 if version_a < version_b, 0 if equal, 1 if version_a > version_b.",
		},
		"semverconstraint": {
			Params: []function.Parameter{
				{
					Name:        "version",
					Type:        cty.String,
					Description: "Semantic version string to evaluate",
				},
				{
					Name:        "constraint",
					Type:        cty.String,
					Description: "Version constraint expression (e.g. \">= 1.0.0\", \"<= 5.0.0\")",
				},
			},
			ReturnType:  cty.Bool,
			Description: "`semverconstraint` returns true if the given version satisfies the constraint expression.",
		},
		"gethttprequest": {
			Params: []function.Parameter{
				{
					Name:        "url",
					Type:        cty.String,
					Description: "The URL to make the HTTP request.",
				},
			},
			VarParam: &function.Parameter{
				Name:        "headers",
				Type:        cty.DynamicPseudoType,
				Description: "Optional map of HTTP headers to include in the request",
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`gethttprequest` makes a HTTP GET request using the go net/http package.",
		},
	}
}
