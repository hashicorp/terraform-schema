// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/policytest"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// decodedPolicyTest is the type representing a decoded Terraform policytest.
type decodedPolicyTest struct {
	Variables   map[string]*policytest.Variable
	DataSources map[string]*policytest.DataSource
}

func newDecodedPolicyTest() *decodedPolicyTest {
	return &decodedPolicyTest{
		Variables:   make(map[string]*policytest.Variable),
		DataSources: make(map[string]*policytest.DataSource),
	}
}

// loadPolicyTestFromFile reads given file, interprets it and stores in given PolicyTest
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadPolicyTestFromFile(file *hcl.File, decodedPolicyTest *decodedPolicyTest) hcl.Diagnostics {
	var diags hcl.Diagnostics
	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
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

			if attr, defined := content.Attributes["sensitive"]; defined {
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &isSensitive)
				diags = append(diags, valDiags...)
			}

			varType := cty.DynamicPseudoType
			var defaults *typeexpr.Defaults
			if attr, defined := content.Attributes["type"]; defined {
				varType, defaults, valDiags = typeexpr.TypeConstraintWithDefaults(attr.Expr)
				diags = append(diags, valDiags...)
			}

			defaultValue := cty.NilVal
			if attr, defined := content.Attributes["default"]; defined {
				val, vDiags := attr.Expr.Value(nil)
				diags = append(diags, vDiags...)
				if !vDiags.HasErrors() {
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

			decodedPolicyTest.Variables[name] = &policytest.Variable{
				Type:         varType,
				Description:  description,
				DefaultValue: defaultValue,
				TypeDefaults: defaults,
				IsSensitive:  isSensitive,
			}
		}

	}

	return diags
}
