// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_2_mod "github.com/hashicorp/terraform-schema/internal/schema/1.2"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_2_mod.ModuleSchema(v)
	bs.Blocks["resource"].Body.Blocks["provisioner"].DependentBody = ProvisionerDependentBodies(v)

	return bs
}
