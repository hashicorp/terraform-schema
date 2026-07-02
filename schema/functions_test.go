// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestFunctionsForVersion_resolvesTerraformMarkdownLinks(t *testing.T) {
	testCases := []struct {
		name    string
		version string
	}{
		{
			name:    "pre product path",
			version: "1.5.0",
		},
		{
			name:    "with product path",
			version: "1.10.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			functions, err := FunctionsForVersion(version.Must(version.NewVersion(tc.version)))
			if err != nil {
				t.Fatal(err)
			}

			sensitiveDescription := functions["sensitive"].Description
			if !strings.Contains(sensitiveDescription, "[sensitive input variables](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output)") {
				t.Fatalf("sensitive description did not include absolute Terraform docs link: %q", sensitiveDescription)
			}
		})
	}
}
