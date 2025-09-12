// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_12_mod "github.com/hashicorp/terraform-schema/internal/schema/1.12"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_12_mod.ModuleSchema(v)

	bs.Blocks["action"] = actionBlockSchema()
	bs.Blocks["resource"].Body.Blocks["lifecycle"].Body.Blocks["action_trigger"] = resourceLifecycleActionTriggerBlock()

	return bs
}
