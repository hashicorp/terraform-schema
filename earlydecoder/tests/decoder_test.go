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
	tftest "github.com/hashicorp/terraform-schema/test"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

type testCase struct {
	name          string
	cfg           string
	expectedMeta  *tftest.Meta
	expectedError hcl.Diagnostics
}

var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
}

func TestLoadTest(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"empty config",
			``,
			&tftest.Meta{
				Path:      path,
				Filenames: []string{"test.tftest.hcl"},
			},
			nil,
		},
		// TODO: add more test once we do more early decoding
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), "test.tftest.hcl", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				"test.tftest.hcl": f,
			}

			meta, diags := LoadTest(path, files)

			if diff := cmp.Diff(tc.expectedError, diags, customComparer...); diff != "" {
				t.Fatalf("expected errors doesn't match: %s", diff)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("test meta doesn't match: %s", diff)
			}
		})
	}
}

func compareVersionConstraint(x, y *version.Constraint) bool {
	return x.Equals(y)
}
