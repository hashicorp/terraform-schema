package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	funcs_v0_12 "github.com/hashicorp/terraform-schema/internal/funcs/0.12"
	funcs_v0_13 "github.com/hashicorp/terraform-schema/internal/funcs/0.13"
	funcs_v0_14 "github.com/hashicorp/terraform-schema/internal/funcs/0.14"
)

func FunctionsForVersion(v *version.Version) (map[string]schema.FunctionSignature, error) {
	ver, err := semVer(v)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
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

	return nil, NoCompatibleFunctionsErr{Constraints: vc}
}
