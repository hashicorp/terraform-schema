// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import "github.com/hashicorp/hcl-lang/reference"

type referenceTag string

func (m referenceTag) GoString() string {
	return "tags." + string(m)
}

// Has returns true if the reference target has the given tag.
func Has(target reference.Target, tag referenceTag) bool {
	// TODO: Make first argument less specific?
	_, ok := target.Tags[tag]

	return ok
}

// Ephemeral indicates that this reference is marked as ephemeral
const Ephemeral = referenceTag("Ephemeral")
