package nodeman

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/semver"
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
	switch v := n.LTS.(type) {
	case string:
		if v == "false" {
			return false
		}
		return true
	case bool:
		return v
	}
	return false
}

// GetLatestNodeVersion gets the latest even numbered node version
func GetLatestNodeVersion() (string, error) {
	releases, err := getNodeReleases()
	if err != nil {
		return "", err
	}
	latest, _ := semver.NewVersion("8")
	for _, schedule := range *releases {
		version, err := semver.NewVersion(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if version.Major()%2 == 0 && latest.LessThan(version) {
			latest = version
		}
	}
	return latest.String(), nil
}

// GetNodeVersionByRangeOrLTS return the latest matching node version in the range,
// or the latest LTS version if the range is invalid
func GetNodeVersionByRangeOrLTS(engine string) (string, error) {
	versionRange, err := semver.NewConstraint(engine)
	if err != nil {
		return GetLatestNodeVersion()
	}
	releases, err := getNodeReleases()
	if err != nil {
		return "", err
	}
	latest, _ := semver.NewVersion("8")
	for _, schedule := range *releases {
		version, err := semver.NewVersion(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if version.Major()%2 == 0 && latest.LessThan(version) && versionRange.Check(version) {
			latest = version
		}
	}
	return latest.String(), nil
}

// GetLatestLTSNodeVersion gets the latest LTS version of node.js
func GetLatestLTSNodeVersion() (string, error) {
	releases, err := getNodeReleases()
	if err != nil {
		return "", err
	}
	latest, _ := semver.NewVersion("8")
	for _, schedule := range *releases {
		version, err := semver.NewVersion(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if schedule.isLTS() && latest.Compare(version) == -1 {
			latest = version
		}
	}
	return latest.String(), nil
}

func getNodeReleases() (*[]nodeLTSSchedule, error) {
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
		return nil, err
	}
	return &jsonSchedules, nil
}
