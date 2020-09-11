package addrs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProviderReferences_LocalNameByAddr(t *testing.T) {
	ref := LocalProviderConfig{LocalName: "customname"}
	addr := Provider{
		Type:      "aws",
		Hostname:  "registry.terraform.io",
		Namespace: "hashicorp",
	}
	refs := ProviderReferences{ref: addr}

	foundRefs := refs.LocalNamesByAddr(addr)
	if len(foundRefs) == 0 {
		t.Fatal("expected to find the reference")
	}
	if diff := cmp.Diff(ref, foundRefs[0]); diff != "" {
		t.Fatalf("reference mismatch: %s", diff)
	}
}
