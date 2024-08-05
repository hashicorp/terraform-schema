// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	test_v1_6 "github.com/hashicorp/terraform-schema/internal/schema/tests/1.6"
	test_v1_7 "github.com/hashicorp/terraform-schema/internal/schema/tests/1.7"
	test_v1_9 "github.com/hashicorp/terraform-schema/internal/schema/tests/1.9"
	tfschema "github.com/hashicorp/terraform-schema/schema"
)

var (
	v1_6 = version.Must(version.NewVersion("1.6"))
	v1_7 = version.Must(version.NewVersion("1.7"))
	v1_9 = version.Must(version.NewVersion("1.9"))
)

// CoreTestSchemaForVersion finds a schema for test configuration files
// that is relevant for the given Terraform version.
// It will return an error if such schema cannot be found.
func CoreTestSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()

	if ver.GreaterThanOrEqual(v1_9) {
		return test_v1_9.TestSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_7) {
		return test_v1_7.TestSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_6) {
		return test_v1_6.TestSchema(ver), nil
	}

	return nil, tfschema.NoCompatibleSchemaErr{Version: ver}
}
