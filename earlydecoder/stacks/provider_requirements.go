// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type providerRequirement struct {
	Source             string
	VersionConstraints []string
}

func decodeRequiredProvidersBlock(block *hcl.Block) (map[string]*providerRequirement, hcl.Diagnostics) {
	attrs, diags := block.Body.JustAttributes()
	reqs := make(map[string]*providerRequirement)
	for name, attr := range attrs {
		kvs, mapDiags := hcl.ExprMap(attr.Expr)
		if mapDiags.HasErrors() {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid required_providers object",
				Detail:   "Required providers entries must be objects for stacks configuration files.",
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
					// TODO: simplify this to only allow a single version constraint string
					// if stacks only allows a single required_providers block
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
			}

			reqs[name] = &pr
		}
	}

	return reqs, diags
}
