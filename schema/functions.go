// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	funcs_v0_12 "github.com/hashicorp/terraform-schema/internal/funcs/0.12"
	funcs_v0_13 "github.com/hashicorp/terraform-schema/internal/funcs/0.13"
	funcs_v0_14 "github.com/hashicorp/terraform-schema/internal/funcs/0.14"
	funcs_v0_15 "github.com/hashicorp/terraform-schema/internal/funcs/0.15"
	funcs_v1_3 "github.com/hashicorp/terraform-schema/internal/funcs/1.3"
	funcs_generated "github.com/hashicorp/terraform-schema/internal/funcs/generated"
)

func FunctionsForVersion(v *version.Version) (map[string]schema.FunctionSignature, error) {
	ver := v.Core()
	if ver.GreaterThanOrEqual(v1_4) {
		return funcs_generated.Functions(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_3) {
		return funcs_v1_3.Functions(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_15) {
		return funcs_v0_15.Functions(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_14) {
		return funcs_v0_14.Functions(ver), nil
	}
	if ver.GreaterThanOrEqual(v0_13) {
		return funcs_v0_13.Functions(ver), nil
	}

	// Return the 0.12 functions for any version <= 0.12
	return funcs_v0_12.Functions(ver), nil
}

func FunctionsForConstraint(vc version.Constraints) (map[string]schema.FunctionSignature, error) {
	for _, v := range terraformVersions {
		if vc.Check(v) {
			return FunctionsForVersion(v)
		}
	}

	return nil, fmt.Errorf("no compatible functions found for %s", vc)
}
