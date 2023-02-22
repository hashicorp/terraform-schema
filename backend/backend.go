// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backend

type BackendData interface {
	Copy() BackendData
	Equals(BackendData) bool
}

type UnknownBackendData struct{}

func (*UnknownBackendData) Copy() BackendData {
	return &UnknownBackendData{}
}

func (*UnknownBackendData) Equals(d BackendData) bool {
	_, ok := d.(*UnknownBackendData)
	return ok
}

type Remote struct {
	Hostname string
}

func (r *Remote) Copy() BackendData {
	return &Remote{
		Hostname: r.Hostname,
	}
}

func (r *Remote) Equals(d BackendData) bool {
	data, ok := d.(*Remote)
	if !ok {
		return false
	}

	return data.Hostname == r.Hostname
}
