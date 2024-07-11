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
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

type testCase struct {
	name          string
	cfg           string
	expectedMeta  *stack.Meta
	expectedError hcl.Diagnostics
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
				Filenames:            []string{"test.tf"},
				Components:           map[string]stack.Component{},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			nil,
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
				Path:                 path,
				Filenames:            []string{"test.tf"},
				Components:           map[string]stack.Component{"test": {Source: "github.com/acme/infra/core", Version: version.MustConstraints(version.NewConstraint(">= 1.0, < 2.0"))}},
				Variables:            map[string]stack.Variable{},
				Outputs:              map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{},
			},
			nil,
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
				Filenames:  []string{"test.tf"},
				Components: map[string]stack.Component{},
				Variables:  map[string]stack.Variable{},
				Outputs:    map[string]stack.Output{},
				ProviderRequirements: map[string]stack.ProviderRequirement{
					"aws":    {Source: "hashicorp/aws", VersionConstraints: []string{"~> 5.7.0"}},
					"random": {Source: "hashicorp/random", VersionConstraints: []string{"~> 3.5.1"}},
				},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), "test.tf", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				"test.tf": f,
			}

			meta, diags := LoadStack(path, files)

			if diff := cmp.Diff(tc.expectedError, diags, customComparer...); diff != "" {
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
