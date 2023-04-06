// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v014_mod "github.com/hashicorp/terraform-schema/internal/schema/0.14"
)

func ConnectionDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	bodies := v014_mod.ConnectionDependentBodies(v)

	ssh := schema.NewSchemaKey(schema.DependencyKeys{
		Attributes: []schema.AttributeDependent{
			{
				Name: "type",
				Expr: schema.ExpressionValue{Static: cty.StringVal("ssh")},
			},
		},
	})

	// See https://github.com/hashicorp/terraform/commit/5b99a56f
	bodies[ssh].Attributes["target_platform"] = &schema.AttributeSchema{
		Constraint: schema.OneOf{
			schema.LiteralValue{Value: cty.StringVal("windows")},
			schema.LiteralValue{Value: cty.StringVal("unix")},
		},
		IsOptional: true,
		Description: lang.Markdown("The target platform to connect to. " +
			"Defaults to `unix` if not set. If the platform is set to `windows`, the default `script_path`" +
			" is `" + `c:\windows\temp\terraform_%RAND%.cmd` + ", assuming the SSH default shell is `cmd.exe`."),
	}

	return bodies
}
