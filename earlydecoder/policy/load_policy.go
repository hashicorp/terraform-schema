// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/policy"
)

// decodedPolicy is the type representing a decoded Terraform policy.
type decodedPolicy struct {
	RequiredCore     []string
	ResourcePolicies map[string]*policy.ResourcePolicy
	ProviderPolicies map[string]*policy.ProviderPolicy
	ModulePolicies   map[string]*policy.ModulePolicy
}

func newDecodedPolicy() *decodedPolicy {
	return &decodedPolicy{
		RequiredCore:     make([]string, 0),
		ResourcePolicies: make(map[string]*policy.ResourcePolicy),
		ProviderPolicies: make(map[string]*policy.ProviderPolicy),
		ModulePolicies:   make(map[string]*policy.ModulePolicy),
	}
}

// loadPolicyFromFile reads given file, interprets it and stores in given Policy
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadPolicyFromFile(file *hcl.File, decodedPolicy *decodedPolicy) hcl.Diagnostics {
	var diags hcl.Diagnostics
	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "policy":
			content, _, contentDiags := block.Body.PartialContent(policyBlockSchema)
			diags = append(diags, contentDiags...)

			for _, block := range content.Blocks {
				switch block.Type {
				case "terraform_config":
					content, _, contentDiags := block.Body.PartialContent(terraformConfigBlockSchema)
					diags = append(diags, contentDiags...)

					if attr, defined := content.Attributes["required_version"]; defined {
						var version string
						valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
						diags = append(diags, valDiags...)
						if !valDiags.HasErrors() {
							decodedPolicy.RequiredCore = append(decodedPolicy.RequiredCore, version)
						}
					}
				}
			}

		case "resource_policy":
			_, _, contentDiags := block.Body.PartialContent(resourcePolicyBlockSchema)
			diags = append(diags, contentDiags...)

			if len(block.Labels) != 2 || block.Labels[0] == "" || block.Labels[1] == "" {
				continue
			}

			body, ok := block.Body.(*hclsyntax.Body)
			if !ok {
				continue
			}
			resType := block.Labels[0]
			resName := block.Labels[1]
			key := resType + "." + resName

			decodedPolicy.ResourcePolicies[key] = &policy.ResourcePolicy{
				Type:  resType,
				Name:  resName,
				Range: body.SrcRange,
			}

		case "provider_policy":
			_, _, contentDiags := block.Body.PartialContent(providerPolicyBlockSchema)
			diags = append(diags, contentDiags...)

			if len(block.Labels) != 2 || block.Labels[0] == "" || block.Labels[1] == "" {
				continue
			}

			body, ok := block.Body.(*hclsyntax.Body)
			if !ok {
				continue
			}
			resType := block.Labels[0]
			resName := block.Labels[1]
			key := resType + "." + resName

			decodedPolicy.ProviderPolicies[key] = &policy.ProviderPolicy{
				Type:  resType,
				Name:  resName,
				Range: body.SrcRange,
			}

		case "module_policy":
			_, _, contentDiags := block.Body.PartialContent(modulePolicyBlockSchema)
			diags = append(diags, contentDiags...)

			if len(block.Labels) != 2 || block.Labels[0] == "" || block.Labels[1] == "" {
				continue
			}

			body, ok := block.Body.(*hclsyntax.Body)
			if !ok {
				continue
			}
			resType := block.Labels[0]
			resName := block.Labels[1]
			key := resType + "." + resName

			decodedPolicy.ModulePolicies[key] = &policy.ModulePolicy{
				Type:  resType,
				Name:  resName,
				Range: body.SrcRange,
			}
		}

	}

	return diags
}
