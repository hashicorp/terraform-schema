// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/policy"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// decodedPolicy is the type representing a decoded Terraform policy.
type decodedPolicy struct {
	RequiredCore     []string
	Variables        map[string]*policy.Variable
	ResourcePolicies map[string]*policy.ResourcePolicy
	ProviderPolicies map[string]*policy.ProviderPolicy
	ModulePolicies   map[string]*policy.ModulePolicy
}

func newDecodedPolicy() *decodedPolicy {
	return &decodedPolicy{
		RequiredCore:     make([]string, 0),
		Variables:        make(map[string]*policy.Variable),
		ResourcePolicies: make(map[string]*policy.ResourcePolicy),
		ProviderPolicies: make(map[string]*policy.ProviderPolicy),
		ModulePolicies:   make(map[string]*policy.ModulePolicy),
	}
}

// loadPolicyFromFile reads given file, interprets it and stores in given Policy
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadPolicyFromFile(file *hcl.File, mod *decodedPolicy) hcl.Diagnostics {
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
							mod.RequiredCore = append(mod.RequiredCore, version)
						}
					}
				}
			}

		case "variable":
			content, _, contentDiags := block.Body.PartialContent(variableSchema)
			diags = append(diags, contentDiags...)
			if len(block.Labels) != 1 || block.Labels[0] == "" {
				continue
			}
			name := block.Labels[0]
			description := ""
			isSensitive := false
			var valDiags hcl.Diagnostics
			if attr, defined := content.Attributes["description"]; defined {
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &description)
				diags = append(diags, valDiags...)
			}
			varType := cty.DynamicPseudoType
			var defaults *typeexpr.Defaults
			if attr, defined := content.Attributes["type"]; defined {
				varType, defaults, valDiags = typeexpr.TypeConstraintWithDefaults(attr.Expr)
				diags = append(diags, valDiags...)
			}
			if attr, defined := content.Attributes["sensitive"]; defined {
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &isSensitive)
				diags = append(diags, valDiags...)
			}
			defaultValue := cty.NilVal
			if attr, defined := content.Attributes["default"]; defined {
				val, diags := attr.Expr.Value(nil)
				if !diags.HasErrors() {
					if varType != cty.NilType {
						var err error
						val, err = convert.Convert(val, varType)
						if err != nil {
							diags = append(diags, &hcl.Diagnostic{
								Severity: hcl.DiagError,
								Summary:  "Invalid default value for variable",
								Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
								Subject:  attr.Expr.Range().Ptr(),
							})
							val = cty.DynamicVal
						}
					}

					defaultValue = val
				}
			}
			mod.Variables[name] = &policy.Variable{
				Type:         varType,
				Description:  description,
				IsSensitive:  isSensitive,
				DefaultValue: defaultValue,
				TypeDefaults: defaults,
			}

		}

	}

	return diags
}
