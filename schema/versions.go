// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import "github.com/hashicorp/go-version"

// ResolveVersion returns Terraform version for which we have schema available
// based on either given version and/or constraint.
// Lack of constraint and version implies latest known version.
//
//go:generate go run ../internal/versiongen -w ./versions_gen.go
func ResolveVersion(tfVersion *version.Version, tfCons version.Constraints) *version.Version {
	if tfVersion != nil {
		coreVersion := tfVersion.Core()
		if coreVersion.LessThan(OldestAvailableVersion) {
			return OldestAvailableVersion
		}
		if coreVersion.GreaterThan(LatestAvailableVersion) {
			// We could simply return coreVersion or tfVersion here
			// but we ensure the return version is actually known, i.e.
			// that we actually have schema for it.
			//
			// Also we strip the pre-release part as it simplifies
			// the version comparisons downstream (we don't care
			// about differences between individual pre-releases
			// of the same patch version).
			for _, v := range terraformVersions {
				if tfVersion.Equal(v) {
					return coreVersion
				}
			}
			return LatestAvailableVersion
		}
		if tfCons.Check(coreVersion) {
			return coreVersion
		}
	}

	for _, v := range terraformVersions {
		if len(tfCons) > 0 && tfCons.Check(v) && v.LessThan(OldestAvailableVersion) {
			return OldestAvailableVersion
		}
		if tfVersion != nil && tfVersion.Core().Equal(v) {
			return tfVersion.Core()
		}
		if len(tfCons) > 0 && tfCons.Check(v) {
			return v
		}
	}

	return LatestAvailableVersion
}
