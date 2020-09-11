package addrs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestParseProviderConfigCompact_empty(t *testing.T) {
	lProviderCfg, err := ParseProviderConfigCompact(nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := LocalProviderConfig{}
	if diff := cmp.Diff(expected, lProviderCfg); diff != "" {
		t.Fatalf("mismatch: %s", diff)
	}
}

func TestParseProviderConfigCompact_nameOnly(t *testing.T) {
	lProviderCfg, err := parseProviderConfigCompactStr("justname")
	if err != nil {
		t.Fatal(err)
	}
	expected := LocalProviderConfig{LocalName: "justname"}
	if diff := cmp.Diff(expected, lProviderCfg); diff != "" {
		t.Fatalf("mismatch: %s", diff)
	}
}

func TestParseProviderConfigCompact_fullRef(t *testing.T) {
	lProviderCfg, err := parseProviderConfigCompactStr("aws.uswest")
	if err != nil {
		t.Fatal(err)
	}
	expected := LocalProviderConfig{LocalName: "aws", Alias: "uswest"}
	if diff := cmp.Diff(expected, lProviderCfg); diff != "" {
		t.Fatalf("mismatch: %s", diff)
	}
}

func parseProviderConfigCompactStr(str string) (LocalProviderConfig, error) {
	traversal, parseDiags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	if parseDiags.HasErrors() {
		return LocalProviderConfig{}, parseDiags
	}

	return ParseProviderConfigCompact(traversal)
}
