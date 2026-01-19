// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

var expectedMergedSchema_v013 = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"provider": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"grafana"}]}`: {
					Detail:   "grafana/grafana",
					HoverURL: "https://registry.terraform.io/providers/grafana/grafana/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/grafana/grafana/latest/docs",
						Tooltip: "grafana/grafana Documentation",
					},
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"auth": {
							Description: lang.MarkupContent{
								Value: "Credentials for accessing the Grafana API.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired:  true,
							IsSensitive: true,
							Constraint:  schema.AnyExpression{OfType: cty.String},
						},
						"url": {
							Description: lang.MarkupContent{
								Value: "URL of the root of the target Grafana server.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
				},
				`{"labels":[{"index":0,"value":"null"}]}`: {
					Detail:   "hashicorp/null",
					HoverURL: "https://registry.terraform.io/providers/hashicorp/null/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/hashicorp/null/latest/docs",
						Tooltip: "hashicorp/null Documentation",
					},
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
				},
				`{"labels":[{"index":0,"value":"rand"}]}`: {
					Detail:   "hashicorp/random",
					HoverURL: "https://registry.terraform.io/providers/hashicorp/random/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/hashicorp/random/latest/docs",
						Tooltip: "hashicorp/random Documentation",
					},
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
				},
			},
		},
		"resource": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "type",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				},
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Constraint: schema.AnyExpression{OfType: cty.Number},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"grafana_alert_notification"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"frequency": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"is_default": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"send_reminder": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"settings": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional:  true,
							IsSensitive: true,
						},
						"type": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"uid": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_dashboard"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"config_json": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"folder": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"slug": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_data_source"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{
						"json_data": {
							Labels: []*schema.LabelSchema{},
							Type:   schema.BlockTypeList,
							Body: &schema.BodySchema{
								Blocks: map[string]*schema.BlockSchema{},
								Attributes: map[string]*schema.AttributeSchema{
									"assume_role_arn": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"auth_type": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"conn_max_lifetime": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"custom_metrics_namespaces": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"default_region": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"encrypt": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"es_version": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"graphite_version": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"http_method": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"interval": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"log_level_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"log_message_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"max_idle_conns": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"max_open_conns": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"postgres_version": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"query_timeout": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"ssl_mode": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"time_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"time_interval": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"timescaledb": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_auth": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_auth_with_ca_cert": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_skip_verify": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tsdb_resolution": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"tsdb_version": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
								},
							},
						},
						"secure_json_data": {
							Labels: []*schema.LabelSchema{},
							Type:   schema.BlockTypeList,
							Body: &schema.BodySchema{
								Blocks: map[string]*schema.BlockSchema{},
								Attributes: map[string]*schema.AttributeSchema{
									"access_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"basic_auth_password": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"password": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"private_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"secret_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_ca_cert": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_client_cert": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_client_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
								},
							},
						},
					},
					Attributes: map[string]*schema.AttributeSchema{
						"access_mode": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"basic_auth_enabled": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"basic_auth_password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsOptional:  true,
							IsSensitive: true,
						},
						"basic_auth_username": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"database_name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"is_default": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsOptional:  true,
							IsSensitive: true,
						},
						"type": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"url": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"username": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_folder"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"title": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"uid": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_organization"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"admin_user": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"admins": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"create_users": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"editors": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"org_id": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsComputed: true,
						},
						"viewers": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_team"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"members": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"team_id": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_user"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"login": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsRequired:  true,
							IsSensitive: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_alert_notification"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"frequency": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"is_default": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"send_reminder": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"settings": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional:  true,
							IsSensitive: true,
						},
						"type": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"uid": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_dashboard"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"config_json": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"folder": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"slug": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_data_source"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{
						"json_data": {
							Labels: []*schema.LabelSchema{},
							Type:   schema.BlockTypeList,
							Body: &schema.BodySchema{
								Blocks: map[string]*schema.BlockSchema{},
								Attributes: map[string]*schema.AttributeSchema{
									"assume_role_arn": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"auth_type": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"conn_max_lifetime": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"custom_metrics_namespaces": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"default_region": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"encrypt": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"es_version": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"graphite_version": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"http_method": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"interval": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"log_level_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"log_message_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"max_idle_conns": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"max_open_conns": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"postgres_version": {
										Constraint: schema.AnyExpression{OfType: cty.Number},
										IsOptional: true,
									},
									"query_timeout": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"ssl_mode": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"time_field": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"time_interval": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"timescaledb": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_auth": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_auth_with_ca_cert": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tls_skip_verify": {
										Constraint: schema.AnyExpression{OfType: cty.Bool},
										IsOptional: true,
									},
									"tsdb_resolution": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
									"tsdb_version": {
										Constraint: schema.AnyExpression{OfType: cty.String},
										IsOptional: true,
									},
								},
							},
						},
						"secure_json_data": {
							Labels: []*schema.LabelSchema{},
							Type:   schema.BlockTypeList,
							Body: &schema.BodySchema{
								Blocks: map[string]*schema.BlockSchema{},
								Attributes: map[string]*schema.AttributeSchema{
									"access_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"basic_auth_password": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"password": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"private_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"secret_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_ca_cert": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_client_cert": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
									"tls_client_key": {
										Constraint:  schema.AnyExpression{OfType: cty.String},
										IsOptional:  true,
										IsSensitive: true,
									},
								},
							},
						},
					},
					Attributes: map[string]*schema.AttributeSchema{
						"access_mode": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"basic_auth_enabled": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"basic_auth_password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsOptional:  true,
							IsSensitive: true,
						},
						"basic_auth_username": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"database_name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"is_default": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsOptional:  true,
							IsSensitive: true,
						},
						"type": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"url": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"username": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_folder"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"title": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"uid": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_organization"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"admin_user": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"admins": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"create_users": {
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"editors": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"org_id": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsComputed: true,
						},
						"viewers": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_team"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"members": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"team_id": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsComputed: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_user"}],"attrs":[{"name":"provider","expr":{"addr":"grafana"}}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsRequired: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"login": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"name": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"password": {
							Constraint:  schema.AnyExpression{OfType: cty.String},
							IsRequired:  true,
							IsSensitive: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"triggers": {
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.MarkdownKind,
							},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}],"attrs":[{"name":"provider","expr":{"addr":"null"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}],"attrs":[{"name":"provider","expr":{"addr":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_id"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"b64_std": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64 without additional transformations.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"b64_url": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64, using the URL-friendly character set: case-sensitive letters, digits and the characters `_` and `-`.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"byte_length": {
							Description: lang.MarkupContent{
								Value: "The number of random bytes to produce. The minimum value is 1, which produces eight bits of randomness.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"dec": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in non-padded decimal digits.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"hex": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in padded hexadecimal digits. This result will always be twice as long as the requested byte length.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string to prefix the output value with. This string is supplied as-is, meaning it is not guaranteed to be URL-safe or base64 encoded.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "\nThe resource `random_id` generates random numbers that are intended to be\nused as unique identifiers for other resources.\n\nThis resource *does* use a cryptographic random number generator in order\nto minimize the chance of collisions, making the results of this resource\nwhen a 16-byte identifier is requested of equivalent uniqueness to a\ntype-4 UUID.\n\nThis resource can be used in conjunction with resources that have\nthe `create_before_destroy` lifecycle flag set to avoid conflicts with\nunique names during the brief period where both the old and new resources\nexist concurrently.\n",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_integer"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"max": {
							Description: lang.MarkupContent{
								Value: "The maximum inclusive value of the range.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"min": {
							Description: lang.MarkupContent{
								Value: "The minimum inclusive value of the range.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The random integer result.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "A custom seed to always produce the same value.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_integer` generates random values from a given range, described by the `min` and `max` attributes of a given resource.\n\nThis resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set, to avoid conflicts with unique names during the brief period where both the old and new resources exist concurrently.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_password"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed:  true,
							IsSensitive: true,
							Constraint:  schema.AnyExpression{OfType: cty.String},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "Identical to [random_string](string.html) with the exception that the result is treated as sensitive and, thus, _not_ displayed in console output.\n\nThis resource *does* use a cryptographic random number generator.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_pet"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length (in words) of the pet name.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "A string to prefix the name with.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"separator": {
							Description: lang.MarkupContent{
								Value: "The character to separate words in the pet name.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_pet` generates random pet names that are intended to be used as unique identifiers for other resources.\n\nThis resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set, to avoid conflicts with unique names during the brief period where both the old and new resources exist concurrently.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_shuffle"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"input": {
							Description: lang.MarkupContent{
								Value: "The list of strings to shuffle.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "Random permutation of the list of strings given in `input`.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.List(cty.String), SkipLiteralComplexTypes: true},
								schema.List{
									Elem: schema.AnyExpression{OfType: cty.String},
								},
							},
						},
						"result_count": {
							Description: lang.MarkupContent{
								Value: "The number of results to return. Defaults to the number of items in the `input` list. If fewer items are requested, some elements will be excluded from the result. If more items are requested, items will be repeated in the result but not more frequently than the number of items in the input list.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string with which to seed the random number generator, in order to produce less-volatile permutations of the list.\n\n**Important:** Even with an identical seed, it is not guaranteed that the same permutation will be produced across different versions of Terraform. This argument causes the result to be *less volatile*, but not fixed for all time.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_shuffle` generates a random permutation of a list of strings given as an argument.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_string"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							Constraint: schema.AnyExpression{OfType: cty.Number},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.Bool},
							IsOptional: true,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_string` generates a random permutation of alphanumeric characters and optionally special characters.\n\nThis resource *does* use a cryptographic random number generator.\n\nHistorically this resource's intended usage has been ambiguous as the original example used it in a password. For backwards compatibility it will continue to exist. For unique ids please use [random_id](id.html), for sensitive random values please use [random_password](password.html).",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_uuid"}],"attrs":[{"name":"provider","expr":{"addr":"rand"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated uuid presented in string format.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_uuid` generates random uuid string that is intended to be used as unique identifiers for other resources.\n\nThis resource uses [hashicorp/go-uuid](https://github.com/hashicorp/go-uuid) to generate a UUID-formatted string for use with services needed a unique string identifier.",
						Kind:  lang.MarkdownKind,
					},
				},
			},
		},
		"data": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "type",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				},
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {
						Constraint: schema.AnyExpression{OfType: cty.Number},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"null_data_source"}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_data_source` data source implements the standard data source lifecycle but does not interact with any external APIs.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_data_source"}],"attrs":[{"name":"provider","expr":{"addr":"null"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsComputed: true,
							IsOptional: true,
						},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_data_source` data source implements the standard data source lifecycle but does not interact with any external APIs.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_data_source"}],"attrs":[{"name":"provider","expr":{"addr":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"id": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
							IsComputed: true,
						},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.MarkdownKind,
							},
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.OneOf{
								schema.AnyExpression{OfType: cty.Map(cty.String), SkipLiteralComplexTypes: true},
								schema.Map{
									Elem:                  schema.AnyExpression{OfType: cty.String},
									AllowInterpolatedKeys: true,
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.MarkdownKind,
							},
							IsComputed: true,
							Constraint: schema.AnyExpression{OfType: cty.String},
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_data_source` data source implements the standard data source lifecycle but does not interact with any external APIs.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}]}`: {
					Detail: "(builtin)",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"backend": {
							Constraint:             backends.BackendTypesAsOneOfConstraint(v0_13_0),
							IsDepKey:               true,
							IsRequired:             true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
						},
						"defaults": {
							Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							IsOptional: true,
						},
						"outputs": {
							Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							IsComputed: true,
						},
						"workspace": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Detail: "(builtin)",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"backend": {
							Constraint:             backends.BackendTypesAsOneOfConstraint(v0_13_0),
							IsRequired:             true,
							IsDepKey:               true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
						},
						"defaults": {
							Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							IsOptional: true,
						},
						"outputs": {
							Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							IsComputed: true,
						},
						"workspace": {
							Constraint: schema.AnyExpression{OfType: cty.String},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"artifactory"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["artifactory"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"artifactory"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["artifactory"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"atlas"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["atlas"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"atlas"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["atlas"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azure"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["azure"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azure"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["azure"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azurerm"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["azurerm"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azurerm"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["azurerm"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"consul"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["consul"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"consul"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["consul"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"cos"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["cos"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"cos"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["cos"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcd"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["etcd"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcd"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["etcd"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcdv3"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["etcdv3"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcdv3"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["etcdv3"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"gcs"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["gcs"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"gcs"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["gcs"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"http"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["http"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"http"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["http"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"kubernetes"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["kubernetes"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"kubernetes"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["kubernetes"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"local"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["local"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"local"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["local"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"manta"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["manta"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"manta"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["manta"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"oss"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["oss"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"oss"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["oss"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"pg"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["pg"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"pg"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["pg"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"remote"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["remote"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"remote"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["remote"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"s3"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["s3"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"s3"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["s3"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"swift"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["swift"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"swift"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{
						"backend":   {IsRequired: true, Constraint: backends.BackendTypesAsOneOfConstraint(v0_13_0), SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent}},
						"config":    {IsOptional: true, Constraint: backends.ConfigsAsObjectConstraint(v0_13_0)["swift"]},
						"defaults":  {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"outputs":   {IsComputed: true, Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType}},
						"workspace": {IsOptional: true, Constraint: schema.AnyExpression{OfType: cty.String}},
					},
					Blocks: map[string]*schema.BlockSchema{},
					Detail: "(builtin)",
				},
			},
		},
		"module": {
			Labels: []*schema.LabelSchema{
				{
					Name:                   "name",
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				},
			},
			SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"source": {
						Constraint:             schema.LiteralType{Type: cty.String},
						IsRequired:             true,
						IsDepKey:               true,
						SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
					},
					"version": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
	},
}
