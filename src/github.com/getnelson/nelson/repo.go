//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//:
//:   Licensed under the Apache License, Version 2.0 (the "License");
//:   you may not use this file except in compliance with the License.
//:   You may obtain a copy of the License at
//:
//:       http://www.apache.org/licenses/LICENSE-2.0
//:
//:   Unless required by applicable law or agreed to in writing, software
//:   distributed under the License is distributed on an "AS IS" BASIS,
//:   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//:   See the License for the specific language governing permissions and
//:   limitations under the License.
//:
//: ----------------------------------------------------------------------------
package main

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
)

/**
 * {
 * 	 "repository": "xs4s",
 * 	 "slug": "iptv/xs4s",
 * 	 "id": 8272,
 * 	 "hook": {
 *     "is_active": true,
 *     "id": 3775
 *   },
 * 	 "owner": "iptv",
 * 	 "access": "push"
 * }
 */

type RepoHook struct {
	IsActive bool `json:"is_active"`
	Id       int  `json:"id"`
}

type RepoSummary struct {
	Repository string    `json:"repository"`
	Slug       string    `json:"slug"`
	Id         int       `json:"id"`
	Hook       *RepoHook `json:"hook"`
	Owner      string    `json:"owner"`
	Access     string    `json:"access"`
}

func SyncRepos(http *gorequest.SuperAgent, cfg *Config) []error {
	uri := cfg.Endpoint + "/v1/profile/sync"
	r, _, errs := AugmentRequest(http.Post(uri), cfg).EndBytes()

	if errs != nil {
		return errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return errs
	}
	return nil
}

func ListRepos(owner string, http *gorequest.SuperAgent, cfg *Config) (list []RepoSummary, err []error) {
	uri := cfg.Endpoint + "/v1/repos?owner=" + owner
	r, body, errs := AugmentRequest(http.Get(uri), cfg).EndBytes()

	if errs != nil {
		return nil, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return nil, errs
	} else {
		var list []RepoSummary
		if err := json.Unmarshal(body, &list); err != nil {
			panic(err)
		}
		return list, errs
	}
}

type EnableRepoRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func Enable(req EnableRepoRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	uri := cfg.Endpoint + "/v1/repos/" + req.Owner + "/" + req.Repo + "/hook"
	r, body, errs := AugmentRequest(http.Post(uri), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "The project " + req.Owner + "/" + req.Repo + " has been enabled.", errs
	}
}

func Disable(req EnableRepoRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {

	uri := cfg.Endpoint + "/v1/repos/" + req.Owner + "/" + req.Repo + "/hook"
	r, body, errs := AugmentRequest(http.Delete(uri), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "The project " + req.Owner + "/" + req.Repo + " has been disabled.", errs
	}
}

func PrintListRepos(repos []RepoSummary) {
	var tabulized = [][]string{}
	for _, x := range repos {
		var enabled bool = x.Hook != nil && x.Hook.IsActive
		tabulized = append(tabulized, []string{
			x.Repository,
			x.Owner,
			x.Access,
			formatEnabled(enabled),
		})
	}

	RenderTableToStdout([]string{
		"Repository",
		"Owner",
		"Access",
		"Status",
	}, tabulized)
}

func formatEnabled(enabled bool) string {
	if enabled {
		return "enabled"
	} else {
		return "disabled"
	}
}
