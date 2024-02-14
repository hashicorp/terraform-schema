// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"errors"
	"testing"
)

func TestSchemaMerger_FunctionsForModule_noCoreFunctions(t *testing.T) {
	sm := NewFunctionsMerger(nil)

	_, err := sm.FunctionsForModule(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, coreFunctionsRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

// TODO: add more tests, like in schema_merge_test.go
