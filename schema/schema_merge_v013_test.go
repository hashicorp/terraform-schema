package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var expectedMergedSchema_v013 = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"provider": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {ValueType: cty.String},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"grafana"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"auth": {
							Description: lang.MarkupContent{
								Value: "Credentials for accessing the Grafana API.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							ValueType:  cty.String,
						},
						"url": {
							Description: lang.MarkupContent{
								Value: "URL of the root of the target Grafana server.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							ValueType:  cty.String,
						},
					},
				},
				`{"labels":[{"index":0,"value":"null"}]}`: {
					Detail:     "hashicorp/null",
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
				},
				`{"labels":[{"index":0,"value":"random"}]}`: {
					Detail:     "hashicorp/random",
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
				},
				`{"labels":[{"index":0,"value":"terraform"}]}`: {
					Detail:     "(builtin)",
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
				},
			},
		},
		"resource": {
			Labels: []*schema.LabelSchema{
				{Name: "type"},
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {ValueType: cty.Number},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"grafana_alert_notification"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"frequency":     {ValueType: cty.String},
						"id":            {ValueType: cty.String},
						"is_default":    {ValueType: cty.Bool},
						"name":          {IsRequired: true, ValueType: cty.String},
						"send_reminder": {ValueType: cty.Bool},
						"settings": {
							ValueType: cty.Map(cty.String),
						},
						"type": {IsRequired: true, ValueType: cty.String},
						"uid":  {ValueType: cty.String},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_dashboard"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"config_json": {IsRequired: true, ValueType: cty.String},
						"folder":      {ValueType: cty.Number},
						"id":          {ValueType: cty.String},
						"slug":        {IsReadOnly: true, ValueType: cty.String},
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
									"assume_role_arn":           {ValueType: cty.String},
									"auth_type":                 {ValueType: cty.String},
									"conn_max_lifetime":         {ValueType: cty.Number},
									"custom_metrics_namespaces": {ValueType: cty.String},
									"default_region":            {ValueType: cty.String},
									"encrypt":                   {ValueType: cty.String},
									"es_version":                {ValueType: cty.Number},
									"graphite_version":          {ValueType: cty.String},
									"http_method":               {ValueType: cty.String},
									"interval":                  {ValueType: cty.String},
									"log_level_field":           {ValueType: cty.String},
									"log_message_field":         {ValueType: cty.String},
									"max_idle_conns":            {ValueType: cty.Number},
									"max_open_conns":            {ValueType: cty.Number},
									"postgres_version":          {ValueType: cty.Number},
									"query_timeout":             {ValueType: cty.String},
									"ssl_mode":                  {ValueType: cty.String},
									"time_field":                {ValueType: cty.String},
									"time_interval":             {ValueType: cty.String},
									"timescaledb":               {ValueType: cty.Bool},
									"tls_auth":                  {ValueType: cty.Bool},
									"tls_auth_with_ca_cert":     {ValueType: cty.Bool},
									"tls_skip_verify":           {ValueType: cty.Bool},
									"tsdb_resolution":           {ValueType: cty.String},
									"tsdb_version":              {ValueType: cty.String},
								},
							},
						},
						"secure_json_data": {
							Labels: []*schema.LabelSchema{},
							Type:   schema.BlockTypeList,
							Body: &schema.BodySchema{
								Blocks: map[string]*schema.BlockSchema{},
								Attributes: map[string]*schema.AttributeSchema{
									"access_key":          {ValueType: cty.String},
									"basic_auth_password": {ValueType: cty.String},
									"password":            {ValueType: cty.String},
									"private_key":         {ValueType: cty.String},
									"secret_key":          {ValueType: cty.String},
									"tls_ca_cert":         {ValueType: cty.String},
									"tls_client_cert":     {ValueType: cty.String},
									"tls_client_key":      {ValueType: cty.String},
								},
							},
						},
					},
					Attributes: map[string]*schema.AttributeSchema{
						"access_mode":         {ValueType: cty.String},
						"basic_auth_enabled":  {ValueType: cty.Bool},
						"basic_auth_password": {ValueType: cty.String},
						"basic_auth_username": {ValueType: cty.String},
						"database_name":       {ValueType: cty.String},
						"id":                  {ValueType: cty.String},
						"is_default":          {ValueType: cty.Bool},
						"name":                {IsRequired: true, ValueType: cty.String},
						"password":            {ValueType: cty.String},
						"type":                {IsRequired: true, ValueType: cty.String},
						"url":                 {ValueType: cty.String},
						"username":            {ValueType: cty.String},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_folder"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id":    {ValueType: cty.String},
						"title": {IsRequired: true, ValueType: cty.String},
						"uid":   {IsReadOnly: true, ValueType: cty.String},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_organization"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"admin_user": {ValueType: cty.String},
						"admins": {
							ValueType: cty.List(cty.String),
						},
						"create_users": {ValueType: cty.Bool},
						"editors": {
							ValueType: cty.List(cty.String),
						},
						"id":     {ValueType: cty.String},
						"name":   {IsRequired: true, ValueType: cty.String},
						"org_id": {IsReadOnly: true, ValueType: cty.Number},
						"viewers": {
							ValueType: cty.List(cty.String),
						},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_team"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email": {ValueType: cty.String},
						"id":    {ValueType: cty.String},
						"members": {
							ValueType: cty.List(cty.String),
						},
						"name":    {IsRequired: true, ValueType: cty.String},
						"team_id": {IsReadOnly: true, ValueType: cty.Number},
					},
				},
				`{"labels":[{"index":0,"value":"grafana_user"}]}`: {
					Detail: "grafana/grafana",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"email":    {IsRequired: true, ValueType: cty.String},
						"id":       {ValueType: cty.String},
						"login":    {ValueType: cty.String},
						"name":     {ValueType: cty.String},
						"password": {IsRequired: true, ValueType: cty.String},
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}],"attrs":[{"name":"provider","expr":{"ref":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_id"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"b64_std": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64 without additional transformations.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"b64_url": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64, using the URL-friendly character set: case-sensitive letters, digits and the characters `_` and `-`.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"byte_length": {
							Description: lang.MarkupContent{
								Value: "The number of random bytes to produce. The minimum value is 1, which produces eight bits of randomness.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.Number,
						},
						"dec": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in non-padded decimal digits.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"hex": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in padded hexadecimal digits. This result will always be twice as long as the requested byte length.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string to prefix the output value with. This string is supplied as-is, meaning it is not guaranteed to be URL-safe or base64 encoded.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
					},
					Description: lang.MarkupContent{
						Value: "\nThe resource `random_id` generates random numbers that are intended to be\nused as unique identifiers for other resources.\n\nThis resource *does* use a cryptographic random number generator in order\nto minimize the chance of collisions, making the results of this resource\nwhen a 16-byte identifier is requested of equivalent uniqueness to a\ntype-4 UUID.\n\nThis resource can be used in conjunction with resources that have\nthe `create_before_destroy` lifecycle flag set to avoid conflicts with\nunique names during the brief period where both the old and new resources\nexist concurrently.\n",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_integer"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"max": {
							Description: lang.MarkupContent{
								Value: "The maximum inclusive value of the range.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.Number,
						},
						"min": {
							Description: lang.MarkupContent{
								Value: "The minimum inclusive value of the range.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.Number,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The random integer result.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.Number,
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "A custom seed to always produce the same value.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_integer` generates random values from a given range, described by the `min` and `max` attributes of a given resource.\n\nThis resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set, to avoid conflicts with unique names during the brief period where both the old and new resources exist concurrently.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_password"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.Number,
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
					},
					Description: lang.MarkupContent{
						Value: "Identical to [random_string](string.html) with the exception that the result is treated as sensitive and, thus, _not_ displayed in console output.\n\nThis resource *does* use a cryptographic random number generator.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_pet"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length (in words) of the pet name.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "A string to prefix the name with.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
						"separator": {
							Description: lang.MarkupContent{
								Value: "The character to separate words in the pet name.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_pet` generates random pet names that are intended to be used as unique identifiers for other resources.\n\nThis resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set, to avoid conflicts with unique names during the brief period where both the old and new resources exist concurrently.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_shuffle"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"input": {
							Description: lang.MarkupContent{
								Value: "The list of strings to shuffle.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.List(cty.String),
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "Random permutation of the list of strings given in `input`.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.List(cty.String),
						},
						"result_count": {
							Description: lang.MarkupContent{
								Value: "The number of results to return. Defaults to the number of items in the `input` list. If fewer items are requested, some elements will be excluded from the result. If more items are requested, items will be repeated in the result but not more frequently than the number of items in the input list.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string with which to seed the random number generator, in order to produce less-volatile permutations of the list.\n\n**Important:** Even with an identical seed, it is not guaranteed that the same permutation will be produced across different versions of Terraform. This argument causes the result to be *less volatile*, but not fixed for all time.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_shuffle` generates a random permutation of a list of strings given as an argument.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_string"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.MarkdownKind,
							},
							IsRequired: true,
							ValueType:  cty.Number,
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Number,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Bool,
						},
					},
					Description: lang.MarkupContent{
						Value: "The resource `random_string` generates a random permutation of alphanumeric characters and optionally special characters.\n\nThis resource *does* use a cryptographic random number generator.\n\nHistorically this resource's intended usage has been ambiguous as the original example used it in a password. For backwards compatibility it will continue to exist. For unique ids please use [random_id](id.html), for sensitive random values please use [random_password](password.html).",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"random_uuid"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {ValueType: cty.String},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated uuid presented in string format.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
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
				{Name: "type"},
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"count": {ValueType: cty.Number},
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
							ValueType: cty.String,
						},
						"id": {ValueType: cty.String},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.Map(cty.String),
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
						},
					},
					Description: lang.MarkupContent{
						Value: "The `null_data_source` data source implements the standard data source lifecycle but does not interact with any external APIs.",
						Kind:  lang.MarkdownKind,
					},
				},
				`{"labels":[{"index":0,"value":"null_data_source"}],"attrs":[{"name":"provider","expr":{"ref":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.String,
						},
						"id": {ValueType: cty.String},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.MarkdownKind,
							},
							ValueType: cty.Map(cty.String),
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.Map(cty.String),
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.MarkdownKind,
							},
							IsReadOnly: true,
							ValueType:  cty.String,
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
						"backend":   {IsRequired: true, ValueType: cty.String},
						"config":    {ValueType: cty.DynamicPseudoType},
						"defaults":  {ValueType: cty.DynamicPseudoType},
						"outputs":   {IsReadOnly: true, ValueType: cty.DynamicPseudoType},
						"workspace": {ValueType: cty.String},
					},
				},
			},
		},
	},
}
