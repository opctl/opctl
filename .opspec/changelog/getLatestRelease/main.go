package main

import (
	"encoding/json"
	"golang.org/x/mod/semver"
	"os"

	changelog "github.com/anton-yurchenko/go-changelog"
)

type Release struct {
	Description  string `json:"description"`
	IsPrerelease bool   `json:"isPreRelease"`
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

	lateReleaseJSON, err := json.Marshal(
		Release{
			Description:  latestRelease.Changes.ToString(),
			IsPrerelease: semver.Prerelease(*latestRelease.Version) != "",
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
