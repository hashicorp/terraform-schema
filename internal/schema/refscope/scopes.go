// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package refscope

import (
	"github.com/hashicorp/hcl-lang/lang"
)

var (
	BuiltinScope  = lang.ScopeId("builtin")
	DataScope     = lang.ScopeId("data")
	LocalScope    = lang.ScopeId("local")
	ModuleScope   = lang.ScopeId("module")
	OutputScope   = lang.ScopeId("output")
	ProviderScope = lang.ScopeId("provider")
	ResourceScope = lang.ScopeId("resource")
	VariableScope = lang.ScopeId("variable")
)
