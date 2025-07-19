// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/hcl-lang/lang"
)

func TestImportBlock(t *testing.T) {
	schema := importBlock()

	// Test basic schema properties
	if schema.Description != lang.Markdown("Import resources into Terraform to bring them under Terraform's management") {
		t.Errorf("unexpected description: %v", schema.Description)
	}

	if schema.Body.HoverURL != "https://developer.hashicorp.com/terraform/language/import" {
		t.Errorf("unexpected hover URL: %v", schema.Body.HoverURL)
	}

	// Test that all expected attributes are present
	expectedAttrs := []string{"provider", "id", "identity", "to"}
	for _, attr := range expectedAttrs {
		if _, ok := schema.Body.Attributes[attr]; !ok {
			t.Errorf("missing expected attribute: %s", attr)
		}
	}

	// Test id attribute
	idAttr := schema.Body.Attributes["id"]
	if !idAttr.IsOptional {
		t.Error("id attribute should be optional")
	}
	if idAttr.Description != lang.Markdown("ID of the resource to be imported. e.g. `i-abcd1234`. Either `id` or `identity` must be specified, but not both.") {
		t.Errorf("unexpected id description: %v", idAttr.Description)
	}

	// Test identity attribute
	identityAttr := schema.Body.Attributes["identity"]
	if !identityAttr.IsOptional {
		t.Error("identity attribute should be optional")
	}
	if identityAttr.Description != lang.Markdown("Key-value pairs to identify the resource to be imported. Either `id` or `identity` must be specified, but not both.") {
		t.Errorf("unexpected identity description: %v", identityAttr.Description)
	}
}

func TestImportBlock_Attributes(t *testing.T) {
	schema := importBlock()

	// Test that all expected attributes are present
	expectedAttrs := []string{"provider", "id", "identity", "to"}
	for _, attr := range expectedAttrs {
		if _, ok := schema.Body.Attributes[attr]; !ok {
			t.Errorf("missing expected attribute: %s", attr)
		}
	}

	// Test that id and identity are optional
	if !schema.Body.Attributes["id"].IsOptional {
		t.Error("id attribute should be optional")
	}
	if !schema.Body.Attributes["identity"].IsOptional {
		t.Error("identity attribute should be optional")
	}

	// Test that to is required
	if !schema.Body.Attributes["to"].IsRequired {
		t.Error("to attribute should be required")
	}

	// Test that provider is optional
	if !schema.Body.Attributes["provider"].IsOptional {
		t.Error("provider attribute should be optional")
	}
}
