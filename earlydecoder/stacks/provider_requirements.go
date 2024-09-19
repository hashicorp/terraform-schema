// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/go-version"
)

type providerRequirement struct {
	Source             string
	VersionConstraints string
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
				parsedVersion, valDiags := kv.Value.Value(nil)
				if valDiags.HasErrors() || !parsedVersion.Type().Equals(cty.String) {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Unsuitable value type",
						Detail:   "Unsuitable value: string required",
						Subject:  attr.Expr.Range().Ptr(),
					})
					continue
				}
				if !parsedVersion.IsNull() {
					pr.VersionConstraints = parsedVersion.AsString()
					_, err := version.NewConstraint(pr.VersionConstraints)
					if err != nil {
						// TODO
						diags = append(diags, &hcl.Diagnostic{
							Severity: hcl.DiagError,
							Summary:  fmt.Sprintf("Unable to parse %q provider requirements", name),
							Detail:   fmt.Sprintf("Constraint %q is not a valid constraint: %s", pr.VersionConstraints, err),
							Subject: attr.Expr.Range().Ptr(),
						})
						continue
					}
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

					_, err := tfaddr.ParseProviderSource(pr.Source)
					if err != nil {
						diags = append(diags, &hcl.Diagnostic{
							Severity: hcl.DiagError,
							Summary:  fmt.Sprintf("Unable to parse provider source for %q", name),
							Detail:   fmt.Sprintf("%q provider source (%q) is not a valid source string", name, pr.Source),
							Subject:  attr.Expr.Range().Ptr(),
						})
						// continue
					}
				}
			}

			reqs[name] = &pr
		}
	}

	return reqs, diags
}
