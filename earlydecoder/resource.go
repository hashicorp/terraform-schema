// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-schema/module"
)

type resource struct {
	Type     string
	Name     string
	Provider module.ProviderRef
}

// MapKey returns a string that can be used to uniquely identify the receiver
// in a map[string]*resource.
func (r *resource) MapKey() string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

func inferProviderNameFromType(typeName string) string {
	if underPos := strings.IndexByte(typeName, '_'); underPos != -1 {
		return typeName[:underPos]
	}
	return typeName
}
