// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	v0_12_2  = version.Must(version.NewVersion("0.12.2"))
	v0_12_4  = version.Must(version.NewVersion("0.12.4"))
	v0_12_7  = version.Must(version.NewVersion("0.12.7"))
	v0_12_8  = version.Must(version.NewVersion("0.12.8"))
	v0_12_10 = version.Must(version.NewVersion("0.12.10"))
	v0_12_17 = version.Must(version.NewVersion("0.12.17"))
	v0_12_20 = version.Must(version.NewVersion("0.12.20"))
	v0_12_21 = version.Must(version.NewVersion("0.12.21"))
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
	f := BaseFunctions()

	if v.GreaterThanOrEqual(v0_12_2) {
		f["range"] = schema.FunctionSignature{
			VarParam: &function.Parameter{
				Name: "params",
				Type: cty.Number,
			},
			ReturnType:  cty.List(cty.Number),
			Description: "`range` generates a list of numbers using a start value, a limit value, and a step value.",
		}
		f["uuidv5"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "namespace",
					Type: cty.String,
				},
				{
					Name: "name",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`uuidv5` generates a _name-based_ UUID, as described in [RFC 4122 section 4.3](https://tools.ietf.org/html/rfc4122#section-4.3), also known as a \"version 5\" UUID.",
		}
		f["yamldecode"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "src",
					Type: cty.String,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`yamldecode` parses a string as a subset of YAML, and produces a representation of its value.",
		}
		f["yamlencode"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "value",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.String,
			Description: "`yamlencode` encodes a given value to a string using [YAML 1.2](https://yaml.org/spec/1.2/spec.html) block syntax.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_4) {
		f["abspath"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`abspath` takes a string containing a filesystem path and converts it to an absolute path. That is, if the path is not absolute, it will be joined with the current working directory.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_7) {
		f["regex"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "pattern",
					Type: cty.String,
				},
				{
					Name: "string",
					Type: cty.String,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`regex` applies a [regular expression](https://en.wikipedia.org/wiki/Regular_expression) to a string and returns the matching substrings.",
		}
		f["regexall"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "pattern",
					Type: cty.String,
				},
				{
					Name: "string",
					Type: cty.String,
				},
			},
			ReturnType:  cty.List(cty.DynamicPseudoType),
			Description: "`regexall` applies a [regular expression](https://en.wikipedia.org/wiki/Regular_expression) to a string and returns a list of all matches.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_8) {
		f["fileset"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
				{
					Name: "pattern",
					Type: cty.String,
				},
			},
			ReturnType:  cty.Set(cty.String),
			Description: "`fileset` enumerates a set of regular file names given a path and pattern. The path is automatically removed from the resulting set of file names and any result still containing path separators always returns forward slash (`/`) as the path separator for cross-system compatibility.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_10) {
		f["parseint"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "number",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "base",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`parseint` parses the given string as a representation of an integer in the specified base and returns the resulting number. The base must be between 2 and 62 inclusive.",
		}
		f["cidrsubnets"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name:        "prefix",
					Description: "`prefix` must be given in CIDR notation, as defined in [RFC 4632 section 3.1](https://tools.ietf.org/html/rfc4632#section-3.1).",
					Type:        cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name: "newbits",
				Type: cty.Number,
			},
			ReturnType:  cty.List(cty.String),
			Description: "`cidrsubnets` calculates a sequence of consecutive IP address ranges within a particular CIDR prefix.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_17) {
		f["trim"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
				{
					Name:        "cutset",
					Description: "A string containing all of the characters to trim. Each character is taken separately, so the order of characters is insignificant.",
					Type:        cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`trim` removes the specified set of characters from the start and end of the given string.",
		}
		f["trimprefix"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
				{
					Name: "prefix",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`trimprefix` removes the specified prefix from the start of the given string. If the string does not start with the prefix, the string is returned unchanged.",
		}
		f["trimspace"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`trimspace` removes any space characters from the start and end of the given string.",
		}
		f["trimsuffix"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
				{
					Name: "suffix",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`trimsuffix` removes the specified suffix from the end of the given string.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_20) {
		f["can"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "expression",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Bool,
			Description: "`can` evaluates the given expression and returns a boolean value indicating whether the expression produced a result without any errors.",
		}
		f["try"] = schema.FunctionSignature{
			VarParam: &function.Parameter{
				Name:        "expressions",
				Description: "",
				Type:        cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`try` evaluates all of its argument expressions in turn and returns the result of the first one that does not produce any errors.",
		}
	}
	if v.GreaterThanOrEqual(v0_12_21) {
		f["setsubtract"] = schema.FunctionSignature{
			Params: []function.Parameter{
				{
					Name: "a",
					Type: cty.Set(cty.DynamicPseudoType),
				},
				{
					Name: "b",
					Type: cty.Set(cty.DynamicPseudoType),
				},
			},
			ReturnType:  cty.Set(cty.DynamicPseudoType),
			Description: "The `setsubtract` function returns a new set containing the elements from the first set that are not present in the second set. In other words, it computes the [relative complement](https://en.wikipedia.org/wiki/Complement_\\(set_theory\\)#Relative_complement) of the second set.",
		}
	}

	return f
}
