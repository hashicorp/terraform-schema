// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func connectionBlock(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		Description:            lang.Markdown("Connection block describing how the provisioner connects to the given instance"),
		MaxItems:               1,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Connection},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				SelfRefs: true,
			},
			HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/connection.html",
			Attributes: map[string]*schema.AttributeSchema{
				"type": {
					Constraint: schema.OneOf{
						schema.LiteralValue{
							Value:       cty.StringVal("ssh"),
							Description: lang.Markdown("Use SSH to connect and provision the instance"),
						},
						schema.LiteralValue{
							Value:       cty.StringVal("winrm"),
							Description: lang.Markdown("Use WinRM to connect and provision the instance"),
						},
					},
					IsOptional:             true,
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					Description:            lang.Markdown("Connection type to use - `ssh` (default) or `winrm`"),
				},
				"user": {
					Constraint: schema.AnyExpression{OfType: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The user that we should use for the connection. " +
						"Defaults to `root` when using type `ssh` and defaults to `Administrator` " +
						"when using type `winrm`."),
				},
				"password": {
					Constraint: schema.AnyExpression{OfType: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The password we should use for the connection. " +
						"In some cases this is specified by the provider."),
				},
				"host": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("The address of the resource to connect to"),
				},
				"port": {
					Constraint: schema.AnyExpression{OfType: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The port to connect to. Defaults to `22` " +
						"when using type `ssh` and defaults to `5985` when using type `winrm`."),
				},
				"timeout": {
					Constraint: schema.AnyExpression{OfType: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The timeout to wait for the connection to become " +
						"available. Should be provided as a string like `30s` or `5m`. " +
						"Defaults to 5 minutes."),
				},
				"script_path": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("The path used to copy scripts meant for remote execution."),
				},
			},
		},
		DependentBody: ConnectionDependentBodies(v),
	}
}

func ConnectionDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	m := make(map[schema.SchemaKey]*schema.BodySchema, 0)

	ssh := schema.NewSchemaKey(schema.DependencyKeys{
		Attributes: []schema.AttributeDependent{
			{
				Name: "type",
				Expr: schema.ExpressionValue{Static: cty.StringVal("ssh")},
			},
		},
	})

	m[ssh] = &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"private_key": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The contents of an SSH key to use for the connection. " +
					"This takes preference over the password if provided."),
			},
			"certificate": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The contents of a signed CA Certificate. The argument " +
					"must be used in conjunction with a `private_key`."),
			},
			"agent": {
				Constraint: schema.AnyExpression{OfType: cty.Bool},
				IsOptional: true,
				Description: lang.Markdown("Set to `false` to disable using `ssh-agent` to authenticate. " +
					"On Windows the only supported SSH authentication agent is " +
					"[Pageant](http://the.earth.li/~sgtatham/putty/0.66/htmldoc/Chapter9.html#pageant)."),
			},
			"agent_identity": {
				Constraint:  schema.AnyExpression{OfType: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The preferred identity from the ssh agent for authentication."),
			},
			"host_key": {
				Constraint:  schema.AnyExpression{OfType: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The public key from the remote host or the signing CA, used to verify the connection."),
			},
			"bastion_host": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("Setting this enables the bastion host connection. " +
					"This host will be connected to first, and then the `host` connection will be made from there."),
			},
			"bastion_host_key": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The public key from the remote host or the signing CA, " +
					"used to verify the host connection."),
			},
			"bastion_port": {
				Constraint: schema.AnyExpression{OfType: cty.Number},
				IsOptional: true,
				Description: lang.Markdown("The port to use connect to the bastion host. " +
					"Defaults to the value of the `port` field."),
			},
			"bastion_user": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The user for the connection to the bastion host. " +
					"Defaults to the value of the `user` field."),
			},
			"bastion_password": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The password we should use for the bastion host. " +
					"Defaults to the value of the `password` field."),
			},
			"bastion_private_key": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The contents of an SSH key file to use for the bastion host. " +
					"Defaults to the value of the `private_key` field."),
			},
		},
	}

	// See https://github.com/hashicorp/terraform/commit/3031aca9
	if v.GreaterThanOrEqual(v0_12_7) {
		m[ssh].Attributes["bastion_certificate"] = &schema.AttributeSchema{
			Constraint: schema.AnyExpression{OfType: cty.String},
			IsOptional: true,
			Description: lang.Markdown("The contents of a signed CA Certificate. The `certificate` argument " +
				"must be used in conjunction with a `bastion_private_key`."),
		}
	}

	winRm := schema.NewSchemaKey(schema.DependencyKeys{
		Attributes: []schema.AttributeDependent{
			{
				Name: "type",
				Expr: schema.ExpressionValue{Static: cty.StringVal("winrm")},
			},
		},
	})
	m[winRm] = &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"https": {
				Constraint:  schema.AnyExpression{OfType: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Set to `true` to connect using HTTPS instead of HTTP."),
			},
			"insecure": {
				Constraint:  schema.AnyExpression{OfType: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Set to `true` to not validate the HTTPS certificate chain."),
			},
			"use_ntlm": {
				Constraint: schema.AnyExpression{OfType: cty.Bool},
				IsOptional: true,
				Description: lang.Markdown("Set to `true` to use NTLM authentication, rather than default " +
					"(basic authentication), removing the requirement for basic authentication to be enabled " +
					"within the target guest. Read more about remote connection authentication at " +
					"[docs.microsoft.com](https://docs.microsoft.com/en-gb/windows/win32/winrm/authentication-for-remote-connections)."),
			},
			"cacert": {
				Constraint:  schema.AnyExpression{OfType: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The CA certificate to validate against."),
			},
		},
	}

	return m
}
