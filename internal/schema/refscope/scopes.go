// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package refscope

import (
	"github.com/hashicorp/hcl-lang/lang"
)

var (
	BuiltinScope   = lang.ScopeId("builtin")
	DataScope      = lang.ScopeId("data")
	LocalScope     = lang.ScopeId("local")
	ModuleScope    = lang.ScopeId("module")
	OutputScope    = lang.ScopeId("output")
	ProviderScope  = lang.ScopeId("provider")
	ResourceScope  = lang.ScopeId("resource")
	EphemeralScope = lang.ScopeId("ephemeral")
	VariableScope  = lang.ScopeId("variable")

	ComponentScope     = lang.ScopeId("component")
	IdentityTokenScope = lang.ScopeId("identity_token")
	StoreScope         = lang.ScopeId("store")
	OrchestrateContext = lang.ScopeId("orchestrate_context")

	ListScope = lang.ScopeId("list")

	ActionScope = lang.ScopeId("action")
)
