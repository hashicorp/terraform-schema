// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

func StackSchema(v *version.Version) *schema.BodySchema {
	// TODO: This will likely change after the refactor
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"component":         componentBlockSchema(v),
			"provider":          providerBlockSchema(v),
			"required_provider": providerBlockSchema(v),
			"variable":          variableBlockSchema(v),
			"output":            outputBlockSchema(v),
		},
	}
}

func DeploymentSchema(v *version.Version) *schema.BodySchema {
	// TODO: This will likely change after the refactor
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			// tfdeploy.hcl
			"deployment":     deploymentBlockSchema(v),
			"identity_token": identityTokenBlockSchema(v),
			"orchestrate":    orchestrateBlockSchema(v),
		},
	}
}
