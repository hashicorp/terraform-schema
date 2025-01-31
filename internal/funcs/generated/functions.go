// Code generated by "gen"; DO NOT EDIT.
package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

var (
	v1_10_0 = version.Must(version.NewVersion("1.10.0"))
	v1_9_0  = version.Must(version.NewVersion("1.9.0"))
	v1_8_0  = version.Must(version.NewVersion("1.8.0"))
	v1_5_0  = version.Must(version.NewVersion("1.5.0"))
	v1_4_0  = version.Must(version.NewVersion("1.4.0"))
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	if v.GreaterThanOrEqual(v1_10_0) {
		return v1_10_0_Functions()
	}
	if v.GreaterThanOrEqual(v1_9_0) {
		return v1_9_0_Functions()
	}
	if v.GreaterThanOrEqual(v1_8_0) {
		return v1_8_0_Functions()
	}
	if v.GreaterThanOrEqual(v1_5_0) {
		return v1_5_0_Functions()
	}
	if v.GreaterThanOrEqual(v1_4_0) {
		return v1_4_0_Functions()
	}

	return v1_4_0_Functions()
}
