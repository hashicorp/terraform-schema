// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/backend"
	"github.com/hashicorp/terraform-schema/internal/addr"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

type testCase struct {
	name          string
	cfg           string
	expectedMeta  *module.Meta
	expectedError hcl.Diagnostics
}

var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
}

func TestLoadModule(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"empty config",
			``,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"core requirements only",
			`
terraform {
  required_version = "~> 0.12"
}`,
			&module.Meta{
				Path:                 path,
				CoreRequirements:     version.MustConstraints(version.NewConstraint("~> 0.12")),
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"legacy inferred provider requirements",
			`
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

data "blah_foobar" "test" {
  name = "something"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     addr.NewLegacyProvider("aws"),
					{LocalName: "blah"}:    addr.NewLegacyProvider("blah"),
					{LocalName: "google"}:  addr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: addr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					addr.NewLegacyProvider("aws"):     {},
					addr.NewLegacyProvider("blah"):    {},
					addr.NewLegacyProvider("google"):  {},
					addr.NewLegacyProvider("grafana"): {},
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"simplified 0.12 provider requirements",
			`
terraform {
  required_providers {
    aws = "1.2.0"
    google = ">= 3.0.0"
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     addr.NewLegacyProvider("aws"),
					{LocalName: "google"}:  addr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: addr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					addr.NewLegacyProvider("aws"):     version.MustConstraints(version.NewConstraint("1.2.0")),
					addr.NewLegacyProvider("google"):  version.MustConstraints(version.NewConstraint(">= 3.0.0")),
					addr.NewLegacyProvider("grafana"): {},
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"version-only 0.12 provider requirements",
			`
terraform {
  required_providers {
    aws = {
    	version = "1.2.0"
    }
    google = {
    	version = ">= 3.0.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     addr.NewLegacyProvider("aws"),
					{LocalName: "google"}:  addr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: addr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					addr.NewLegacyProvider("aws"):     version.MustConstraints(version.NewConstraint("1.2.0")),
					addr.NewLegacyProvider("google"):  version.MustConstraints(version.NewConstraint(">= 3.0.0")),
					addr.NewLegacyProvider("grafana"): {},
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"0.13+ provider requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "2.1.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
					{LocalName: "grafana"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "grafana",
						Type:      "grafana",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: version.MustConstraints(version.NewConstraint("1.0.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "grafana",
						Type:      "grafana",
					}: version.MustConstraints(version.NewConstraint("2.1.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"multiple valid version requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}

terraform {
  required_providers {
    aws = {
      version = "1.1.0"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: version.MustConstraints(version.NewConstraint(">= 1.0.0,1.1.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"multiple invalid version requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}

terraform {
  required_providers {
    aws = {
    	source = "hashicorp/aws"
      version = "1.1.0"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: version.MustConstraints(version.NewConstraint(">= 1.0.0,1.1.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"0.13+ provider aliases",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}
provider "aws" {
  alias = "euwest"
  region = "eu-west-2"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "euwest"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: version.MustConstraints(version.NewConstraint("1.0.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"0.15+ provider aliases",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
      configuration_aliases = [aws.east]
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}
provider "aws" {
  alias = "west"
  region = "eu-west-2"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "east"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "west"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: version.MustConstraints(version.NewConstraint("1.0.0")),
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"explicit provider association",
			`
terraform {
  required_providers {
    goo = {
      source  = "hashicorp/google-beta"
      version = "2.0.0"
    }
  }
}

resource "google_something" "test" {
	provider = goo
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "goo"}: {
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google-beta",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultProviderRegistryHost,
						Namespace: "hashicorp",
						Type:      "google-beta",
					}: version.MustConstraints(version.NewConstraint("2.0.0")),
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_Variables(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"no name variables",
			`
variable "" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"no name variables",
			`
variable {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing name for variable",
					Detail:   "All variable blocks must have 1 labels (name).",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 10,
							Byte:   10,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 11,
							Byte:   11,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 11,
							Byte:   11,
						},
					},
				},
			},
		},
		{
			"double label variables",
			`
variable "one" "two" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Extraneous label for variable",
					Detail:   "Only 1 labels (name) are expected for variable blocks.",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 16,
							Byte:   16,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 21,
							Byte:   21,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 23,
							Byte:   23,
						},
					},
				},
			},
		},
		{
			"empty variables",
			`
variable "name" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type: cty.DynamicPseudoType,
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with type",
			`
variable "name" {
	type = string
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type: cty.String,
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with description",
			`
variable "name" {
	description = "description"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.DynamicPseudoType,
						Description: "description",
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with sensitive",
			`
variable "name" {
	sensitive = true
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.DynamicPseudoType,
						IsSensitive: true,
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with type and description and sensitive",
			`
variable "name" {
	type = string
	description = "description"
	sensitive = true
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.String,
						Description: "description",
						IsSensitive: true,
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with default",
			`
variable "name" {
  	default = {}
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:         cty.DynamicPseudoType,
						DefaultValue: cty.EmptyObjectVal,
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"variables with optional type values",
			`
variable "name" {
  type = object({
    foo = optional(string, "food")
    bar = optional(number)
  })
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type: cty.Object(map[string]cty.Type{
							"foo": cty.String,
							"bar": cty.Number,
						}),
						TypeDefaults: &typeexpr.Defaults{
							Type: cty.Object(map[string]cty.Type{
								"foo": cty.String,
								"bar": cty.Number,
							}),
							DefaultValues: map[string]cty.Value{
								"foo": cty.StringVal("food"),
							},
						},
					},
				},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"empty output",
			`
output "name" {
}
`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs: map[string]module.Output{
					"name": {Value: cty.NilVal},
				},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_backend(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"no backend",
			`
terraform {

}`,
			&module.Meta{
				Path:                 path,
				Backend:              nil,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"s3 backend",
			`
terraform {
  backend "s3" {
  	blah = "test"
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "s3",
					Data: &backend.UnknownBackendData{},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"empty remote backend",
			`
terraform {
  backend "remote" {}
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"remote backend with hostname",
			`
terraform {
  backend "remote" {
  	hostname = "app.terraform.io"
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{Hostname: "app.terraform.io"},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"remote backend with hostname and more attributes",
			`
terraform {
  backend "remote" {
    hostname = "app.terraform.io"
    organization = "test"

    workspaces {
      name = "test"
    }
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{Hostname: "app.terraform.io"},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_cloud(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"cloud backend",
			`
terraform {
	cloud {
		hostname = "app.terraform.io"
		organization = "example_corp"
	}
}`,
			&module.Meta{
				Path:    path,
				Backend: nil,
				Cloud: &backend.Cloud{
					Hostname: "app.terraform.io",
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"cloud backend empy hostname",
			`
terraform {
	cloud {
		organization = "example_corp"
	}
}`,
			&module.Meta{
				Path:                 path,
				Backend:              nil,
				Cloud:                &backend.Cloud{},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"additional block",
			`
terraform {
	cloud {
		hostname = "foo.com"
		workspaces {
			tags = ["app"]
		}
	}
}`,
			&module.Meta{
				Path:    path,
				Backend: nil,
				Cloud: &backend.Cloud{
					Hostname: "foo.com",
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_Modules(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"no name module",
			`
module "" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			nil,
		},
		{
			"no name modules",
			`
module {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing name for module",
					Detail:   "All module blocks must have 1 labels (name).",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 8,
							Byte:   8,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 9,
							Byte:   9,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 9,
							Byte:   9,
						},
					},
				},
			},
		},
		{
			"double label modules",
			`
module "one" "two" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls:          map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Extraneous label for module",
					Detail:   "Only 1 labels (name) are expected for module blocks.",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 14,
							Byte:   14,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 19,
							Byte:   19,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 21,
							Byte:   21,
						},
					},
				},
			},
		},
		{
			"empty modules",
			`
module "name" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 3, Column: 2, Byte: 18},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with source",
			`
module "name" {
	source = "registry.terraform.io/terraform-aws-modules/vpc/aws"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						SourceAddr: tfaddr.MustParseModuleSource("registry.terraform.io/terraform-aws-modules/vpc/aws"),
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 4, Column: 2, Byte: 82},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with version",
			`
module "name" {
	version = "> 3.0.0, < 4.0.0"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						Version:    version.MustConstraints(version.NewConstraint("> 3.0.0, < 4.0.0")),
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 4, Column: 2, Byte: 48},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with source and version",
			`
module "name" {
	source = "terraform-aws-modules/vpc/aws"
	version = "1.0.0"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						SourceAddr: tfaddr.MustParseModuleSource("terraform-aws-modules/vpc/aws"),
						Version:    version.MustConstraints(version.NewConstraint("1.0.0")),
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 5, Column: 2, Byte: 79},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with local source",
			`
module "name" {
	source = "./local"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						SourceAddr: module.LocalSourceAddr("./local"),
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 4, Column: 2, Byte: 38},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with local source and inputs",
			`
module "name" {
	source = "./local"
	one = "one"
	two = 42
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						SourceAddr: module.LocalSourceAddr("./local"),
						InputNames: []string{
							"one", "two",
						},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 6, Column: 2, Byte: 61},
						},
					},
				},
			},
			nil,
		},
		{
			"modules with unknown source",
			`
module "name" {
	source = "github.com/terraform-aws-modules/terraform-aws-security-group"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
				Filenames:            []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{
					"name": {
						LocalName:  "name",
						SourceAddr: module.UnknownSourceAddr("github.com/terraform-aws-modules/terraform-aws-security-group"),
						InputNames: []string{},
						RangePtr: &hcl.Range{
							Filename: "test.tf",
							Start:    hcl.Pos{Line: 2, Column: 15, Byte: 15},
							End:      hcl.Pos{Line: 4, Column: 2, Byte: 92},
						},
					},
				},
			},
			nil,
		},
		{
			"invalid provider name",
			`
provider "-" {
}
provider "valid" {
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "valid"}: addr.NewLegacyProvider("valid"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					addr.NewLegacyProvider("valid"): {},
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid provider name",
					Detail:   `"-" is not a valid provider name: must contain only letters, digits, and dashes, and may not use leading or trailing dashes`,
				},
			},
		},
		{
			"invalid implied provider name",
			`
resource "-invalid_foo" "name" {
}
resource "valid_foo" "name" {
}
data "-invalid_bar" "name" {
}
data "valid_bar" "name" {
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "valid"}: addr.NewLegacyProvider("valid"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					addr.NewLegacyProvider("valid"): {},
				},
				Variables:   map[string]module.Variable{},
				Outputs:     map[string]module.Output{},
				Filenames:   []string{"test.tf"},
				ModuleCalls: map[string]module.DeclaredModuleCall{},
			},
			hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid provider name",
					Detail:   `"-invalid" is not a valid implied provider name: must contain only letters, digits, and dashes, and may not use leading or trailing dashes`,
				},
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid provider name",
					Detail:   `"-invalid" is not a valid implied provider name: must contain only letters, digits, and dashes, and may not use leading or trailing dashes`,
				},
			},
		},
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), "test.tf", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				"test.tf": f,
			}

			meta, diags := LoadModule(path, files)

			if diff := cmp.Diff(tc.expectedError, diags, customComparer...); diff != "" {
				t.Fatalf("expected errors doesn't match: %s", diff)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("module meta doesn't match: %s", diff)
			}
		})
	}
}

func compareVersionConstraint(x, y *version.Constraint) bool {
	return x.Equals(y)
}
