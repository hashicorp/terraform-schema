// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/policy"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

type testCase struct {
	name          string
	cfg           string
	fileName      string
	expectedMeta  *policy.Meta
	expectedDiags hcl.Diagnostics
}

// Using your specific comparator requirements
var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
	// We ignore the 'Range' fields in these tests because they track exact
	// line, column, and byte offsets. Testing these would make the tests
	// brittle—changing a single space or newline in the test HCL string
	// would cause the test to fail even if the logic is correct.
	cmpopts.IgnoreFields(policy.ResourcePolicy{}, "Range"),
	cmpopts.IgnoreFields(policy.ProviderPolicy{}, "Range"),
	cmpopts.IgnoreFields(policy.ModulePolicy{}, "Range"),
}

func TestLoadPolicy(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			name:     "empty config",
			fileName: "empty.policy.hcl",
			cfg:      ``,
			expectedMeta: &policy.Meta{
				Path:             path,
				Filenames:        []string{"empty.policy.hcl"},
				ResourcePolicies: map[string]policy.ResourcePolicy{},
				ProviderPolicies: map[string]policy.ProviderPolicy{},
				ModulePolicies:   map[string]policy.ModulePolicy{},
			},
		},
		{
			name:     "core requirements",
			fileName: "version.policy.hcl",
			cfg: `
policy {
  terraform_config {
    required_version = ">= 1.12"
  }
}`,
			expectedMeta: &policy.Meta{
				Path:             path,
				Filenames:        []string{"version.policy.hcl"},
				CoreRequirements: version.MustConstraints(version.NewConstraint(">= 1.12")),
				ResourcePolicies: map[string]policy.ResourcePolicy{},
				ProviderPolicies: map[string]policy.ProviderPolicy{},
				ModulePolicies:   map[string]policy.ModulePolicy{},
			},
		},
		{
			name:     "resource policy",
			fileName: "resource.policy.hcl",
			cfg: `
resource_policy "aws_instance" "web" {
  filter = attrs.type == "t3.micro"
}`,
			expectedMeta: &policy.Meta{
				Path:      path,
				Filenames: []string{"resource.policy.hcl"},
				ResourcePolicies: map[string]policy.ResourcePolicy{
					"aws_instance.web": {Type: "aws_instance", Name: "web"},
				},
				ProviderPolicies: map[string]policy.ProviderPolicy{},
				ModulePolicies:   map[string]policy.ModulePolicy{},
			},
		},
		{
			name:     "module policy",
			fileName: "module.policy.hcl",
			cfg:      `module_policy "./modules/vpc" "net" {}`,
			expectedMeta: &policy.Meta{
				Path:             path,
				Filenames:        []string{"module.policy.hcl"},
				ResourcePolicies: map[string]policy.ResourcePolicy{},
				ProviderPolicies: map[string]policy.ProviderPolicy{},
				ModulePolicies: map[string]policy.ModulePolicy{
					"./modules/vpc.net": {Type: "./modules/vpc", Name: "net"},
				},
			},
		},
		{
			name:     "provider policy",
			fileName: "provider.policy.hcl",
			cfg:      `provider_policy "hashicorp/aws" "main" {}`,
			expectedMeta: &policy.Meta{
				Path:             path,
				Filenames:        []string{"provider.policy.hcl"},
				ResourcePolicies: map[string]policy.ResourcePolicy{},
				ModulePolicies:   map[string]policy.ModulePolicy{},
				ProviderPolicies: map[string]policy.ProviderPolicy{
					"hashicorp/aws.main": {Type: "hashicorp/aws", Name: "main"},
				},
			},
		},
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), tc.fileName, hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("hcl syntax error: %s", diags)
			}

			files := map[string]*hcl.File{
				tc.fileName: f,
			}

			meta, fdiags := LoadPolicy(path, files)

			// Match Diagnostics Summary
			if len(tc.expectedDiags) != len(fdiags) {
				t.Errorf("expected %d diagnostics, got %d: %s", len(tc.expectedDiags), len(fdiags), fdiags)
			} else {
				for j, expected := range tc.expectedDiags {
					if fdiags[j].Summary != expected.Summary {
						t.Errorf("diag %d summary mismatch: want %q, got %q", j, expected.Summary, fdiags[j].Summary)
					}
				}
			}

			// Match Meta
			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("policy meta mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func compareVersionConstraint(x, y version.Constraints) bool {
	return x.String() == y.String()
}
