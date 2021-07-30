package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/zclconf/go-cty/cty"
)

var expectedMergedSchema_v012 = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"provider": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {
						Expr: schema.ExprConstraints{
							schema.LiteralTypeExpr{Type: cty.String},
						},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"null"}]}`: {
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
					Detail:     "hashicorp/null",
					HoverURL:   "https://registry.terraform.io/providers/hashicorp/null/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/hashicorp/null/latest/docs",
						Tooltip: "hashicorp/null Documentation",
					},
				},
				`{"labels":[{"index":0,"value":"random"}]}`: {
					Blocks:     map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{},
					Detail:     "hashicorp/random",
					HoverURL:   "https://registry.terraform.io/providers/hashicorp/random/latest/docs",
					DocsLink: &schema.DocsLink{
						URL:     "https://registry.terraform.io/providers/hashicorp/random/latest/docs",
						Tooltip: "hashicorp/random Documentation",
					},
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
					"count": {
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.Number},
							schema.LiteralTypeExpr{Type: cty.Number},
						},
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{
				`{"labels":[{"index":0,"value":"null_resource"}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
							IsComputed: true,
						},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}],"attrs":[{"name":"provider","expr":{"addr":"null"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"null_resource"}],"attrs":[{"name":"provider","expr":{"addr":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"triggers": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_id"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"b64_std": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64 without additional transformations.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"b64_url": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64, using the URL-friendly character set: case-sensitive letters, digits and the characters `_` and `-`.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"byte_length": {
							Description: lang.MarkupContent{
								Value: "The number of random bytes to produce. The minimum value is 1, which produces eight bits of randomness.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"dec": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in non-padded decimal digits.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"hex": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in padded hexadecimal digits. This result will always be twice as long as the requested byte length.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string to prefix the output value with. This string is supplied as-is, meaning it is not guaranteed to be URL-safe or base64 encoded.",
								Kind:  lang.PlainTextKind,
							},
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_integer"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"max": {
							Description: lang.MarkupContent{
								Value: "The maximum inclusive value of the range.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"min": {
							Description: lang.MarkupContent{
								Value: "The minimum inclusive value of the range.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The random integer result.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "A custom seed to always produce the same value.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_password"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed:  true,
							IsSensitive: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_pet"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length (in words) of the pet name.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "A string to prefix the name with.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"separator": {
							Description: lang.MarkupContent{
								Value: "The character to separate words in the pet name.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_shuffle"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"input": {
							Description: lang.MarkupContent{
								Value: "The list of strings to shuffle.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.List(cty.String)},
								schema.LiteralTypeExpr{Type: cty.List(cty.String)},
								schema.ListExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "Random permutation of the list of strings given in `input`.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.List(cty.String)},
								schema.LiteralTypeExpr{Type: cty.List(cty.String)},
								schema.ListExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"result_count": {
							Description: lang.MarkupContent{
								Value: "The number of results to return. Defaults to the number of items in the `input` list. If fewer items are requested, some elements will be excluded from the result. If more items are requested, items will be repeated in the result but not more frequently than the number of items in the input list.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string with which to seed the random number generator, in order to produce less-volatile permutations of the list.\n\n**Important:** Even with an identical seed, it is not guaranteed that the same permutation will be produced across different versions of Terraform. This argument causes the result to be *less volatile*, but not fixed for all time.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_string"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_uuid"}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated uuid presented in string format.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_id"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"b64_std": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64 without additional transformations.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"b64_url": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in base64, using the URL-friendly character set: case-sensitive letters, digits and the characters `_` and `-`.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"byte_length": {
							Description: lang.MarkupContent{
								Value: "The number of random bytes to produce. The minimum value is 1, which produces eight bits of randomness.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"dec": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in non-padded decimal digits.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"hex": {
							Description: lang.MarkupContent{
								Value: "The generated id presented in padded hexadecimal digits. This result will always be twice as long as the requested byte length.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string to prefix the output value with. This string is supplied as-is, meaning it is not guaranteed to be URL-safe or base64 encoded.",
								Kind:  lang.PlainTextKind,
							},
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_integer"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"max": {
							Description: lang.MarkupContent{
								Value: "The maximum inclusive value of the range.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"min": {
							Description: lang.MarkupContent{
								Value: "The minimum inclusive value of the range.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The random integer result.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "A custom seed to always produce the same value.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_password"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed:  true,
							IsSensitive: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_pet"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length (in words) of the pet name.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"prefix": {
							Description: lang.MarkupContent{
								Value: "A string to prefix the name with.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"separator": {
							Description: lang.MarkupContent{
								Value: "The character to separate words in the pet name.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_shuffle"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"input": {
							Description: lang.MarkupContent{
								Value: "The list of strings to shuffle.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.List(cty.String)},
								schema.LiteralTypeExpr{Type: cty.List(cty.String)},
								schema.ListExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "Random permutation of the list of strings given in `input`.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.List(cty.String)},
								schema.LiteralTypeExpr{Type: cty.List(cty.String)},
								schema.ListExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"result_count": {
							Description: lang.MarkupContent{
								Value: "The number of results to return. Defaults to the number of items in the `input` list. If fewer items are requested, some elements will be excluded from the result. If more items are requested, items will be repeated in the result but not more frequently than the number of items in the input list.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"seed": {
							Description: lang.MarkupContent{
								Value: "Arbitrary string with which to seed the random number generator, in order to produce less-volatile permutations of the list.\n\n**Important:** Even with an identical seed, it is not guaranteed that the same permutation will be produced across different versions of Terraform. This argument causes the result to be *less volatile*, but not fixed for all time.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_string"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"length": {
							Description: lang.MarkupContent{
								Value: "The length of the string desired.",
								Kind:  lang.PlainTextKind,
							},
							IsRequired: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
						},
						"lower": {
							Description: lang.MarkupContent{
								Value: "Include lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"min_lower": {
							Description: lang.MarkupContent{
								Value: "Minimum number of lowercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_numeric": {
							Description: lang.MarkupContent{
								Value: "Minimum number of numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_special": {
							Description: lang.MarkupContent{
								Value: "Minimum number of special characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"min_upper": {
							Description: lang.MarkupContent{
								Value: "Minimum number of uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
						"number": {
							Description: lang.MarkupContent{
								Value: "Include numeric characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"override_special": {
							Description: lang.MarkupContent{
								Value: "Supply your own list of special characters to use for string generation.  This overrides the default character list in the special argument.  The `special` argument must still be set to true for any overwritten characters to be used in generation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated random string.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
						"special": {
							Description: lang.MarkupContent{
								Value: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
						"upper": {
							Description: lang.MarkupContent{
								Value: "Include uppercase alphabet characters in the result.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Bool},
								schema.LiteralTypeExpr{Type: cty.Bool},
							},
							IsOptional: true,
						},
					},
				},
				`{"labels":[{"index":0,"value":"random_uuid"}],"attrs":[{"name":"provider","expr":{"addr":"random"}}]}`: {
					Detail: "hashicorp/random",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"keepers": {
							Description: lang.MarkupContent{
								Value: "Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"result": {
							Description: lang.MarkupContent{
								Value: "The generated uuid presented in string format.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
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
					"count": {
						Expr: schema.ExprConstraints{
							schema.TraversalExpr{OfType: cty.Number},
							schema.LiteralTypeExpr{Type: cty.Number},
						},
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
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsComputed: true,
							IsOptional: true,
						},
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"null_data_source"}],"attrs":[{"name":"provider","expr":{"addr":"null"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
							IsComputed: true,
						},
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"null_data_source"}],"attrs":[{"name":"provider","expr":{"addr":"null.foobar"}}]}`: {
					Detail: "hashicorp/null",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"has_computed_default": {
							Description: lang.MarkupContent{
								Value: "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true,
							IsComputed: true,
						},
						"id": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
							IsOptional: true, IsComputed: true},
						"inputs": {
							Description: lang.MarkupContent{
								Value: "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
								Kind:  lang.PlainTextKind,
							},
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
							IsOptional: true,
						},
						"outputs": {
							Description: lang.MarkupContent{
								Value: "After the data source is \"read\", a copy of the `inputs` map.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Map(cty.String)},
								schema.LiteralTypeExpr{Type: cty.Map(cty.String)},
								schema.MapExpr{
									Elem: schema.ExprConstraints{
										schema.TraversalExpr{OfType: cty.String},
										schema.LiteralTypeExpr{Type: cty.String},
									},
								},
							},
						},
						"random": {
							Description: lang.MarkupContent{
								Value: "A random value. This is primarily for testing and has little practical use; prefer the [random provider](https://www.terraform.io/docs/providers/random/) for more practical random number use-cases.",
								Kind:  lang.PlainTextKind,
							},
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}]}`: {
					Detail: "(builtin)",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"backend": {
							IsRequired: true,
							IsDepKey:   true,
							Expr:       backends.BackendTypesAsExprConstraints(v0_12_0),
						},
						"defaults": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.DynamicPseudoType},
								schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
							},
						},
						"outputs": {
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.DynamicPseudoType},
								schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
							},
						},
						"workspace": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Detail: "(builtin)",
					Blocks: map[string]*schema.BlockSchema{},
					Attributes: map[string]*schema.AttributeSchema{
						"backend": {
							IsRequired: true,
							IsDepKey:   true,
							Expr:       backends.BackendTypesAsExprConstraints(v0_12_0),
						},
						"defaults": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.DynamicPseudoType},
								schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
							},
						},
						"outputs": {
							IsComputed: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.DynamicPseudoType},
								schema.LiteralTypeExpr{Type: cty.DynamicPseudoType},
							},
						},
						"workspace": {
							IsOptional: true,
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.String},
								schema.LiteralTypeExpr{Type: cty.String},
							},
						},
					},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"artifactory"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["artifactory"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"artifactory"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["artifactory"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"atlas"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["atlas"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"atlas"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["atlas"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azure"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["azure"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azure"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["azure"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azurerm"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["azurerm"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"azurerm"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["azurerm"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"consul"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["consul"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"consul"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["consul"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcd"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["etcd"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcd"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["etcd"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcdv3"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["etcdv3"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"etcdv3"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["etcdv3"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"gcs"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["gcs"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"gcs"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["gcs"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"http"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["http"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"http"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["http"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"local"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["local"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"local"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["local"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"manta"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["manta"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"manta"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["manta"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"pg"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["pg"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"pg"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["pg"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"remote"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["remote"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"remote"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["remote"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"s3"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["s3"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"s3"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["s3"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"swift"}},{"name":"provider","expr":{"addr":"terraform"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["swift"]}},
				},
				`{"labels":[{"index":0,"value":"terraform_remote_state"}],"attrs":[{"name":"backend","expr":{"static":"swift"}}]}`: {
					Attributes: map[string]*schema.AttributeSchema{"config": {IsOptional: true, Expr: backends.ConfigsAsExprConstraints(v0_12_0)["swift"]}},
				},
			},
		},
		"module": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"source": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsRequired: true,
						IsDepKey:   true,
					},
					"version": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsOptional: true,
					},
				},
			},
			DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
		},
	},
}
