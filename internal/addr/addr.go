// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package addr

import (
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

// NewLegacyProvider returns a mock address for a provider.
func NewLegacyProvider(name string) tfaddr.Provider {
	return tfaddr.Provider{
		Type:      tfaddr.MustParseProviderPart(name),
		Namespace: tfaddr.LegacyProviderNamespace,
		Hostname:  tfaddr.DefaultProviderRegistryHost,
	}
}

// NewDefaultProvider returns the default address of a HashiCorp-maintained,
// Registry-hosted provider.
func NewDefaultProvider(name string) tfaddr.Provider {
	return tfaddr.Provider{
		Type:      tfaddr.MustParseProviderPart(name),
		Namespace: "hashicorp",
		Hostname:  tfaddr.DefaultProviderRegistryHost,
	}
}

// NewBuiltInProvider returns the address of a "built-in" provider. See
// the docs for Provider.IsBuiltIn for more information.
func NewBuiltInProvider(name string) tfaddr.Provider {
	return tfaddr.Provider{
		Type:      tfaddr.MustParseProviderPart(name),
		Namespace: tfaddr.BuiltInProviderNamespace,
		Hostname:  tfaddr.BuiltInProviderHost,
	}
}
