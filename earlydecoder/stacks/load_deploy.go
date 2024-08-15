// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty/cty"
)

// loadDeployFromFile reads given file, interprets it and stores in given stack
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadDeployFromFile(file *hcl.File, ds *decodedStack) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, _, contentDiags := file.Body.PartialContent(deploymentRootSchema)
	diags = append(diags, contentDiags...)
	for _, block := range content.Blocks {
		switch block.Type {
		case "deployment":
			content, _, contentDiags := block.Body.PartialContent(deploymentSchema)
			diags = append(diags, contentDiags...)

			if len(block.Labels) != 1 || block.Labels[0] == "" {
				continue
			}

			name := block.Labels[0]

			inputs := make(map[string]cty.Value)
			if attr, defined := content.Attributes["inputs"]; defined {
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &inputs)
				diags = append(diags, valDiags...)
			}

			ds.Deployments[name] = &stack.Deployment{
				Inputs: inputs,
			}
		}
	}

	return diags
}