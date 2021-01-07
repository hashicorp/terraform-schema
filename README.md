# terraform-schema [![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/terraform-schema.svg)](https://pkg.go.dev/github.com/hashicorp/terraform-schema)

This library helps assembling a complete [`hcl-lang`](https://github.com/hashicorp/hcl-lang)
schema for decoding Terraform config based on static Terraform core schema
and relevant provider schemas.

**There is more than one schema?**

Yes.

 - Terraform Core defines top-level schema
   - `provider`, `resource` or `data` blocks incl. meta attributes, such as `alias` or `count`
   - `variable`, `output` blocks etc.
 - each Terraform provider defines its own schema for the body of some of these blocks
   - attributes and nested blocks inside `resource`, `data` or `provider` blocks

Each of these can also differ between (core / provider) version.

## Current Status

This project is in use by the Terraform Language Server and could _in theory_
be used by other projects which need to decode the _whole_ configuration.

However it has not been tested in any other scenarios.

Please note that this library depends on [`hcl-lang`](https://github.com/hashicorp/hcl-lang)
which itself is not considered stable yet.

**Breaking changes may be introduced.**

## Alternative Solution

If you only need to perform shallow config decoding, e.g. you just need to
get a list of variables, outputs, provider names/source etc. and you don't care
as much about version-specific or provider-specific details then you should
explore [terraform-config-inspect](https://github.com/hashicorp/terraform-config-inspect)
instead, which will likely be sufficient for your needs.

## How It Works

### Usage

```go
import (
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-json"
)

// parse files e.g. via hclsyntax
parsedFiles := map[string]*hcl.File{ /* ... */ }

// obtain relevant core schema
coreSchema := tfschema.UniversalCoreModuleSchema()

// obtain relevant provider schemas e.g. via terraform-exec
// and marshal them into terraform-json type
providerSchemas := &tfjson.ProviderSchemas{ /* ... */ }

mergedSchema, err := tfschema.MergeCoreWithJsonProviderSchemas(parsedFiles, coreSchema, providerSchemas)
if err != nil {
	// ...
}

```

### Provider Schemas

The only reliable way of obtaining provider schemas at the time of writing is via
Terraform CLI by running `terraform providers schema -json` (0.12+).

[`terraform-exec`](https://github.com/hashicorp/terraform-exec) can help automating
the process of obtaining the schema.

[`terraform-json`](https://github.com/hashicorp/terraform-json) provides types
that the JSON output can be marshalled into, which also used by `terraform-exec`
and is considered as standard way of representing the output.


#### Known Issues

At the time of writing there is a known issue affecting the above command
where it requires the following to be true in order to produce schemas:

 - configuration is valid (e.g. contains no incomplete blocks)
 - authentication is provided for any remote backend

Read more at [hashicorp/terraform#24261](https://github.com/hashicorp/terraform/issues/24261).

Other ways of obtaining schemas are also being explored.

## Why a Separate Repository (from Terraform Core)?

As demonstrated by the [recent separation of Plugin SDK](https://www.terraform.io/docs/extend/guides/v1-upgrade-guide.html),
Terraform Core is not intended to be consumed as a Go module
and any packages which happen to be importable were not tested
nor designed for use outside the context of Terraform Core.

Terraform Core's versioning reflects mainly the experience of end-users
interacting with Terraform via the CLI. Functionality which is consumed
as Go package and imported into other Go programs deserves its own
dedicated versioning.

Terraform Core's schema internally always supports the latest version
(e.g. Terraform 0.12 can parse 0.12 configuration). Consumers of this library
however need to parse different versions of configuration at the same time.

## Experimental Status

By using the software in this repository (the "Software"), you acknowledge that: (1) the Software is still in development, may change, and has not been released as a commercial product by HashiCorp and is not currently supported in any way by HashiCorp; (2) the Software is provided on an "as-is" basis, and may include bugs, errors, or other issues; (3) the Software is NOT INTENDED FOR PRODUCTION USE, use of the Software may result in unexpected results, loss of data, or other unexpected results, and HashiCorp disclaims any and all liability resulting from use of the Software; and (4) HashiCorp reserves all rights to make all decisions about the features, functionality and commercial release (or non-release) of the Software, at any time and without any obligation or liability whatsoever.

