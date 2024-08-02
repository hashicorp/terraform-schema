// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	v1_9_mod "github.com/hashicorp/terraform-schema/internal/schema/stacks/1.9"
)

// StackSchema returns the static schema for a stack
// configuration (*.tfstack.hcl) file.
func StackSchema(_ *version.Version) *schema.BodySchema {
	bs := v1_9_mod.StackSchema(nil)

	bs.Blocks["variable"].IsDeprecated = true
	bs.Blocks["variable"].Description = lang.Markdown("The `variables` attribute has been replaced with the `inputs` attribute. Please update the configuration to use `inputs` instead of `variables`, as support for the `variables` attribute will be removed entirely before the final release of Terraform Stacks.")

	bs.Blocks["input"] = inputBlockSchema()

	return bs
}

// DeploymentSchema returns the static schema for a deployment
// configuration (*.tfdeploy.hcl) file.
func DeploymentSchema(_ *version.Version) *schema.BodySchema {
	bs := v1_9_mod.DeploymentSchema(nil)

	return bs
}
