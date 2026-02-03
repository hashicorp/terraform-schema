// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
)

func TestGetTerraformReleases(t *testing.T) {
	releases, err := GetTerraformReleases()
	if err != nil {
		t.Fatal(err)
	}

	minExpectedLength := 234
	if len(releases) < minExpectedLength {
		t.Fatalf("expected >= %d releases, %d given", minExpectedLength, len(releases))
	}

	// The oldest release should really be 0.1.0. We're however getting
	// releases sorted by dates and those dates were backfilled as part
	// of some older data migrations where the original dates were lost.
	expectedDate := time.Date(2017, 3, 1, 17, 36, 49, 0, time.UTC)
	expectedOldestRelease := release{
		Version: version.Must(version.NewVersion("0.6.4")),
		Created: &expectedDate,
	}
	oldestRelease := releases[len(releases)-1]
	if diff := cmp.Diff(expectedOldestRelease, oldestRelease); diff != "" {
		t.Fatalf("unexpected oldest release: %s", diff)
	}
}
