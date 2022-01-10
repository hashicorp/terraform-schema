package module

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/backend"
)

func TestBackend_Equals(t *testing.T) {
	testCases := []struct {
		first, second *Backend
		expectEqual   bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&Backend{
				Type: "s3",
				Data: &backend.UnknownBackendData{},
			},
			&Backend{
				Type: "s3",
				Data: &backend.UnknownBackendData{},
			},
			true,
		},
		{
			&Backend{
				Type: "s3",
				Data: &backend.UnknownBackendData{},
			},
			&Backend{
				Type: "s4",
				Data: &backend.UnknownBackendData{},
			},
			false,
		},
		{
			&Backend{
				Type: "remote",
				Data: &backend.Remote{},
			},
			&Backend{
				Type: "remote",
				Data: &backend.Remote{},
			},
			true,
		},
		{
			&Backend{
				Type: "remote",
				Data: &backend.Remote{
					Hostname: "foobar",
				},
			},
			&Backend{
				Type: "remote",
				Data: &backend.Remote{
					Hostname: "foobar",
				},
			},
			true,
		},
		{
			&Backend{
				Type: "remote",
				Data: &backend.Remote{
					Hostname: "foobar",
				},
			},
			&Backend{
				Type: "remote",
				Data: &backend.Remote{
					Hostname: "bar",
				},
			},
			false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			equals := tc.first.Equals(tc.second)
			if tc.expectEqual != equals {
				if tc.expectEqual {
					t.Fatalf("expected backends to be equal\nfirst: %#v\nsecond: %#v", tc.first, tc.second)
				}
				t.Fatalf("expected backends to mismatch\nfirst: %#v\nsecond: %#v", tc.first, tc.second)
			}
		})
	}
}

func TestProviderRequirements(t *testing.T) {
	testCases := []struct {
		first, second ProviderRequirements
		expectEqual   bool
	}{
		{
			ProviderRequirements{},
			ProviderRequirements{},
			true,
		},
		{
			ProviderRequirements{
				tfaddr.NewBuiltInProvider("terraform"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			ProviderRequirements{
				tfaddr.NewBuiltInProvider("terraform"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			true,
		},
		{
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			ProviderRequirements{
				tfaddr.NewDefaultProvider("bar"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			false,
		},
		{
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.1")),
			},
			false,
		},
		{
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
				tfaddr.NewDefaultProvider("bar"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			false,
		},
		{
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			ProviderRequirements{
				tfaddr.NewDefaultProvider("foo"): version.MustConstraints(version.NewConstraint("1.0")),
				tfaddr.NewDefaultProvider("bar"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			equals := tc.first.Equals(tc.second)
			if tc.expectEqual != equals {
				if tc.expectEqual {
					t.Fatalf("expected requirements to be equal\nfirst: %#v\nsecond: %#v", tc.first, tc.second)
				}
				t.Fatalf("expected requirements to mismatch\nfirst: %#v\nsecond: %#v", tc.first, tc.second)
			}
		})
	}
}
