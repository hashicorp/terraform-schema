// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import "github.com/hashicorp/hcl-lang/schema"

func patchMockProviderBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["override_during"] = overrideDuringAttributeSchema()

	bs.Body.Blocks["mock_resource"] = patchMockResourceBlockSchema(bs.Body.Blocks["mock_resource"])
	bs.Body.Blocks["mock_data"] = patchMockDataBlockSchema(bs.Body.Blocks["mock_data"])
	bs.Body.Blocks["override_resource"].Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	bs.Body.Blocks["override_data"].Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	return bs
}
