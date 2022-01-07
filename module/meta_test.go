package module

import (
	"fmt"
	"testing"

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
