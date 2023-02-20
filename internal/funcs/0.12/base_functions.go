package funcs

import (
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func BaseFunctions() map[string]schema.FunctionSignature {
	return map[string]schema.FunctionSignature{
		"abs": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`abs` returns the absolute value of the given number. In other words, if the number is zero or positive then it is returned as-is, but if it is negative then it is multiplied by -1 to make it positive before returning it.",
		},
		"base64decode": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`base64decode` takes a string containing a Base64 character sequence and returns the original string.",
		},
		"base64encode": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`base64encode` applies Base64 encoding to a string.",
		},
		"base64gzip": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`base64gzip` compresses a string with gzip and then encodes the result in Base64 encoding.",
		},
		"base64sha256": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`base64sha256` computes the SHA256 hash of a given string and encodes it with Base64. This is not equivalent to `base64encode(sha256(\"test\"))` since `sha256()` returns hexadecimal representation.",
		},
		"base64sha512": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`base64sha512` computes the SHA512 hash of a given string and encodes it with Base64. This is not equivalent to `base64encode(sha512(\"test\"))` since `sha512()` returns hexadecimal representation.",
		},
		"basename": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`basename` takes a string containing a filesystem path and removes all except the last portion from it.",
		},
		"bcrypt": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name:        "cost",
				Description: "The `cost` argument is optional and will default to 10 if unspecified.",
				Type:        cty.Number,
			},
			ReturnType:  cty.String,
			Description: "`bcrypt` computes a hash of the given string using the Blowfish cipher, returning a string in [the _Modular Crypt Format_](https://passlib.readthedocs.io/en/stable/modular_crypt_format.html) usually expected in the shadow password file on many Unix systems.",
		},
		"ceil": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`ceil` returns the closest whole number that is greater than or equal to the given value, which may be a fraction.",
		},
		"chomp": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`chomp` removes newline characters at the end of a string.",
		},
		"chunklist": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.List(cty.DynamicPseudoType),
				},
				{
					Name:        "size",
					Description: "The maximum length of each chunk. All but the last element of the result is guaranteed to be of exactly this size.",
					Type:        cty.Number,
				},
			},
			ReturnType:  cty.List(cty.List(cty.DynamicPseudoType)),
			Description: "`chunklist` splits a single list into fixed-size chunks, returning a list of lists.",
		},
		"cidrhost": {
			Params: []function.Parameter{
				{
					Name:        "prefix",
					Description: "`prefix` must be given in CIDR notation, as defined in [RFC 4632 section 3.1](https://tools.ietf.org/html/rfc4632#section-3.1).",
					Type:        cty.String,
				},
				{
					Name:        "hostnum",
					Description: "`hostnum` is a whole number that can be represented as a binary integer with no more than the number of digits remaining in the address after the given prefix.",
					Type:        cty.Number,
				},
			},
			ReturnType:  cty.String,
			Description: "`cidrhost` calculates a full host IP address for a given host number within a given IP network address prefix.",
		},
		"cidrnetmask": {
			Params: []function.Parameter{
				{
					Name:        "prefix",
					Description: "`prefix` must be given in CIDR notation, as defined in [RFC 4632 section 3.1](https://tools.ietf.org/html/rfc4632#section-3.1).",
					Type:        cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`cidrnetmask` converts an IPv4 address prefix given in CIDR notation into a subnet mask address.",
		},
		"cidrsubnet": {
			Params: []function.Parameter{
				{
					Name:        "prefix",
					Description: "`prefix` must be given in CIDR notation, as defined in [RFC 4632 section 3.1](https://tools.ietf.org/html/rfc4632#section-3.1).",
					Type:        cty.String,
				},
				{
					Name:        "newbits",
					Description: "`newbits` is the number of additional bits with which to extend the prefix.",
					Type:        cty.Number,
				},
				{
					Name:        "netnum",
					Description: "`netnum` is a whole number that can be represented as a binary integer with no more than `newbits` binary digits, which will be used to populate the additional bits added to the prefix.",
					Type:        cty.Number,
				},
			},
			ReturnType:  cty.String,
			Description: "`cidrsubnet` calculates a subnet address within given IP network address prefix.",
		},
		"coalesce": {
			VarParam: &function.Parameter{
				Name: "vals",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`coalesce` takes any number of arguments and returns the first one that isn't null or an empty string.",
		},
		"coalescelist": {
			VarParam: &function.Parameter{
				Name:        "vals",
				Description: "List or tuple values to test in the given order.",
				Type:        cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`coalescelist` takes any number of list arguments and returns the first one that isn't empty.",
		},
		"compact": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.List(cty.String),
				},
			},
			ReturnType:  cty.List(cty.String),
			Description: "`compact` takes a list of strings and returns a new list with any empty string elements removed.",
		},
		"concat": {
			VarParam: &function.Parameter{
				Name: "seqs",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`concat` takes two or more lists and combines them into a single list.",
		},
		"contains": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "value",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`contains` determines whether a given list or set contains a given single value as one of its elements.",
		},
		"csvdecode": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`csvdecode` decodes a string containing CSV-formatted data and produces a list of maps representing that data.",
		},
		"dirname": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`dirname` takes a string containing a filesystem path and removes the last portion from it.",
		},
		"distinct": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.List(cty.DynamicPseudoType),
				},
			},
			ReturnType:  cty.List(cty.DynamicPseudoType),
			Description: "`distinct` takes a list and returns a new list with any duplicate elements removed.",
		},
		"element": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "index",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`element` retrieves a single element from a list.",
		},
		"file": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`file` reads the contents of a file at the given path and returns them as a string.",
		},
		"filebase64": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filebase64` reads the contents of a file at the given path and returns them as a base64-encoded string.",
		},
		"filebase64sha256": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filebase64sha256` is a variant of `base64sha256` that hashes the contents of a given file rather than a literal string.",
		},
		"filebase64sha512": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filebase64sha512` is a variant of `base64sha512` that hashes the contents of a given file rather than a literal string.",
		},
		"fileexists": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.Bool,
			Description: "`fileexists` determines whether a file exists at a given path.",
		},
		"filemd5": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filemd5` is a variant of `md5` that hashes the contents of a given file rather than a literal string.",
		},
		"filesha1": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filesha1` is a variant of `sha1` that hashes the contents of a given file rather than a literal string.",
		},
		"filesha256": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filesha256` is a variant of `sha256` that hashes the contents of a given file rather than a literal string.",
		},
		"filesha512": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`filesha512` is a variant of `sha512` that hashes the contents of a given file rather than a literal string.",
		},
		"flatten": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`flatten` takes a list and replaces any elements that are lists with a flattened sequence of the list contents.",
		},
		"floor": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`floor` returns the closest whole number that is less than or equal to the given value, which may be a fraction.",
		},
		"format": {
			Params: []function.Parameter{
				{
					Name: "format",
					Type: cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name: "args",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "The `format` function produces a string by formatting a number of other values according to a specification string. It is similar to the `printf` function in C, and other similar functions in other programming languages.",
		},
		"formatdate": {
			Params: []function.Parameter{
				{
					Name: "format",
					Type: cty.String,
				},
				{
					Name: "time",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`formatdate` converts a timestamp into a different time format.",
		},
		"formatlist": {
			Params: []function.Parameter{
				{
					Name: "format",
					Type: cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name: "args",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`formatlist` produces a list of strings by formatting a number of other values according to a specification string.",
		},
		"indent": {
			Params: []function.Parameter{
				{
					Name:        "spaces",
					Description: "Number of spaces to add after each newline character.",
					Type:        cty.Number,
				},
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`indent` adds a given number of spaces to the beginnings of all but the first line in a given multi-line string.",
		},
		"index": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "value",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`index` finds the element index for a given value in a list.",
		},
		"join": {
			Params: []function.Parameter{
				{
					Name:        "separator",
					Description: "Delimiter to insert between the given strings.",
					Type:        cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name:        "lists",
				Description: "One or more lists of strings to join.",
				Type:        cty.List(cty.String),
			},
			ReturnType:  cty.String,
			Description: "`join` produces a string by concatenating together all elements of a given list of strings with the given delimiter.",
		},
		"jsondecode": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`jsondecode` interprets a given string as JSON, returning a representation of the result of decoding that string.",
		},
		"jsonencode": {
			Params: []function.Parameter{
				{
					Name: "val",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.String,
			Description: "`jsonencode` encodes a given value to a string using JSON syntax.",
		},
		"keys": {
			Params: []function.Parameter{
				{
					Name:        "inputMap",
					Description: "The map to extract keys from. May instead be an object-typed value, in which case the result is a tuple of the object attributes.",
					Type:        cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`keys` takes a map and returns a list containing the keys from that map.",
		},
		"length": {
			Params: []function.Parameter{
				{
					Name: "value",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Number,
			Description: "`length` determines the length of a given list, map, or string.",
		},
		"list": {
			VarParam: &function.Parameter{
				Name: "vals",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.List(cty.DynamicPseudoType),
			Description: "`list` takes any number of list arguments and returns a list containing those values in the same order.",
		},
		"log": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
				{
					Name: "base",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`log` returns the logarithm of a given number in a given base.",
		},
		"lookup": {
			Params: []function.Parameter{
				{
					Name: "inputMap",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "key",
					Type: cty.String,
				},
			},
			VarParam: &function.Parameter{
				Name: "default",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`lookup` retrieves the value of a single element from a map, given its key. If the given key does not exist, the given default value is returned instead.",
		},
		"lower": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`lower` converts all cased letters in the given string to lowercase.",
		},
		"map": {
			VarParam: &function.Parameter{
				Name: "vals",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.Map(cty.DynamicPseudoType),
			Description: "`map` takes an even number of arguments and returns a map whose elements are constructed from consecutive pairs of arguments.",
		},
		"matchkeys": {
			Params: []function.Parameter{
				{
					Name: "values",
					Type: cty.List(cty.DynamicPseudoType),
				},
				{
					Name: "keys",
					Type: cty.List(cty.DynamicPseudoType),
				},
				{
					Name: "searchset",
					Type: cty.List(cty.DynamicPseudoType),
				},
			},
			ReturnType:  cty.List(cty.DynamicPseudoType),
			Description: "`matchkeys` constructs a new list by taking a subset of elements from one list whose indexes match the corresponding indexes of values in another list.",
		},
		"max": {
			VarParam: &function.Parameter{
				Name: "numbers",
				Type: cty.Number,
			},
			ReturnType:  cty.Number,
			Description: "`max` takes one or more numbers and returns the greatest number from the set.",
		},
		"md5": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`md5` computes the MD5 hash of a given string and encodes it with hexadecimal digits.",
		},
		"merge": {
			VarParam: &function.Parameter{
				Name: "maps",
				Type: cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`merge` takes an arbitrary number of maps or objects, and returns a single map or object that contains a merged set of elements from all arguments.",
		},
		"min": {
			VarParam: &function.Parameter{
				Name: "numbers",
				Type: cty.Number,
			},
			ReturnType:  cty.Number,
			Description: "`min` takes one or more numbers and returns the smallest number from the set.",
		},
		"pathexpand": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`pathexpand` takes a filesystem path that might begin with a `~` segment, and if so it replaces that segment with the current user's home directory path.",
		},
		"pow": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
				{
					Name: "power",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`pow` calculates an exponent, by raising its first argument to the power of the second argument.",
		},
		"replace": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
				{
					Name: "substr",
					Type: cty.String,
				},
				{
					Name: "replace",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`replace` searches a given string for another given substring, and replaces each occurrence with a given replacement string.",
		},
		"reverse": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`reverse` takes a sequence and produces a new sequence of the same length with all of the same elements as the given sequence but in reverse order.",
		},
		"rsadecrypt": {
			Params: []function.Parameter{
				{
					Name: "ciphertext",
					Type: cty.String,
				},
				{
					Name: "privatekey",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`rsadecrypt` decrypts an RSA-encrypted ciphertext, returning the corresponding cleartext.",
		},
		"setintersection": {
			Params: []function.Parameter{
				{
					Name: "first_set",
					Type: cty.Set(cty.DynamicPseudoType),
				},
			},
			VarParam: &function.Parameter{
				Name: "other_sets",
				Type: cty.Set(cty.DynamicPseudoType),
			},
			ReturnType:  cty.Set(cty.DynamicPseudoType),
			Description: "The `setintersection` function takes multiple sets and produces a single set containing only the elements that all of the given sets have in common. In other words, it computes the [intersection](https://en.wikipedia.org/wiki/Intersection_\\(set_theory\\)) of the sets.",
		},
		"setproduct": {
			VarParam: &function.Parameter{
				Name:        "sets",
				Description: "The sets to consider. Also accepts lists and tuples, and if all arguments are of list or tuple type then the result will preserve the input ordering",
				Type:        cty.DynamicPseudoType,
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "The `setproduct` function finds all of the possible combinations of elements from all of the given sets by computing the [Cartesian product](https://en.wikipedia.org/wiki/Cartesian_product).",
		},
		"setunion": {
			Params: []function.Parameter{
				{
					Name: "first_set",
					Type: cty.Set(cty.DynamicPseudoType),
				},
			},
			VarParam: &function.Parameter{
				Name: "other_sets",
				Type: cty.Set(cty.DynamicPseudoType),
			},
			ReturnType:  cty.Set(cty.DynamicPseudoType),
			Description: "The `setunion` function takes multiple sets and produces a single set containing the elements from all of the given sets. In other words, it computes the [union](https://en.wikipedia.org/wiki/Union_\\(set_theory\\)) of the sets.",
		},
		"sha1": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`sha1` computes the SHA1 hash of a given string and encodes it with hexadecimal digits.",
		},
		"sha256": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`sha256` computes the SHA256 hash of a given string and encodes it with hexadecimal digits.",
		},
		"sha512": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`sha512` computes the SHA512 hash of a given string and encodes it with hexadecimal digits.",
		},
		"signum": {
			Params: []function.Parameter{
				{
					Name: "num",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.Number,
			Description: "`signum` determines the sign of a number, returning a number between -1 and 1 to represent the sign.",
		},
		"slice": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.DynamicPseudoType,
				},
				{
					Name: "start_index",
					Type: cty.Number,
				},
				{
					Name: "end_index",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`slice` extracts some consecutive elements from within a list.",
		},
		"sort": {
			Params: []function.Parameter{
				{
					Name: "list",
					Type: cty.List(cty.String),
				},
			},
			ReturnType:  cty.List(cty.String),
			Description: "`sort` takes a list of strings and returns a new list with those strings sorted lexicographically.",
		},
		"split": {
			Params: []function.Parameter{
				{
					Name: "separator",
					Type: cty.String,
				},
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.List(cty.String),
			Description: "`split` produces a list by dividing a given string at all occurrences of a given separator.",
		},
		"strrev": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`strrev` reverses the characters in a string. Note that the characters are treated as _Unicode characters_ (in technical terms, Unicode [grapheme cluster boundaries](https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries) are respected).",
		},
		"substr": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
				{
					Name: "offset",
					Type: cty.Number,
				},
				{
					Name: "length",
					Type: cty.Number,
				},
			},
			ReturnType:  cty.String,
			Description: "`substr` extracts a substring from a given string by offset and (maximum) length.",
		},
		"templatefile": {
			Params: []function.Parameter{
				{
					Name: "path",
					Type: cty.String,
				},
				{
					Name: "vars",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`templatefile` reads the file at the given path and renders its content as a template using a supplied set of template variables.",
		},
		"timeadd": {
			Params: []function.Parameter{
				{
					Name: "timestamp",
					Type: cty.String,
				},
				{
					Name: "duration",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`timeadd` adds a duration to a timestamp, returning a new timestamp.",
		},
		"timestamp": {
			ReturnType:  cty.String,
			Description: "`timestamp` returns a UTC timestamp string in [RFC 3339](https://tools.ietf.org/html/rfc3339) format.",
		},
		"title": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`title` converts the first letter of each word in the given string to uppercase.",
		},
		"tobool": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Bool,
			Description: "`tobool` converts its argument to a boolean value.",
		},
		"tolist": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.List(cty.DynamicPseudoType),
			Description: "`tolist` converts its argument to a list value.",
		},
		"tomap": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Map(cty.DynamicPseudoType),
			Description: "`tomap` converts its argument to a map value.",
		},
		"tonumber": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Number,
			Description: "`tonumber` converts its argument to a number value.",
		},
		"toset": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.Set(cty.DynamicPseudoType),
			Description: "`toset` converts its argument to a set value.",
		},
		"tostring": {
			Params: []function.Parameter{
				{
					Name: "v",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.String,
			Description: "`tostring` converts its argument to a string value.",
		},
		"transpose": {
			Params: []function.Parameter{
				{
					Name: "values",
					Type: cty.Map(cty.List(cty.String)),
				},
			},
			ReturnType:  cty.Map(cty.List(cty.String)),
			Description: "`transpose` takes a map of lists of strings and swaps the keys and values to produce a new map of lists of strings.",
		},
		"upper": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`upper` converts all cased letters in the given string to uppercase.",
		},
		"urlencode": {
			Params: []function.Parameter{
				{
					Name: "str",
					Type: cty.String,
				},
			},
			ReturnType:  cty.String,
			Description: "`urlencode` applies URL encoding to a given string.",
		},
		"uuid": {
			ReturnType:  cty.String,
			Description: "`uuid` generates a unique identifier string.",
		},
		"values": {
			Params: []function.Parameter{
				{
					Name: "mapping",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`values` takes a map and returns a list containing the values of the elements in that map.",
		},
		"zipmap": {
			Params: []function.Parameter{
				{
					Name: "keys",
					Type: cty.List(cty.String),
				},
				{
					Name: "values",
					Type: cty.DynamicPseudoType,
				},
			},
			ReturnType:  cty.DynamicPseudoType,
			Description: "`zipmap` constructs a map from a list of keys and a corresponding list of values.",
		},
	}
}
