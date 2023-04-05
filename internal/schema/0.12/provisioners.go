// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

// See https://github.com/hashicorp/terraform/blob/v0.12.0/command/internal_plugin_list.go

func ProvisionerDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	m := map[schema.SchemaKey]*schema.BodySchema{
		labelKey("file"):        FileProvisioner,
		labelKey("local-exec"):  LocalExecProvisioner,
		labelKey("remote-exec"): RemoteExecProvisioner,
	}

	// Vendor provisioners are deprecated in 0.13.4+
	// See https://discuss.hashicorp.com/t/notice-terraform-to-begin-deprecation-of-vendor-tool-specific-provisioners-starting-in-terraform-0-13-4/13997
	// Some of these provisioners have complex schemas
	// but we can at least helpfully list their names
	m[labelKey("chef")] = &schema.BodySchema{}
	m[labelKey("salt-masterless")] = &schema.BodySchema{}
	m[labelKey("habitat")] = &schema.BodySchema{}
	if v.GreaterThanOrEqual(v0_12_2) {
		// See https://github.com/hashicorp/terraform/commit/615110e13
		m[labelKey("puppet")] = &schema.BodySchema{}
	}

	return m
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}

var FileProvisioner = &schema.BodySchema{
	Description: lang.Markdown("Provisioner used to copy files or directories from the machine executing Terraform" +
		" to the newly created resource."),
	HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/file.html",
	Attributes: map[string]*schema.AttributeSchema{
		"source": {
			IsOptional: true,
			Constraint: schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("The source file or folder. It can be specified as relative " +
				"to the current working directory or as an absolute path. This attribute cannot be " +
				"specified with `content`."),
		},
		"content": {
			IsOptional: true,
			Constraint: schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("The content to copy on the destination. If destination is a file," +
				" the content will be written on that file, in case of a directory a file named `tf-file-content`" +
				" is created. It's recommended to use a file as the destination. This attribute cannot be " +
				"specified with `source`."),
		},
		"destination": {
			IsRequired:  true,
			Constraint:  schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("The destination path. It must be specified as an absolute path."),
		},
	},
}

var LocalExecProvisioner = &schema.BodySchema{
	Description: lang.Markdown("Invokes a local executable after a resource is created. " +
		"This invokes a process on the machine running Terraform, not on the resource."),
	HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/local-exec.html",
	Attributes: map[string]*schema.AttributeSchema{
		"command": {
			IsRequired: true,
			Constraint: schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("This is the command to execute. It can be provided as a relative path " +
				"to the current working directory or as an absolute path. It is evaluated in a shell, " +
				"and can use environment variables or Terraform variables."),
		},
		"interpreter": {
			IsOptional: true,
			Constraint: schema.List{
				Elem: schema.AnyExpression{OfType: cty.String},
			},
			Description: lang.Markdown("If provided, this is a list of interpreter arguments used to execute " +
				"the command. The first argument is the interpreter itself. It can be provided as a relative " +
				"path to the current working directory or as an absolute path. The remaining arguments are " +
				"appended prior to the command. This allows building command lines of the form " +
				"`\"/bin/bash\", \"-c\", \"echo foo\"`. If interpreter is unspecified, sensible defaults " +
				"will be chosen based on the system OS."),
		},
		"working_dir": {
			IsOptional: true,
			Constraint: schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("If provided, specifies the working directory where command will be executed. " +
				"It can be provided as as a relative path to the current working directory or as an absolute path. " +
				"The directory must exist."),
		},
		"environment": {
			IsOptional: true,
			Constraint: schema.Map{
				Elem: schema.AnyExpression{OfType: cty.String},
			},
			Description: lang.Markdown("Map of key value pairs representing the environment of the executed command. " +
				"Inherits the current process environment."),
		},
	},
}

var RemoteExecProvisioner = &schema.BodySchema{
	Description: lang.Markdown("Invokes a script on a remote resource after it is created. " +
		"This can be used to run a configuration management tool, bootstrap into a cluster, etc."),
	HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/remote-exec.html",
	Attributes: map[string]*schema.AttributeSchema{
		"inline": {
			IsOptional: true,
			Constraint: schema.List{
				Elem: schema.AnyExpression{OfType: cty.String},
			},
			Description: lang.Markdown("A list of command strings. They are executed in the order they are provided." +
				" This cannot be provided with `script` or `scripts`."),
		},
		"script": {
			IsOptional: true,
			Constraint: schema.AnyExpression{OfType: cty.String},
			Description: lang.Markdown("A path (relative or absolute) to a local script that will be copied " +
				"to the remote resource and then executed. This cannot be provided with `inline` or `scripts`."),
		},
		"scripts": {
			IsOptional: true,
			Constraint: schema.List{
				Elem: schema.AnyExpression{OfType: cty.String},
			},
			Description: lang.Markdown("A list of paths (relative or absolute) to local scripts that will be copied " +
				"to the remote resource and then executed. They are executed in the order they are provided." +
				" This cannot be provided with `inline` or `script`."),
		},
	},
}
