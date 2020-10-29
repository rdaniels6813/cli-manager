package nodeman

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/blang/semver/v4"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
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
func GetLatestNodeVersion(client HTTPClient) (string, error) {
	releases, err := getNodeReleases(client)
	if err != nil {
		return "", err
	}
	latest, _ := parseSemver("8")
	for _, schedule := range *releases {
		version, err := parseSemver(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if version.Major%2 == 0 && latest.LT(version) {
			latest = version
		}
	}
	return latest.String(), nil
}

// GetNodeVersionByRangeOrLTS return the latest matching node version in the range,
// or the latest LTS version if the range is invalid
func GetNodeVersionByRangeOrLTS(engine string, client HTTPClient) (string, error) {
	versionRange, err := semver.ParseRange(engine)
	if err != nil {
		fmt.Printf("Error parsing engines range: %s\nUsing latest LTS\n", engine)
		return GetLatestNodeVersion(client)
	}
	releases, err := getNodeReleases(client)
	if err != nil {
		return "", err
	}
	latest, _ := parseSemver("8")
	for _, schedule := range *releases {
		version, err := parseSemver(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if version.Major%2 == 0 && latest.LT(version) && versionRange(version) {
			latest = version
		}
	}
	return latest.String(), nil
}

// GetLatestLTSNodeVersion gets the latest LTS version of node.js
func GetLatestLTSNodeVersion(client HTTPClient) (string, error) {
	releases, err := getNodeReleases(client)
	if err != nil {
		return "", err
	}
	latest, _ := parseSemver("8")
	for _, schedule := range *releases {
		version, err := parseSemver(schedule.Version)
		if err != nil {
			log.Fatal(err)
		}
		if schedule.isLTS() && latest.Compare(version) == -1 {
			latest = version
		}
	}
	return latest.String(), nil
}

func getNodeReleases(client HTTPClient) (*[]nodeLTSSchedule, error) {
	jsonURL := "https://nodejs.org/dist/index.json"
	req, err := http.NewRequest("GET", jsonURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
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

func parseSemver(v string) (semver.Version, error) {
	return semver.Parse(strings.TrimPrefix(v, "v"))
}
