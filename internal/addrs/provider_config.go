package addrs

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
)

// LocalProviderConfig is the address of a provider configuration from the
// perspective of references in a particular module.
//
// Finding the corresponding AbsProviderConfig will require looking up the
// LocalName in the providers table in the module's configuration; there is
// no syntax-only translation between these types.
type LocalProviderConfig struct {
	LocalName string

	// If not empty, Alias identifies which non-default (aliased) provider
	// configuration this address refers to.
	Alias string
}

// ParseProviderConfigCompact parses the given absolute traversal as a relative
// provider address in compact form. The following are examples of traversals
// that can be successfully parsed as compact relative provider configuration
// addresses:
//
//     aws
//     aws.foo
//
// This function will panic if given a relative traversal.
//
// If the returned diagnostics contains errors then the result value is invalid
// and must not be used.
func ParseProviderConfigCompact(traversal hcl.Traversal) (LocalProviderConfig, error) {
	var errs *multierror.Error

	if len(traversal) == 0 {
		return LocalProviderConfig{}, nil
	}

	ret := LocalProviderConfig{
		LocalName: traversal.RootName(),
	}

	if len(traversal) < 2 {
		// Just a type name, then.
		return ret, errs.ErrorOrNil()
	}

	aliasStep := traversal[1]
	switch ts := aliasStep.(type) {
	case hcl.TraverseAttr:
		ret.Alias = ts.Name
		return ret, errs.ErrorOrNil()
	default:
		errs = multierror.Append(&ParserError{
			Summary: "Invalid provider configuration address",
			Detail:  "The provider type name must either stand alone or be followed by an alias name separated with a dot.",
		})
	}

	if len(traversal) > 2 {
		errs = multierror.Append(&ParserError{
			Summary: "Invalid provider configuration address",
			Detail:  "Extraneous extra operators after provider configuration address.",
		})
	}

	return ret, errs.ErrorOrNil()
}
