// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

type providerRequirement struct {
	Source               string
	VersionConstraints   []string
	ConfigurationAliases []module.ProviderRef
}

func decodeRequiredProvidersBlock(block *hcl.Block) (map[string]*providerRequirement, hcl.Diagnostics) {
	attrs, diags := block.Body.JustAttributes()
	reqs := make(map[string]*providerRequirement)
	for name, attr := range attrs {
		// Look for a legacy version in the attribute first
		if expr, err := attr.Expr.Value(nil); err == nil && expr.Type().IsPrimitiveType() {
			var version string
			valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
			diags = append(diags, valDiags...)
			if !valDiags.HasErrors() {
				reqs[name] = &providerRequirement{
					VersionConstraints: []string{version},
				}
			}
			continue
		}

		kvs, mapDiags := hcl.ExprMap(attr.Expr)
		if mapDiags.HasErrors() {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid required_providers object",
				Detail:   "Required providers entries must be strings or objects.",
				Subject:  attr.Expr.Range().Ptr(),
			})
			continue
		}

		var pr providerRequirement

		for _, kv := range kvs {
			key, keyDiags := kv.Key.Value(nil)
			if keyDiags.HasErrors() {
				diags = append(diags, keyDiags...)
				continue
			}

			if key.Type() != cty.String {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid Attribute",
					Detail:   fmt.Sprintf("Invalid attribute value for provider requirement: %#v", key),
					Subject:  kv.Key.Range().Ptr(),
				})
				continue
			}

			switch key.AsString() {
			case "version":
				version, valDiags := kv.Value.Value(nil)
				if valDiags.HasErrors() || !version.Type().Equals(cty.String) {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Unsuitable value type",
						Detail:   "Unsuitable value: string required",
						Subject:  attr.Expr.Range().Ptr(),
					})
					continue
				}
				if !version.IsNull() {
					pr.VersionConstraints = append(pr.VersionConstraints, version.AsString())
				}

			case "source":
				source, valDiags := kv.Value.Value(nil)
				if valDiags.HasErrors() || !source.Type().Equals(cty.String) {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Unsuitable value type",
						Detail:   "Unsuitable value: string required",
						Subject:  attr.Expr.Range().Ptr(),
					})
					continue
				}

				if !source.IsNull() {
					pr.Source = source.AsString()
				}
			case "configuration_aliases":
				aliases, valDiags := decodeConfigurationAliases(name, kv.Value)
				if valDiags.HasErrors() {
					diags = append(diags, valDiags...)
					continue
				}
				pr.ConfigurationAliases = append(pr.ConfigurationAliases, aliases...)
			}

			reqs[name] = &pr
		}
	}

	return reqs, diags
}

func decodeConfigurationAliases(localName string, value hcl.Expression) ([]module.ProviderRef, hcl.Diagnostics) {
	aliases := make([]module.ProviderRef, 0)
	var diags hcl.Diagnostics

	exprs, listDiags := hcl.ExprList(value)
	if listDiags.HasErrors() {
		diags = append(diags, listDiags...)
		return aliases, diags
	}

	for _, expr := range exprs {
		traversal, travDiags := hcl.AbsTraversalForExpr(expr)
		if travDiags.HasErrors() {
			diags = append(diags, travDiags...)
			continue
		}

		ref, cfgDiags := parseProviderRef(traversal)
		if cfgDiags.HasErrors() {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid configuration_aliases value",
				Detail:   `Configuration aliases can only contain references to local provider configuration names in the format of provider.alias`,
				Subject:  value.Range().Ptr(),
			})
			continue
		}

		if ref.LocalName != localName {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid configuration_aliases value",
				Detail:   fmt.Sprintf(`Configuration aliases must be prefixed with the provider name. Expected %q, but found %q.`, localName, ref.LocalName),
				Subject:  value.Range().Ptr(),
			})
			continue
		}

		aliases = append(aliases, ref)
	}

	return aliases, diags
}

func parseProviderRef(traversal hcl.Traversal) (module.ProviderRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	ret := module.ProviderRef{
		LocalName: traversal.RootName(),
	}

	if len(traversal) < 2 {
		// Just a local name, then.
		return ret, diags
	}

	aliasStep := traversal[1]
	switch ts := aliasStep.(type) {
	case hcl.TraverseAttr:
		ret.Alias = ts.Name
		return ret, diags
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration address",
			Detail:   "The provider type name must either stand alone or be followed by an alias name separated with a dot.",
			Subject:  aliasStep.SourceRange().Ptr(),
		})
	}

	if len(traversal) > 2 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration address",
			Detail:   "Extraneous extra operators after provider configuration address.",
			Subject:  traversal[2:].SourceRange().Ptr(),
		})
	}

	return ret, diags
}
