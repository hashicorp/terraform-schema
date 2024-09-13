// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

// StackSchema returns the static schema for a stack
// configuration (*.tfstack.hcl) file.
func StackSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"component":          componentBlockSchema(),
			"provider":           providerBlockSchema(),
			"required_providers": requiredProvidersBlockSchema(),
			"variable":           variableBlockSchema(),
			"output":             outputBlockSchema(),
			"locals":             localsBlockSchema(),
			"removed":            removedBlockSchema(),
		},
	}
}

// DeploymentSchema returns the static schema for a deployment
// configuration (*.tfdeploy.hcl) file.
func DeploymentSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment":     deploymentBlockSchema(),
			"identity_token": identityTokenBlockSchema(),
			"orchestrate":    orchestrateBlockSchema(),
			"store":          storeBlockSchema(),
			"locals":         localsBlockSchema(),
		},
	}
}
