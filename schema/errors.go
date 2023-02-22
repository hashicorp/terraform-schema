// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

type coreSchemaRequiredErr struct{}

func (e coreSchemaRequiredErr) Error() string {
	return "core schema required (none provided)"
}

type NoCompatibleSchemaErr struct {
	Version     *version.Version
	Constraints version.Constraints
}

func (e NoCompatibleSchemaErr) Error() string {
	if e.Version != nil {
		return fmt.Sprintf("no compatible schema found for %s", e.Version)
	}
	if e.Constraints != nil && len(e.Constraints) > 0 {
		return fmt.Sprintf("no compatible schema found for %s", e.Constraints)
	}
	return "no compatible schema found"
}
