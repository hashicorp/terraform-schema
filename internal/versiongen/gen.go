package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-version"
)

var baseURL = "https://api.releases.hashicorp.com/v1"

type release struct {
	Version *version.Version `json:"version"`
	Created *time.Time       `json:"timestamp_created"`
}

func main() {
	var writePath string
	flag.StringVar(&writePath, "w", "", "Path to write to")
	flag.Parse()

	output := os.Stdout
	if writePath != "" {
		f, err := os.OpenFile(writePath, os.O_RDWR|os.O_CREATE, 0o755)
		if err != nil {
			log.Fatal(err)
		}
		output = f
	}

	releases, err := GetTerraformReleases()
	if err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].Version.GreaterThan(releases[j].Version)
	})

	outputTpl := `// Code generated by "versiongen"; DO NOT EDIT.
package schema

import (
	"github.com/hashicorp/go-version"
)

var (
	OldestAvailableVersion = version.Must(version.NewVersion("{{ .OldestVersion }}"))
	LatestAvailableVersion = version.Must(version.NewVersion("{{ .LatestVersion }}"))

	terraformVersions = version.Collection{
{{- range .Releases }}
		version.Must(version.NewVersion("{{ .Version }}")),
{{- end }}
	}
)
`
	tpl, err := template.New("output").Parse(outputTpl)
	if err != nil {
		log.Fatal(err)
	}

	type data struct {
		Releases      []release
		OldestVersion *version.Version
		LatestVersion *version.Version
	}

	// we keep this hard-coded to 0.12 since
	// we don't have schema for older versions
	oldestVersion := version.Must(version.NewVersion("0.12.0"))

	err = tpl.Execute(output, data{
		Releases:      releases,
		LatestVersion: releases[0].Version,
		OldestVersion: oldestVersion,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetTerraformReleases() ([]release, error) {
	releases := make([]release, 0)

	var after *time.Time
	for {
		r, err := getTerraformReleasesAfter(after)
		if err != nil {
			return releases, err
		}
		if len(r) == 0 {
			break
		}

		releases = append(releases, r...)
		after = r[len(r)-1].Created
	}

	return releases, nil
}

func getTerraformReleasesAfter(after *time.Time) ([]release, error) {
	u, err := url.Parse(fmt.Sprintf("%s/releases/%s", baseURL, "terraform"))
	if err != nil {
		return nil, err
	}

	params := u.Query()
	params.Set("limit", "20")
	if after != nil {
		params.Set("after", after.Format(time.RFC3339))
	}
	u.RawQuery = params.Encode()

	client := cleanhttp.DefaultClient()
	log.Printf("calling %q", u.String())
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned %q", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []release
	err = json.Unmarshal(b, &releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
