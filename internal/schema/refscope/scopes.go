// Copyright IBM Corp. 2020, 2026
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
	InputScope     = lang.ScopeId("input")

	ComponentScope             = lang.ScopeId("component")
	StackScope                 = lang.ScopeId("stack")
	IdentityTokenScope         = lang.ScopeId("identity_token")
	StoreScope                 = lang.ScopeId("store")
	OrchestrateContext         = lang.ScopeId("orchestrate_context")
	UpstreamInputScope         = lang.ScopeId("upstream_input")
	DeploymentAutoApproveScope = lang.ScopeId("deployment_auto_approve")
	DeploymentGroupScope       = lang.ScopeId("deployment_group")
	DeploymentScope            = lang.ScopeId("deployment")

	ListScope = lang.ScopeId("list")

	ActionScope = lang.ScopeId("action")

	ResourcePolicyScope = lang.ScopeId("resource_policy")
	ProviderPolicyScope = lang.ScopeId("provider_policy")
	ModulePolicyScope   = lang.ScopeId("module_policy")
)
