// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package tokmod

import (
	"github.com/hashicorp/hcl-lang/lang"
)

var (
	Data              = lang.SemanticTokenModifier("terraform-data")
	Locals            = lang.SemanticTokenModifier("terraform-locals")
	Module            = lang.SemanticTokenModifier("terraform-module")
	Output            = lang.SemanticTokenModifier("terraform-output")
	Provider          = lang.SemanticTokenModifier("terraform-provider")
	Resource          = lang.SemanticTokenModifier("terraform-resource")
	Ephemeral         = lang.SemanticTokenModifier("terraform-ephemeral")
	Provisioner       = lang.SemanticTokenModifier("terraform-provisioner")
	Connection        = lang.SemanticTokenModifier("terraform-connection")
	Variable          = lang.SemanticTokenModifier("terraform-variable")
	Terraform         = lang.SemanticTokenModifier("terraform-terraform")
	Backend           = lang.SemanticTokenModifier("terraform-backend")
	Name              = lang.SemanticTokenModifier("terraform-name")
	Type              = lang.SemanticTokenModifier("terraform-type")
	RequiredProviders = lang.SemanticTokenModifier("terraform-requiredProviders")
	Run               = lang.SemanticTokenModifier("terraform-run")
	Variables         = lang.SemanticTokenModifier("terraform-variables")
)

var SupportedModifiers = []lang.SemanticTokenModifier{
	Backend,
	Connection,
	Data,
	Locals,
	Module,
	Name,
	Output,
	Provider,
	Provisioner,
	RequiredProviders,
	Resource,
	Run,
	Terraform,
	Type,
	Variable,
	Variables,
}
