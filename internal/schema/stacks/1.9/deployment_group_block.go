// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func deploymentGroupBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("Defines a group that you can assign individual deployments to join. Deployment groups let you enforce orchestration rules on the deployments within the group"),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "deployment_group"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "deployment_group",
			ScopeId:      refscope.DeploymentGroupScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Deployment Group Name"),
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"auto_approve_checks": {
					Description: lang.Markdown("A list of references to `deployment_auto_approve` blocks. If all the checks in the referenced blocks pass, then plans for the deployments in this group automatically apply"),
					IsRequired:  true,
					Constraint: schema.List{
						Elem: schema.Reference{OfScopeId: refscope.DeploymentAutoApproveScope},
					},
				},
			},
		},
	}
}
