// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

type testCase struct {
	name          string
	cfg           string
	expectedMeta  *stack.Meta
	expectedError map[string]hcl.Diagnostics
}

var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
}

func TestLoadStack(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"empty config",
			``,
			&stack.Meta{
				Path:                 path,
				Filenames:            []string{"test.tfstack.hcl"},
				Components:           map[string]stack.Component{},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{"test.tfstack.hcl": nil},
		},
		{
			"complete component",
			`component "test" {
	source = "github.com/acme/infra/core"
	version = ">= 1.0, < 2.0"
	inputs = {
		"key" = "value"
	}
}`,
			&stack.Meta{
				Path:      path,
				Filenames: []string{"test.tfstack.hcl"},
				Components: map[string]stack.Component{
					"test": {
						Source:     "github.com/acme/infra/core",
						SourceAddr: module.ParseModuleSourceAddr("git::https://github.com/acme/infra.git//core"),
						Version:    version.MustConstraints(version.NewConstraint(">= 1.0, < 2.0")),
					},
				},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{"test.tfstack.hcl": nil},
		},
		{
			"complete required_providers",
			`required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = "~> 5.7.0"
  }
  random = {
    source  = "hashicorp/random"
    version = "~> 3.5.1"
  }
}`,
			&stack.Meta{
				Path:       path,
				Filenames:  []string{"test.tfstack.hcl"},
				Components: map[string]stack.Component{},
				Variables:  map[string]stack.Variable{},
				Outputs:    map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{
					"aws":    {Source: tfaddr.MustParseProviderSource("hashicorp/aws"), VersionConstraints: version.MustConstraints(version.NewConstraint("~> 5.7.0"))},
					"random": {Source: tfaddr.MustParseProviderSource("hashicorp/random"), VersionConstraints: version.MustConstraints(version.NewConstraint("~> 3.5.1"))},
				},
			},
			map[string]hcl.Diagnostics{"test.tfstack.hcl": nil},
		},
		{
			"variables",
			`variable "example" {
  type    = string
  default = "default_value"
}

variable "example2" {
  description = "description"
  sensitive   = true
}`,
			&stack.Meta{
				Path:       path,
				Filenames:  []string{"test.tfstack.hcl"},
				Components: map[string]stack.Component{},
				Variables: map[string]stack.Variable{
					"example": {
						Type:         cty.String,
						DefaultValue: cty.StringVal("default_value"),
					},
					"example2": {
						Type:        cty.DynamicPseudoType,
						Description: "description",
						IsSensitive: true,
					},
				},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{"test.tfstack.hcl": nil},
		},
		{
			"outputs",
			`output "example" {
  value = "output_value"
  sensitive = true
}

output "example2" {
  description = "description"
  value       = "another_output_value"
}`,
			&stack.Meta{
				Path:       path,
				Filenames:  []string{"test.tfstack.hcl"},
				Components: map[string]stack.Component{},
				Variables:  map[string]stack.Variable{},
				Outputs: map[string]stack.Output{
					"example": {
						Value:       cty.StringVal("output_value"),
						IsSensitive: true,
					},
					"example2": {
						Description: "description",
						Value:       cty.StringVal("another_output_value"),
					},
				},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{"test.tfstack.hcl": nil},
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadStackDiagnostics(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"invalid provider source",
			`required_providers {
  aws = {
    source = "test/test/hashicorp/aws"
  }
}`,
			&stack.Meta{
				Path:                 path,
				Filenames:            []string{"test.tfstack.hcl"},
				Components:           map[string]stack.Component{},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{
				"test.tfstack.hcl": {
					{
						Severity: hcl.DiagError,
						Summary:  `Unable to parse provider source for "aws"`,
						Detail:   `"aws" provider source ("test/test/hashicorp/aws") is not a valid source string`,
						Subject: &hcl.Range{
							Filename: "test.tfstack.hcl",
							Start:    hcl.Pos{Line: 2, Column: 9, Byte: 29},
							End:      hcl.Pos{Line: 4, Column: 4, Byte: 73},
						},
					},
				},
			},
		},
		{
			"invalid provider version",
			`required_providers {
  aws = {
    version = "x~> 5.7.0"
  }
}`,
			&stack.Meta{
				Path:                 path,
				Filenames:            []string{"test.tfstack.hcl"},
				Components:           map[string]stack.Component{},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{
				"test.tfstack.hcl": {
					{
						Severity: hcl.DiagError,
						Summary:  `Unable to parse "aws" provider requirements`,
						Detail:   `Constraint "x~> 5.7.0" is not a valid constraint: Malformed constraint: x~> 5.7.0`,
						Subject: &hcl.Range{
							Filename: "test.tfstack.hcl",
							Start:    hcl.Pos{Line: 2, Column: 9, Byte: 29},
							End:      hcl.Pos{Line: 4, Column: 4, Byte: 60},
						},
					},
				},
			},
		},
		{
			"invalid variable default value",
			`variable "example" {
  type	= string
  default = [1]
}`,
			&stack.Meta{
				Path:       path,
				Filenames:  []string{"test.tfstack.hcl"},
				Components: map[string]stack.Component{},
				Variables: map[string]stack.Variable{
					"example": {
						Type:         cty.String,
						DefaultValue: cty.DynamicVal,
					},
				},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			map[string]hcl.Diagnostics{
				"test.tfstack.hcl": {
					{
						Severity: hcl.DiagError,
						Summary:  `Invalid default value for variable`,
						Detail:   `This default value is not compatible with the variable's type constraint: string required.`,
						Subject: &hcl.Range{
							Filename: "test.tfstack.hcl",
							Start:    hcl.Pos{Line: 3, Column: 13, Byte: 49},
							End:      hcl.Pos{Line: 3, Column: 16, Byte: 52},
						},
					},
				},
			},
		},
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), "test.tfstack.hcl", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				"test.tfstack.hcl": f,
			}

			var fdiags map[string]hcl.Diagnostics
			meta, fdiags := LoadStack(path, files)

			if diff := cmp.Diff(tc.expectedError, fdiags, customComparer...); diff != "" {
				t.Fatalf("expected errors doesn't match: %s", diff)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("stack meta doesn't match: %s", diff)
			}
		})
	}
}

func compareVersionConstraint(x, y *version.Constraint) bool {
	return x.Equals(y)
}
