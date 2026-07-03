// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchModuleBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["source"] = &schema.AttributeSchema{
		Constraint: schema.AnyExpression{OfType: cty.String},
		Description: lang.Markdown("Source where to load the module from, " +
			"a local directory (e.g. `./module`) or a remote address - e.g. " +
			"`hashicorp/consul/aws` (Terraform Registry address) or " +
			"`github.com/hashicorp/example` (GitHub)"),
		IsRequired:             true,
		IsDepKey:               true,
		SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
		CompletionHooks: lang.CompletionHooks{
			{Name: "CompleteLocalModuleSources"},
			{Name: "CompleteRegistryModuleSources"},
		},
	}

	bs.Body.Attributes["version"] = &schema.AttributeSchema{
		Constraint: schema.AnyExpression{OfType: cty.String},
		IsOptional: true,
		Description: lang.Markdown("Constraint to set the version of the module, e.g. `~> 1.0`." +
			" Only applicable to modules in a module registry."),
		CompletionHooks: lang.CompletionHooks{
			{Name: "CompleteRegistryModuleVersions"},
		},
	}

	return bs
}
