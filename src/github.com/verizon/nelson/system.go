package main

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"strconv"
)

/*
 * {
 *   "description": "retains the latest version",
 *   "policy": "retain-latest"
 * }
 */
type CleanupPolicy struct {
	Description string `json:"description"`
	Policy      string `json:"policy"`
}

func ListCleanupPolicies(http *gorequest.SuperAgent, cfg *Config) (list []CleanupPolicy, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/cleanup-policies"), cfg).EndBytes()

	if r != nil {
		if r.StatusCode/100 != 2 {
			codeAsStr := strconv.Itoa(r.StatusCode)
			errs = append(errs, errors.New("Unxpected response from Nelson server ["+codeAsStr+"]"))
			return nil, errs
		} else {
			var list []CleanupPolicy
			if err := json.Unmarshal(bytes, &list); err != nil {
				panic(err)
			}
			return list, errs
		}
	} else {
		return nil, errs
	}
}

func PrintCleanupPolicies(policies []CleanupPolicy) {
	var tabulized = [][]string{}
	for _, s := range policies {
		tabulized = append(tabulized, []string{s.Policy, s.Description})
	}
	RenderTableToStdout([]string{"Policy", "Description"}, tabulized)
}

type BuildInfoResponse struct {
	BuildInfo BuildInfo `json:"build_info"`
	Banner    string    `json:"banner"`
}

type BuildInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	ScalaVersion string `json:"scala_version"`
	SbtVersion   string `json:"sbt_version"`
	GitRevision  string `json:"git_revision"`
	BuildDate    string `json:"build_date"`
}

// GET /build-info
func WhoAreYou(http *gorequest.SuperAgent, cfg *Config) (resp BuildInfoResponse, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/build-info"), cfg).EndBytes()

	if r.StatusCode != 200 {
		errs = append(errs, errors.New("Bad response from Nelson server"))
		return BuildInfoResponse{}, errs
	} else {
		var resp BuildInfoResponse
		if err := json.Unmarshal(bytes, &resp); err != nil {
			panic(err)
		}
		return resp, errs
	}
}
