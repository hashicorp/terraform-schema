package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestProviderSchema_SetProviderVersion(t *testing.T) {
	ps := &ProviderSchema{
		Provider: &schema.BodySchema{},
		Resources: map[string]*schema.BodySchema{
			"foo": {
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsOptional: true,
					},
				},
			},
		},
		DataSources: map[string]*schema.BodySchema{
			"bar": {
				Attributes: map[string]*schema.AttributeSchema{
					"num": {
						Expr:       schema.LiteralTypeOnly(cty.Number),
						IsOptional: true,
					},
				},
			},
		},
	}
	expectedSchema := &ProviderSchema{
		Provider: &schema.BodySchema{
			Detail:   "hashicorp/aws 1.2.5",
			HoverURL: "https://registry.terraform.io/providers/hashicorp/aws/1.2.5/docs",
			DocsLink: &schema.DocsLink{
				URL:     "https://registry.terraform.io/providers/hashicorp/aws/1.2.5/docs",
				Tooltip: "hashicorp/aws Documentation",
			},
		},
		Resources: map[string]*schema.BodySchema{
			"foo": {
				Detail: "hashicorp/aws 1.2.5",
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Expr:       schema.LiteralTypeOnly(cty.String),
						IsOptional: true,
					},
				},
			},
		},
		DataSources: map[string]*schema.BodySchema{
			"bar": {
				Detail: "hashicorp/aws 1.2.5",
				Attributes: map[string]*schema.AttributeSchema{
					"num": {
						Expr:       schema.LiteralTypeOnly(cty.Number),
						IsOptional: true,
					},
				},
			},
		},
	}

	pAddr := tfaddr.Provider{
		Hostname:  tfaddr.DefaultRegistryHost,
		Namespace: "hashicorp",
		Type:      "aws",
	}
	pv := version.Must(version.NewVersion("1.2.5"))

	ps.SetProviderVersion(pAddr, pv)

	if diff := cmp.Diff(expectedSchema, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("unexpected schema: %s", diff)
	}
}
