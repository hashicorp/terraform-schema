// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func resourceLifecycleActionTriggerBlock() *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Description: lang.Markdown("A block that defines an action or actions to run depending on the condition (optional) and event (required)."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"events": {
					IsRequired:  true,
					Description: lang.Markdown("A list of events which will trigger the action in this"),
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Keyword{
								Keyword:     "before_create",
								Description: lang.Markdown("Run the trigger before the resource is created"),
							},
							schema.Keyword{
								Keyword:     "before_update",
								Description: lang.Markdown("Run the trigger before the resource is updated"),
							},
							schema.Keyword{
								Keyword:     "after_create",
								Description: lang.Markdown("Run the trigger after the resource is created"),
							},
							schema.Keyword{
								Keyword:     "after_update",
								Description: lang.Markdown("Run the trigger after the resource is updated"),
							},
						},
					},
				},
				"condition": {
					IsOptional:  true,
					Description: lang.Markdown("Condition, a boolean expression which can be used to further control whether or not Terraform invokes the actions defined by the given action_trigger; the actions will be invoked if this condition evaluates to true."),
					Constraint:  schema.AnyExpression{OfType: cty.Bool},
				},
				"actions": {
					IsRequired:  true,
					Description: lang.Markdown("An ordered list of actions which should run during the given event (if the optional `condition` evaluates to true)."),
					Constraint: schema.List{
						Elem: schema.Reference{OfScopeId: refscope.ActionScope},
					},
				},
			},
		},
	}

	return bs
}
