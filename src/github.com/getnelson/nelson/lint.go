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
	"strconv"
)

/*
 * {
 *   "unit": "howdy-http",
 *   "resources": ["s3"],
 *   "template": "credentials.template"
 * }
 */
type LintTemplateRequest struct {
	Unit      string   `json:"unit"`
	Resources []string `json:"resources"`
	Template  string   `json:"template"`
}

type LintTemplateFailure struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func LintTemplate(req LintTemplateRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(http.Post(cfg.Endpoint+"/v1/validate-template"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 == 2 {
		return "Template rendered successfully.\nRendered output discarded for security reasons.", errs
	} else if r.StatusCode == 400 || r.StatusCode == 504 {
		var fail LintTemplateFailure
		if err := json.Unmarshal(body, &fail); err != nil {
			resp := string(body[:])
			errs = append(errs, errors.New("Unexpected response from Nelson server: JSON error"))
			return resp, errs
		}
		errs = append(errs, errors.New(fail.Message))
		return fail.Details, errs
	} else {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server: HTTP status "+strconv.Itoa(r.StatusCode)))
		return resp, errs
	}
}

/*
 * {
 *   "units": [{"kind":"howdy-http", "name":"howdy-http@1.2"}],
 *   "manifest": "CAgICAgIHBsYW5zOg0KICAgICAgICAgIC0gZGVmYXVsdA=="
 * }
 */
type LintManifestRequest struct {
	Units    []ManifestUnit `json:"units"`
	Manifest string         `json:"manifest"`
}

type ManifestUnit struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

func LintManifest(req LintManifestRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(http.Post(cfg.Endpoint+"/v1/lint"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 == 2 {
		return "Nelson manifest validated with no errors.\n", errs
	} else if r.StatusCode == 400 || r.StatusCode == 504 {
		resp := string(body[:])
		errs = append(errs, errors.New(""))
		return resp, errs
	} else {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server: HTTP status "+strconv.Itoa(r.StatusCode)))
		return resp, errs
	}
}
