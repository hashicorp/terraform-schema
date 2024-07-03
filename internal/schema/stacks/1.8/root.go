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
			"component":         componentBlockSchema(),
			"provider":          providerBlockSchema(),
			"required_provider": providerBlockSchema(),
			"variable":          variableBlockSchema(),
			"output":            outputBlockSchema(),
		},
	}
}

func DeploymentSchema(v *version.Version) *schema.BodySchema {
	// TODO: This will likely change after the refactor
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment":     deploymentBlockSchema(),
			"identity_token": identityTokenBlockSchema(),
			"orchestrate":    orchestrateBlockSchema(),
		},
	}
}
