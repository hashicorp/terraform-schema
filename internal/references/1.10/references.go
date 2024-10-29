// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package references

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	refs_v0_12 "github.com/hashicorp/terraform-schema/internal/references/0.12"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func BuiltinReferences(modPath string) reference.Targets {
	refs := refs_v0_12.BuiltinReferences(modPath)

	refs = append(refs, reference.Target{
		Addr: lang.Address{
			lang.RootStep{Name: "terraform"},
			lang.AttrStep{Name: "applying"},
		},
		ScopeId:     refscope.BuiltinScope,
		Type:        cty.Bool,
		Description: lang.Markdown("True if Terraform is currently in the apply phase (including destroy mode), false otherwise"),
	})

	return refs
}
