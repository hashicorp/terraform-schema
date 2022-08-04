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

	expectedLength := 234
	if expectedLength < len(releases) {
		t.Fatalf("expected >= %d releases, %d given", expectedLength, len(releases))
	}

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
