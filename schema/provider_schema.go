// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type ProviderSchema struct {
	Provider           *schema.BodySchema
	Resources          map[string]*schema.BodySchema
	EphemeralResources map[string]*schema.BodySchema
	DataSources        map[string]*schema.BodySchema
	Functions          map[string]*schema.FunctionSignature
	ListResources      map[string]*schema.BodySchema
	ActionResources    map[string]*schema.BodySchema
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

	if ps.EphemeralResources != nil {
		newPs.EphemeralResources = make(map[string]*schema.BodySchema, len(ps.EphemeralResources))
		for name, erSchema := range ps.EphemeralResources {
			newPs.EphemeralResources[name] = erSchema.Copy()
		}
	}

	if ps.DataSources != nil {
		newPs.DataSources = make(map[string]*schema.BodySchema, len(ps.DataSources))
		for name, rSchema := range ps.DataSources {
			newPs.DataSources[name] = rSchema.Copy()
		}
	}

	if ps.Functions != nil {
		newPs.Functions = make(map[string]*schema.FunctionSignature, len(ps.Functions))
		for name, fSig := range ps.Functions {
			newPs.Functions[name] = fSig.Copy()
		}
	}

	if ps.ListResources != nil {
		newPs.ListResources = make(map[string]*schema.BodySchema, len(ps.ListResources))
		for name, lsSchema := range ps.ListResources {
			newPs.ListResources[name] = lsSchema.Copy()
		}
	}

	if ps.ActionResources != nil {
		newPs.ActionResources = make(map[string]*schema.BodySchema, len(ps.ActionResources))
		for name, arSchema := range ps.ActionResources {
			newPs.ActionResources[name] = arSchema.Copy()
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
	for _, erSchema := range ps.EphemeralResources {
		erSchema.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, dsSchema := range ps.DataSources {
		dsSchema.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, fSig := range ps.Functions {
		fSig.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, lsSchema := range ps.ListResources {
		lsSchema.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, arSchema := range ps.ActionResources {
		arSchema.Detail = detailForSrcAddr(pAddr, v)
	}
}
