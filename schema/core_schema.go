// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	mod_v0_12 "github.com/hashicorp/terraform-schema/internal/schema/0.12"
	mod_v0_13 "github.com/hashicorp/terraform-schema/internal/schema/0.13"
	mod_v0_14 "github.com/hashicorp/terraform-schema/internal/schema/0.14"
	mod_v0_15 "github.com/hashicorp/terraform-schema/internal/schema/0.15"
	mod_v1_1 "github.com/hashicorp/terraform-schema/internal/schema/1.1"
	mod_v1_2 "github.com/hashicorp/terraform-schema/internal/schema/1.2"
	mod_v1_4 "github.com/hashicorp/terraform-schema/internal/schema/1.4"
	mod_v1_5 "github.com/hashicorp/terraform-schema/internal/schema/1.5"
)

var (
	v0_12 = version.Must(version.NewVersion("0.12"))
	v0_13 = version.Must(version.NewVersion("0.13"))
	v0_14 = version.Must(version.NewVersion("0.14"))
	v0_15 = version.Must(version.NewVersion("0.15"))
	v1_1  = version.Must(version.NewVersion("1.1"))
	v1_2  = version.Must(version.NewVersion("1.2"))
	v1_3  = version.Must(version.NewVersion("1.3"))
	v1_4  = version.Must(version.NewVersion("1.4"))
	v1_5  = version.Must(version.NewVersion("1.5"))
)

// CoreModuleSchemaForVersion finds a module schema which is relevant
// for the given Terraform version.
// It will return error if such schema cannot be found.
func CoreModuleSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()
	if ver.GreaterThanOrEqual(v1_5) {
		return mod_v1_5.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_4) {
		return mod_v1_4.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_2) {
		return mod_v1_2.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_1) {
		return mod_v1_1.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_15) {
		return mod_v0_15.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_14) {
		return mod_v0_14.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_13) {
		return mod_v0_13.ModuleSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_12) {
		return mod_v0_12.ModuleSchema(ver), nil
	}

	return nil, NoCompatibleSchemaErr{Version: ver}
}

// CoreModuleSchemaForConstraint returns schema relevant to the given constraint.
//
// Since the underlying implementation relies on knowing exact version and
// that is the "happy path" when exact version is known, we pre-generate
// known available versions, to be able to convert from constraints to
// versions easily. This means that e.g. `~> 1` translates to latest
// known version.
//
//go:generate go run ../internal/versiongen -w ./versions_gen.go
func CoreModuleSchemaForConstraint(vc version.Constraints) (*schema.BodySchema, error) {
	for _, v := range terraformVersions {
		if vc.Check(v) {
			return CoreModuleSchemaForVersion(v)
		}
	}

	return nil, NoCompatibleSchemaErr{Constraints: vc}
}
