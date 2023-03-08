// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/backend"
	"github.com/zclconf/go-cty/cty"
)

func decodeBackendsBlock(block *hcl.Block) (backend.BackendData, hcl.Diagnostics) {
	bType := block.Labels[0]
	attrs, diags := block.Body.JustAttributes()

	switch bType {
	case "remote":
		if attr, ok := attrs["hostname"]; ok {
			val, vDiags := attr.Expr.Value(nil)
			diags = append(diags, vDiags...)
			if val.IsWhollyKnown() && val.Type() == cty.String {
				return &backend.Remote{
					Hostname: val.AsString(),
				}, nil
			}
		}

		return &backend.Remote{}, nil
	}

	return &backend.UnknownBackendData{}, diags
}

func decodeCloudBlock(block *hcl.Block) (*backend.Cloud, hcl.Diagnostics) {
	attrs, diags := block.Body.JustAttributes()

	// https://developer.hashicorp.com/terraform/language/settings/terraform-cloud#usage-example
	// Required for Terraform Enterprise
	// Defaults to app.terraform.io for Terraform Cloud
	if attr, ok := attrs["hostname"]; ok {
		val, vDiags := attr.Expr.Value(nil)
		diags = append(diags, vDiags...)
		if val.IsWhollyKnown() && val.Type() == cty.String {
			return &backend.Cloud{
				Hostname: val.AsString(),
			}, nil
		}
	}

	// since it defaults to app.terraform.io, it is safe to return that
	// if hostname is empty
	return &backend.Cloud{
		Hostname: "app.terraform.io",
	}, nil
}
