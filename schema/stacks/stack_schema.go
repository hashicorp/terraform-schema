// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	stack_1_9 "github.com/hashicorp/terraform-schema/internal/schema/stacks/1.9"
)

// CoreStackSchemaForVersion finds a schema for stack configuration files
// that is relevant for the given Terraform version.
// It will return an error if such schema cannot be found.
func CoreStackSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	return stack_1_9.StackSchema(v), nil
}
