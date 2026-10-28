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
		"valuetostring": {
			Params: []function.Parameter{
				{
					Name:         "value",
					Type:         cty.DynamicPseudoType,
					Description:  "The value to convert to string.",
					AllowNull:    true,
					AllowUnknown: true,
				},
			},
			ReturnType:  cty.String,
			Description: "`valuetostring` converts any value to a string. Primitives are rendered directly; complex types (objects, lists, maps, sets, tuples) are encoded as indented JSON. Unknown leaves render as `(unknown)` and null leaves as `(null)`.",
		},
		"contains_substring": {
			Params: []function.Parameter{
				{
					Name:        "str",
					Type:        cty.String,
					Description: "The string to search",
				},
				{
					Name:        "substr",
					Type:        cty.String,
					Description: "The substring to look for",
				},
			},
			ReturnType:  cty.Bool,
			Description: "`contains_substring` returns true if `str` contains `substr`.",
		},
		"regex_match": {
			Params: []function.Parameter{
				{
					Name:        "str",
					Type:        cty.String,
					Description: "The string to match",
				},
				{
					Name:        "pattern",
					Type:        cty.String,
					Description: "The regex pattern to match against",
				},
			},
			ReturnType:  cty.Bool,
			Description: "`regex_match` returns true if the string matches the regex pattern, false otherwise.",
		},
		"type": {
			Params: []function.Parameter{
				{
					Name:        "value",
					Type:        cty.DynamicPseudoType,
					Description: "The value to get the type of.",
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`type` returns the type of the given value. Useful for debugging.",
		},
	}
}
