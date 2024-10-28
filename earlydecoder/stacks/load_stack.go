// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// decodedStack is the type representing a decoded Terraform stack.
type decodedStack struct {
	Components           map[string]*stack.Component
	Variables            map[string]*stack.Variable
	Outputs              map[string]*stack.Output
	ProviderRequirements map[string]*providerRequirement
}

func newDecodedStack() *decodedStack {
	return &decodedStack{
		Components:           make(map[string]*stack.Component),
		Variables:            make(map[string]*stack.Variable),
		Outputs:              make(map[string]*stack.Output),
		ProviderRequirements: make(map[string]*providerRequirement),
	}
}

// loadStackFromFile reads given file, interprets it and stores in given Stack
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadStackFromFile(file *hcl.File, ds *decodedStack) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)
	for _, block := range content.Blocks {
		switch block.Type {
		case "component":
			content, _, contentDiags := block.Body.PartialContent(componentSchema)
			diags = append(diags, contentDiags...)

			if len(block.Labels) != 1 || block.Labels[0] == "" {
				continue
			}

			name := block.Labels[0]
			source := ""
			var versionCons version.Constraints

			var valDiags hcl.Diagnostics
			if attr, defined := content.Attributes["source"]; defined {
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &source)
				diags = append(diags, valDiags...)
			}

			if attr, defined := content.Attributes["version"]; defined {
				var versionStr string
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &versionStr)
				diags = append(diags, valDiags...)
				if versionStr != "" {
					vc, err := version.NewConstraint(versionStr)
					if err == nil {
						versionCons = vc
					}
				}
			}

			ds.Components[name] = &stack.Component{
				Source:     source,
				SourceAddr: module.ParseModuleSourceAddr(source),
				Version:    versionCons,
			}
		case "provider":
			// there is no point to parsing provider blocks here, as they need the full
			// context of the configuration to be parsed correctly
		case "required_providers":
			reqs, reqsDiags := decodeRequiredProvidersBlock(block)
			diags = append(diags, reqsDiags...)
			for name, req := range reqs {
				if _, exists := ds.ProviderRequirements[name]; !exists {
					ds.ProviderRequirements[name] = req
				} else {
					if req.Source != "" {
						source := ds.ProviderRequirements[name].Source
						if source != "" && source != req.Source {
							diags = append(diags, &hcl.Diagnostic{
								Severity: hcl.DiagError,
								Summary:  "Multiple provider source attributes",
								Detail:   fmt.Sprintf("Found multiple source attributes for provider %s: %q, %q", name, source, req.Source),
								Subject:  &block.DefRange,
							})
						} else {
							ds.ProviderRequirements[name].Source = req.Source
						}
					}

					if req.VersionConstraints != "" {
						existingVersionConstraints := ds.ProviderRequirements[name].VersionConstraints
						if existingVersionConstraints != "" && existingVersionConstraints != req.VersionConstraints {
							diags = append(diags, &hcl.Diagnostic{
								Severity: hcl.DiagError,
								Summary:  "Multiple provider version constraints",
								Detail:   fmt.Sprintf("Found multiple version constraints for provider %s: %q, %q", name, existingVersionConstraints, req.VersionConstraints),
								Subject:  &block.DefRange,
							})
						} else {
							ds.ProviderRequirements[name].VersionConstraints = req.VersionConstraints
						}
					}
				}
			}
		case "variable":
			// TODO: rename to input
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

			ds.Variables[name] = &stack.Variable{
				Type:         varType,
				Description:  description,
				DefaultValue: defaultValue,
				TypeDefaults: defaults,
				IsSensitive:  isSensitive,
			}
		case "output":
			content, _, contentDiags := block.Body.PartialContent(outputSchema)
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

			if attr, defined := content.Attributes["type"]; defined {
				valDiags = gohcl.DecodeExpression(attr.Expr, nil, &isSensitive)
				diags = append(diags, valDiags...)
			}

			value := cty.NilVal
			if attr, defined := content.Attributes["value"]; defined {
				// TODO: Provide context w/ funcs and variables
				val, diags := attr.Expr.Value(nil)
				if !diags.HasErrors() {
					value = val
				}
			}

			ds.Outputs[name] = &stack.Output{
				Description: description,
				IsSensitive: isSensitive,
				Value:       value,
			}
		}
	}

	return diags
}
