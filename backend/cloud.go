// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backend

type Cloud struct {
	Hostname string
}

func (be *Cloud) Equals(b *Cloud) bool {
	if be == nil && b == nil {
		return true
	}

	if be == nil || b == nil {
		return false
	}

	if be.Hostname != b.Hostname {
		return false
	}

	return true
}
