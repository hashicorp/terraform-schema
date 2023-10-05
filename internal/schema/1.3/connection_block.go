// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v015_mod "github.com/hashicorp/terraform-schema/internal/schema/0.15"
)

func ConnectionDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	bodies := v015_mod.ConnectionDependentBodies(v)

	ssh := schema.NewSchemaKey(schema.DependencyKeys{
		Attributes: []schema.AttributeDependent{
			{
				Name: "type",
				Expr: schema.ExpressionValue{Static: cty.StringVal("ssh")},
			},
		},
	})

	// See https://github.com/hashicorp/terraform/commit/4cfb6bc8
	bodies[ssh].Attributes["proxy_scheme"] = &schema.AttributeSchema{
		Constraint: schema.OneOf{
			schema.LiteralValue{Value: cty.StringVal("http")},
			schema.LiteralValue{Value: cty.StringVal("https")},
		},
		IsOptional: true,
		Description: lang.Markdown("Scheme to use to connect to the proxy (`http` or `https`). " +
			"Defaults to `http`."),
	}
	bodies[ssh].Attributes["proxy_host"] = &schema.AttributeSchema{
		Constraint: schema.AnyExpression{OfType: cty.String},
		IsOptional: true,
		Description: lang.Markdown("Host to connect to in order to enable SSH over HTTP connection. " +
			"This host will be connected to first, and then the `host` or `bastion_host` connection " +
			"will be made from there."),
	}
	bodies[ssh].Attributes["proxy_port"] = &schema.AttributeSchema{
		Constraint:  schema.AnyExpression{OfType: cty.Number},
		IsOptional:  true,
		Description: lang.Markdown("The port to use connect to the `proxy_host`"),
	}
	bodies[ssh].Attributes["proxy_user_name"] = &schema.AttributeSchema{
		Constraint:  schema.AnyExpression{OfType: cty.String},
		IsOptional:  true,
		Description: lang.Markdown("The username to use to connect to the `proxy_host`"),
	}
	bodies[ssh].Attributes["proxy_user_password"] = &schema.AttributeSchema{
		Constraint:  schema.AnyExpression{OfType: cty.String},
		IsOptional:  true,
		Description: lang.Markdown("The password to use to connect to the `proxy_host`"),
	}

	return bodies
}
