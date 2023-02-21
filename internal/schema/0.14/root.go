// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v013_mod "github.com/hashicorp/terraform-schema/internal/schema/0.13"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v013_mod.ModuleSchema(v)

	bs.Blocks["variable"] = variableBlockSchema
	bs.Blocks["terraform"] = terraformBlockSchema(v)
	bs.Blocks["resource"].Body.Blocks["provisioner"].DependentBody = ProvisionerDependentBodies

	return bs
}
