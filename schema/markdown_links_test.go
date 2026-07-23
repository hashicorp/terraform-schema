// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/hcl-lang/lang"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestResolveTerraformMarkdownLinks(t *testing.T) {
	testCases := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "adds terraform product path",
			markdown: "[sensitive input variables](/language/values/variables#suppressing-values-in-cli-output)",
			expected: "[sensitive input variables](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output)",
		},
		{
			name:     "keeps existing terraform product path",
			markdown: "[sensitive input variables](/terraform/language/values/variables#suppressing-values-in-cli-output)",
			expected: "[sensitive input variables](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output)",
		},
		{
			name:     "preserves absolute links",
			markdown: "[Terraform](https://developer.hashicorp.com/terraform)",
			expected: "[Terraform](https://developer.hashicorp.com/terraform)",
		},
		{
			name:     "preserves non-root relative links",
			markdown: "[the main provider documentation](../index.html)",
			expected: "[the main provider documentation](../index.html)",
		},
		{
			name:     "preserves other root relative links",
			markdown: "[provider docs](/docs/providers/aws/index.html)",
			expected: "[provider docs](/docs/providers/aws/index.html)",
		},
		{
			name:     "preserves fragment links",
			markdown: "[details](#details)",
			expected: "[details](#details)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveTerraformMarkdownLinks(tc.markdown)
			if actual != tc.expected {
				t.Fatalf("unexpected markdown\nexpected: %q\nactual:   %q", tc.expected, actual)
			}
		})
	}
}

func TestMarkupContent_resolvesTerraformMarkdownLinks(t *testing.T) {
	content := markupContent("[sensitive input variables](/language/values/variables#suppressing-values-in-cli-output)", tfjson.SchemaDescriptionKindMarkdown)

	if content.Kind != lang.MarkdownKind {
		t.Fatalf("unexpected markup kind: %s", content.Kind)
	}

	expected := "[sensitive input variables](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output)"
	if content.Value != expected {
		t.Fatalf("unexpected markup content\nexpected: %q\nactual:   %q", expected, content.Value)
	}
}
