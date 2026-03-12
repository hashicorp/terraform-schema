// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import "github.com/hashicorp/hcl-lang/schema"

func patchMockDataBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	return bs
}
