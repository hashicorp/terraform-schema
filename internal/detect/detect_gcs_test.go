// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"testing"
)

func TestGCSDetector(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			"www.googleapis.com/storage/v1/bucket/foo",
			"gcs::https://www.googleapis.com/storage/v1/bucket/foo",
		},
		{
			"www.googleapis.com/storage/v1/bucket/foo/bar",
			"gcs::https://www.googleapis.com/storage/v1/bucket/foo/bar",
		},
		{
			"www.googleapis.com/storage/v1/foo/bar.baz",
			"gcs::https://www.googleapis.com/storage/v1/foo/bar.baz",
		},
	}

	f := new(GCSDetector)
	for i, tc := range cases {
		output, ok, err := f.Detect(tc.Input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if !ok {
			t.Fatal("not ok")
		}

		if output != tc.Output {
			t.Fatalf("%d: bad: %#v", i, output)
		}
	}
}
