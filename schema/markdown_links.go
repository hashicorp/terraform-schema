// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"regexp"
	"strings"

	hclschema "github.com/hashicorp/hcl-lang/schema"
)

const (
	developerDocsBaseURL = "https://developer.hashicorp.com"
	terraformDocsBaseURL = developerDocsBaseURL + "/terraform"
)

var markdownInlineRootLink = regexp.MustCompile(`(\]\()(/[^)\s]+)([^)]*\))`)
var terraformDocsPathPrefixes = []string{
	"/cli/",
	"/commands/",
	"/cloud-docs/",
	"/internals/",
	"/language/",
}

func resolveFunctionMarkdownLinks(functions map[string]hclschema.FunctionSignature) map[string]hclschema.FunctionSignature {
	normalized := make(map[string]hclschema.FunctionSignature, len(functions))
	for name, fSig := range functions {
		fCopy := *fSig.Copy()
		fCopy.Description = resolveTerraformMarkdownLinks(fCopy.Description)
		for i := range fCopy.Params {
			fCopy.Params[i].Description = resolveTerraformMarkdownLinks(fCopy.Params[i].Description)
		}
		if fCopy.VarParam != nil {
			fCopy.VarParam.Description = resolveTerraformMarkdownLinks(fCopy.VarParam.Description)
		}
		normalized[name] = fCopy
	}

	return normalized
}

func resolveTerraformMarkdownLinks(markdown string) string {
	return markdownInlineRootLink.ReplaceAllStringFunc(markdown, func(link string) string {
		matches := markdownInlineRootLink.FindStringSubmatch(link)
		if len(matches) != 4 {
			return link
		}

		return matches[1] + absoluteTerraformDocsURL(matches[2]) + matches[3]
	})
}

func absoluteTerraformDocsURL(href string) string {
	if strings.HasPrefix(href, "/terraform/") {
		return developerDocsBaseURL + href
	}
	for _, prefix := range terraformDocsPathPrefixes {
		if strings.HasPrefix(href, prefix) {
			return terraformDocsBaseURL + href
		}
	}

	return href
}
