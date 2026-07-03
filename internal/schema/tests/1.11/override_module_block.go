// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/schema"
)

func patchOverrideModuleBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	return bs
}
