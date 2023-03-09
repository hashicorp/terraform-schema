// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/backend"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// decodedModule is the type representing a decoded Terraform module.
type decodedModule struct {
	RequiredCore         []string
	Backends             map[string]backend.BackendData
	CloudBackend         *backend.Cloud
	ProviderRequirements map[string]*providerRequirement
	ProviderConfigs      map[string]*providerConfig
	Resources            map[string]*resource
	DataSources          map[string]*dataSource
	Variables            map[string]*module.Variable
	Outputs              map[string]*module.Output
	ModuleCalls          map[string]*module.DeclaredModuleCall
}

func newDecodedModule() *decodedModule {
	return &decodedModule{
		RequiredCore:         make([]string, 0),
		Backends:             make(map[string]backend.BackendData),
		ProviderRequirements: make(map[string]*providerRequirement),
		ProviderConfigs:      make(map[string]*providerConfig),
		Resources:            make(map[string]*resource),
		DataSources:          make(map[string]*dataSource),
		Variables:            make(map[string]*module.Variable),
		Outputs:              make(map[string]*module.Output),
		ModuleCalls:          make(map[string]*module.DeclaredModuleCall),
	}
}

// providerConfig represents a provider block in the configuration
type providerConfig struct {
	Name  string
	Alias string
}

