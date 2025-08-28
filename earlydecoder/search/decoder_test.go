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
	"github.com/hashicorp/terraform-schema/search"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

type testCase struct {
	name          string
	cfg           string
	fileName      string
	expectedMeta  *search.Meta
	expectedError map[string]hcl.Diagnostics
}

var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
}

var fileName = "test.tfquery.hcl"

func TestLoadSearch(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"empty config",
			``,
			fileName,
			&search.Meta{
				Path:               path,
				Filenames:          []string{fileName},
				Variables:          map[string]search.Variable{},
				Lists:              map[string]search.List{},
				ProviderReferences: map[search.ProviderRef]tfaddr.Provider{},
			},
			map[string]hcl.Diagnostics{fileName: nil},
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
			fileName,
			&search.Meta{

				Path:      path,
				Filenames: []string{fileName},
				Variables: map[string]search.Variable{
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
				Lists:              map[string]search.List{},
				ProviderReferences: map[search.ProviderRef]tfaddr.Provider{},
			},
			map[string]hcl.Diagnostics{fileName: nil},
		},
	}

	runTestCases(testCases, t, path)

}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), tc.fileName, hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				tc.fileName: f,
			}

			var fdiags map[string]hcl.Diagnostics
			meta, fdiags := LoadSearch(path, files)

			if diff := cmp.Diff(tc.expectedError, fdiags, customComparer...); diff != "" {
				t.Fatalf("expected errors doesn't match: %s", diff)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("search meta doesn't match: %s", diff)
			}
		})
	}
}

func TestLoadSearchDiagnostics(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"invalid variable default value",
			`variable "example" {
  type	= string
  default = [1]
}`,
			fileName,
			&search.Meta{
				Path:      path,
				Filenames: []string{fileName},
				Variables: map[string]search.Variable{
					"example": {
						Type:         cty.String,
						DefaultValue: cty.DynamicVal,
					},
				},
				Lists:              map[string]search.List{},
				ProviderReferences: map[search.ProviderRef]tfaddr.Provider{},
			},
			map[string]hcl.Diagnostics{
				fileName: {
					{
						Severity: hcl.DiagError,
						Summary:  `Invalid default value for variable`,
						Detail:   `This default value is not compatible with the variable's type constraint: string required, but have tuple.`,
						Subject: &hcl.Range{
							Filename: fileName,
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

func compareVersionConstraint(x, y *version.Constraint) bool {
	return x.Equals(y)
}
