// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backend

type CloudData interface {
	Copy() CloudData
	Equals(CloudData) bool
}

type UnknownCloudData struct{}

func (*UnknownCloudData) Copy() CloudData {
	return &UnknownCloudData{}
}

func (*UnknownCloudData) Equals(d CloudData) bool {
	_, ok := d.(*UnknownCloudData)
	return ok
}

type Cloud struct {
	Organization string
	Hostname     string
}

func (r *Cloud) Copy() CloudData {
	return &Cloud{
		Organization: r.Organization,
	}
}

func (r *Cloud) Equals(d CloudData) bool {
	data, ok := d.(*Cloud)
	if !ok {
		return false
	}

	return data.Organization == r.Organization
}