// loadModuleFromFile reads given file, interprets it and stores in given Module
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadModuleFromFile(file *hcl.File, mod *decodedModule) hcl.Diagnostics {
	var diags hcl.Diagnostics
	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "terraform":
			content, _, contentDiags := block.Body.PartialContent(terraformBlockSchema)
			diags = append(diags, contentDiags...)

			if attr, defined := content.Attributes["required_version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() {
					mod.RequiredCore = append(mod.RequiredCore, version)
				}
			}

			for _, innerBlock := range content.Blocks {
				switch innerBlock.Type {
				case "cloud":
					data, bDiags := decodeCloudBlock(innerBlock)
					diags = append(diags, bDiags...)
					mod.CloudBackend = data
				case "backend":
					bType := innerBlock.Labels[0]

					data, bDiags := decodeBackendsBlock(innerBlock)
					diags = append(diags, bDiags...)

					if _, exists := mod.Backends[bType]; exists {
						diags = append(diags, &hcl.Diagnostic{
							Severity: hcl.DiagError,
							Summary:  "Multiple backend definitions",
							Detail:   fmt.Sprintf("Found multiple backend definitions for %q. Only one is allowed.", bType),
							Subject:  &innerBlock.DefRange,
						})
						continue
					}

					mod.Backends[bType] = data
				case "required_providers":
					reqs, reqsDiags := decodeRequiredProvidersBlock(innerBlock)
					diags = append(diags, reqsDiags...)
					for name, req := range reqs {
						if _, exists := mod.ProviderRequirements[name]; !exists {
							mod.ProviderRequirements[name] = req
						} else {
							if req.Source != "" {
								source := mod.ProviderRequirements[name].Source
								if source != "" && source != req.Source {
									diags = append(diags, &hcl.Diagnostic{
										Severity: hcl.DiagError,
										Summary:  "Multiple provider source attributes",
										Detail:   fmt.Sprintf("Found multiple source attributes for provider %s: %q, %q", name, source, req.Source),
										Subject:  &innerBlock.DefRange,
									})
								} else {
									mod.ProviderRequirements[name].Source = req.Source
								}
							}

							mod.ProviderRequirements[name].VersionConstraints = append(mod.ProviderRequirements[name].VersionConstraints, req.VersionConstraints...)
						}
					}
				}
			}
		case "provider":
			content, _, contentDiags := block.Body.PartialContent(providerConfigSchema)
			diags = append(diags, contentDiags...)

			name := block.Labels[0]
			// Even if there isn't an explicit version required, we still
			// need an entry in our map to signal the unversioned dependency.
			if _, exists := mod.ProviderRequirements[name]; !exists {
				mod.ProviderRequirements[name] = &providerRequirement{}
			}
			if attr, defined := content.Attributes["version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() {
					mod.ProviderRequirements[name].VersionConstraints = append(mod.ProviderRequirements[name].VersionConstraints, version)
				}
			}

			providerKey := name
			var alias string
			if attr, defined := content.Attributes["alias"]; defined {
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &alias)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() && alias != "" {
					providerKey = fmt.Sprintf("%s.%s", name, alias)
				}
			}

			mod.ProviderConfigs[providerKey] = &providerConfig{
				Name:  name,
				Alias: alias,
			}

		case "data":
			content, _, contentDiags := block.Body.PartialContent(resourceSchema)
			diags = append(diags, contentDiags...)

			ds := &dataSource{
				Type: block.Labels[0],
				Name: block.Labels[1],
			}

			mod.DataSources[ds.MapKey()] = ds

			if attr, defined := content.Attributes["provider"]; defined {
				ref, aDiags := decodeProviderAttribute(attr)
				diags = append(diags, aDiags...)
				ds.Provider = ref
			} else {
				// If provider _isn't_ set then we'll infer it from the
				// datasource type.
				ds.Provider = module.ProviderRef{
					LocalName: inferProviderNameFromType(ds.Type),
				}
			}

		case "resource":
			content, _, contentDiags := block.Body.PartialContent(resourceSchema)
			diags = append(diags, contentDiags...)

			r := &resource{
				Type: block.Labels[0],
				Name: block.Labels[1],
			}

			mod.Resources[r.MapKey()] = r

			if attr, defined := content.Attributes["provider"]; defined {
				ref, aDiags := decodeProviderAttribute(attr)
				diags = append(diags, aDiags...)
				r.Provider = ref
			} else {
				// If provider _isn't_ set then we'll infer it from the
				// resource type.
				r.Provider = module.ProviderRef{
					LocalName: inferProviderNameFromType(r.Type),
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
			mod.Variables[name] = &module.Variable{
				Type:         varType,
				Description:  description,
				IsSensitive:  isSensitive,
				DefaultValue: defaultValue,
				TypeDefaults: defaults,
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
			if attr, defined := content.Attributes["sensitive"]; defined {
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
			mod.Outputs[name] = &module.Output{
				Description: description,
				IsSensitive: isSensitive,
				Value:       value,
			}
		case "module":
			content, remainingBody, contentDiags := block.Body.PartialContent(moduleSchema)
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

			inputNames := make([]string, 0)
			remainingAttributes, diags := remainingBody.JustAttributes()
			if !diags.HasErrors() {
				for name := range remainingAttributes {
					inputNames = append(inputNames, name)
				}
			}

			sort.Strings(inputNames)

			var rng *hcl.Range
			hclBody, ok := block.Body.(*hclsyntax.Body)
			if ok {
				rng = hclBody.Range().Ptr()
			}

			mod.ModuleCalls[name] = &module.DeclaredModuleCall{
				LocalName:  name,
				SourceAddr: module.ParseModuleSourceAddr(source),
				Version:    versionCons,
				InputNames: inputNames,
				RangePtr:   rng,
			}
		}

	}

	return diags
}

func decodeProviderAttribute(attr *hcl.Attribute) (module.ProviderRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	// New style here is to provide this as a naked traversal
	// expression, but we also support quoted references for
	// older configurations that predated this convention.
	traversal, travDiags := hcl.AbsTraversalForExpr(attr.Expr)
	if travDiags.HasErrors() {
		traversal = nil // in case we got any partial results

		// Fall back on trying to parse as a string
		var travStr string
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &travStr)
		if !valDiags.HasErrors() {
			var strDiags hcl.Diagnostics
			traversal, strDiags = hclsyntax.ParseTraversalAbs([]byte(travStr), "", hcl.Pos{})
			if strDiags.HasErrors() {
				traversal = nil
			}
		}
	}

	// If we get out here with a nil traversal then we didn't
	// succeed in processing the input.
	if len(traversal) > 0 {
		providerName := traversal.RootName()
		alias := ""
		if len(traversal) > 1 {
			if getAttr, ok := traversal[1].(hcl.TraverseAttr); ok {
				alias = getAttr.Name
			}
		}
		return module.ProviderRef{
			LocalName: providerName,
			Alias:     alias,
		}, diags
	}

	return module.ProviderRef{}, hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider reference",
			Detail:   "Provider argument requires a provider name followed by an optional alias, like \"aws.foo\".",
			Subject:  attr.Expr.Range().Ptr(),
		},
	}
}
