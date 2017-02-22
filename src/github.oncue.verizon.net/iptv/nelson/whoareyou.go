package main

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
)

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
