// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

var (
	v0_12_2  = version.Must(version.NewVersion("0.12.2"))
	v0_12_6  = version.Must(version.NewVersion("0.12.6"))
	v0_12_7  = version.Must(version.NewVersion("0.12.7"))
	v0_12_18 = version.Must(version.NewVersion("0.12.18"))
	v0_12_20 = version.Must(version.NewVersion("0.12.20"))
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"data":      datasourceBlockSchema(v),
			"locals":    localsBlockSchema(),
			"module":    moduleBlockSchema(),
			"output":    outputBlockSchema(),
			"provider":  providerBlockSchema(v),
			"resource":  resourceBlockSchema(v),
			"variable":  variableBlockSchema(v),
			"terraform": terraformBlockSchema(v),
		},
	}
}
