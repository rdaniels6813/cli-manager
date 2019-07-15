package nodeman

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/blang/semver"
)

type nodeLTSSchedule struct {
	Version  string   `json:"version"`
	Date     string   `json:"date"`
	Files    []string `json:"files"`
	NPM      string   `json:"npm"`
	LTS      ltsType  `json:"lts"`
	Security bool     `json:"security"`
}

type ltsType interface{}

func (n *nodeLTSSchedule) isLTS() bool {
	switch n.LTS.(type) {
	case string:
		if n.LTS == "false" {
			return false
		}
		return true
	case bool:
		return n.LTS.(bool)
	}
	return false
}

// GetLatestNodeVersion gets the latest even numbered node version
func GetLatestNodeVersion() string {
	releases := getNodeReleases()
	latest := semver.Version{Major: 8}
	for _, schedule := range *releases {
		version, err := semver.Make(strings.ReplaceAll(schedule.Version, "v", ""))
		if err != nil {
			log.Fatal(err)
		}
		if version.Major%2 == 0 && latest.Compare(version) == -1 {
			latest = version
		}
	}
	return latest.String()
}

// GetLatestLTSNodeVersion gets the latest LTS version of node.js
func GetLatestLTSNodeVersion() string {
	releases := getNodeReleases()
	latest := semver.Version{Major: 8}
	for _, schedule := range *releases {
		version, err := semver.Make(strings.ReplaceAll(schedule.Version, "v", ""))
		if err != nil {
			log.Fatal(err)
		}
		if schedule.isLTS() && latest.Compare(version) == -1 {
			latest = version
		}
	}
	return latest.String()
}

func getNodeReleases() *[]nodeLTSSchedule {
	jsonURL := "https://nodejs.org/dist/index.json"
	req, err := http.NewRequest("GET", jsonURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var jsonSchedules []nodeLTSSchedule
	err = json.NewDecoder(resp.Body).Decode(&jsonSchedules)
	if err != nil {
		log.Fatal(err)
	}
	return &jsonSchedules
}
