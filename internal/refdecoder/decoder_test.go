package refdecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/internal/addrs"
)

func TestDecodeProviderReferences(t *testing.T) {
	testCases := []struct {
		name         string
		src          string
		expectedRefs addrs.ProviderReferences
	}{
		{
			"provider block",
			`
provider "aws" {

}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "aws",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "aws",
				},
			},
		},
		{
			"aliased provider block",
			`
provider "blablah" {
	alias = "foo"
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "blablah",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "blablah",
				},
				addrs.LocalProviderConfig{
					LocalName: "blablah",
					Alias:     "foo",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "blablah",
				},
			},
		},
		{
			"terraform block",
			`
terraform {
  required_providers {
    mycloud = {
      source  = "mycorp/mycloud"
      version = "~> 1.0"
    }
  }
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "mycloud",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "mycorp",
					Type:      "mycloud",
				},
			},
		},
		{
			"default namespace provider under different local name",
			`
terraform {
  required_providers {
    rand = {
      source  = "hashicorp/random"
      version = "3.0.0"
    }
  }
}

resource "random_string" "name" {
  provider = rand
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "rand",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "random",
				},
			},
		},
		{
			"resource block",
			`
resource "mycloud_instance" "foo" {
	count = 2
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "mycloud",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "mycloud",
				},
			},
		},
		{
			"data block",
			`
data "mycloud_instance" "foo" {
	count = 2
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "mycloud",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "hashicorp",
					Type:      "mycloud",
				},
			},
		},
		{
			"configuration aliases",
			`
terraform {
  required_providers {
    xyz = {
      source                = "grafana/grafana"
      version               = "1.9.0"
      configuration_aliases = [xyz.blah]
    }
  }
}
`,
			addrs.ProviderReferences{
				addrs.LocalProviderConfig{
					LocalName: "xyz",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "grafana",
					Type:      "grafana",
				},
				addrs.LocalProviderConfig{
					LocalName: "xyz",
					Alias:     "blah",
				}: addrs.Provider{
					Hostname:  addrs.DefaultRegistryHost,
					Namespace: "grafana",
					Type:      "grafana",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.src), "test.tf", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}

			files := map[string]*hcl.File{
				"test.tf": f,
			}

			refs, diags := DecodeProviderReferences(files)
			if diff := cmp.Diff(tc.expectedRefs, refs); diff != "" {
				t.Fatalf("unexpected provider references: %s", diff)
			}
		})
	}
}
