package main

import (
	"encoding/json"
	"os"

	"github.com/Masterminds/semver/v3"

	changelog "github.com/anton-yurchenko/go-changelog"
)

type Release struct {
	Description  string `json:"description"`
	IsPrerelease bool   `json:"isPrerelease"`
	Version      string `json:"version"`
}

func main() {
	p, err := changelog.NewParser("/CHANGELOG.md")
	if err != nil {
		panic(err)
	}

	c, err := p.Parse()
	if err != nil {
		panic(err)
	}

	latestRelease := c.Releases[0]
	sv, err := semver.NewVersion(*latestRelease.Version)
	if err != nil {
		panic(err)
	}

	lateReleaseJSON, err := json.Marshal(
		Release{
			Description:  latestRelease.Changes.ToString(),
			IsPrerelease: sv.Prerelease() != "",
			Version:      *latestRelease.Version,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("/latestRelease.json", lateReleaseJSON, 0700); err != nil {
		panic(err)
	}
}
