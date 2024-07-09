package stack

import "github.com/hashicorp/go-version"

type Component struct {
	Source  string
	Version version.Constraints
}
