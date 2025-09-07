// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func resourceLifecycleActionTriggerBlockSchema() *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Description: lang.Markdown("a block that defines an action or actions to run depending on the condition (optional) and event (required)."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"events": {
					IsRequired:  true,
					Description: lang.Markdown("(Required, list): A list of events which will trigger the action in this"),
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "before_create",
							Description: lang.Markdown("Run the trigger before the resource is created"),
						},
						schema.Keyword{
							Keyword:     "before_update",
							Description: lang.Markdown("Run the trigger before the resource is updated"),
						},
						schema.Keyword{
							Keyword:     "before_destroy",
							Description: lang.Markdown("Run the trigger before the resource is destroyed"),
						},
						schema.Keyword{
							Keyword:     "after_create",
							Description: lang.Markdown("Run the trigger after the resource is created"),
						},
						schema.Keyword{
							Keyword:     "after_update",
							Description: lang.Markdown("Run the trigger after the resource is updated"),
						},
						schema.Keyword{
							Keyword:     "after_destroy",
							Description: lang.Markdown("Run the trigger after the resource is destroyed"),
						},
						schema.Keyword{
							Keyword:     "on_create",
							Description: lang.Markdown("Run the trigger when the resource is created"),
						},
						schema.Keyword{
							Keyword:     "on_update",
							Description: lang.Markdown("Run the trigger when the resource is updated"),
						},
						schema.Keyword{
							Keyword:     "on_destroy",
							Description: lang.Markdown("Run the trigger when the resource is destroyed"),
						},
					},
				},
				"condition": {
					IsOptional:  true,
					Description: lang.Markdown("(optional): an optional boolean expression which can be used to further control whether or not Terraform invokes the actions defined by the given action_trigger; the actions will be invoked if this condition evaluates to true."),
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
