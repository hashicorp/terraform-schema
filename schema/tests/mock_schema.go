// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	test_v1_7 "github.com/hashicorp/terraform-schema/internal/schema/tests/1.7"
	tfschema "github.com/hashicorp/terraform-schema/schema"
)

// CoreMockSchemaForVersion finds a schema for mock configuration files
// that is relevant for the given Terraform version.
// It will return an error if such schema cannot be found.
func CoreMockSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()

	if ver.GreaterThanOrEqual(v1_7) {
		return test_v1_7.MockSchema(ver), nil
	}

	return nil, tfschema.NoCompatibleSchemaErr{Version: ver}
}
