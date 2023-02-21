// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type ProviderSchema struct {
	Provider    *schema.BodySchema
	Resources   map[string]*schema.BodySchema
	DataSources map[string]*schema.BodySchema
}

func (ps *ProviderSchema) Copy() *ProviderSchema {
	if ps == nil {
		return nil
	}

	newPs := &ProviderSchema{
		Provider: ps.Provider.Copy(),
	}

	if ps.Resources != nil {
		newPs.Resources = make(map[string]*schema.BodySchema, len(ps.Resources))
		for name, rSchema := range ps.Resources {
			newPs.Resources[name] = rSchema.Copy()
		}
	}

	if ps.DataSources != nil {
		newPs.DataSources = make(map[string]*schema.BodySchema, len(ps.DataSources))
		for name, rSchema := range ps.DataSources {
			newPs.DataSources[name] = rSchema.Copy()
		}
	}

	return newPs
}

func (ps *ProviderSchema) SetProviderVersion(pAddr tfaddr.Provider, v *version.Version) {
	if ps.Provider != nil {
		ps.Provider.Detail = detailForSrcAddr(pAddr, v)
		ps.Provider.HoverURL = urlForProvider(pAddr, v)
		ps.Provider.DocsLink = docsLinkForProvider(pAddr, v)
	}
	for _, rSchema := range ps.Resources {
		rSchema.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, dsSchema := range ps.DataSources {
		dsSchema.Detail = detailForSrcAddr(pAddr, v)
	}
}
