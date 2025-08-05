// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/search"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type providerConfig struct {
	Name  string
	Alias string
}

type decodedSearch struct {
	List            map[string]*search.List
	Variables       map[string]*search.Variable
	ProviderConfigs map[string]*providerConfig
}

var providerConfigSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "version",
		},
		{
			Name: "alias",
		},
	},
}

func newDecodedSearch() *decodedSearch {
	return &decodedSearch{
		List:            make(map[string]*search.List),
		Variables:       make(map[string]*search.Variable),
		ProviderConfigs: make(map[string]*providerConfig),
	}
}

func loadSearchFromFile(file *hcl.File, ds *decodedSearch) hcl.Diagnostics {
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

			ds.Variables[name] = &search.Variable{
				Type:         varType,
				Description:  description,
				DefaultValue: defaultValue,
				TypeDefaults: defaults,
				IsSensitive:  isSensitive,
			}

		case "provider":
			content, _, contentDiags := block.Body.PartialContent(providerConfigSchema)
			diags = append(diags, contentDiags...)
			name := block.Labels[0]
			providerKey := name
			var alias string
			if attr, defined := content.Attributes["alias"]; defined {
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &alias)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() && alias != "" {
					providerKey = fmt.Sprintf("%s.%s", name, alias)
				}
			}

			ds.ProviderConfigs[providerKey] = &providerConfig{
				Name:  name,
				Alias: alias,
			}
		}

	}

	return diags
}
