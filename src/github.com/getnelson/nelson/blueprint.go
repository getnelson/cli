//: ----------------------------------------------------------------------------
//: Copyright (C) 2018 Verizon.  All Rights Reserved.
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
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"strconv"
)

/////////////////// PROOFING BLUEPRINTS ///////////////////

/*
 * {
 *   "content": "CAgICAgIHBsYW5zOg0KICAgICAgICAgIC0gZGVmYXVsdA=="
 * }
 */
type ProofBlueprintWire struct {
	Content string `json:"content"`
}

func ProofBlueprint(req ProofBlueprintWire, http *gorequest.SuperAgent, cfg *Config) (res string, err []error) {
	r, body, errs := AugmentRequest(http.Post(cfg.Endpoint+"/v1/blueprints/proof"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 == 2 {
		var result ProofBlueprintWire
		if err := json.Unmarshal(body, &result); err != nil {
			errs = append(errs, err)
		}

		data, err := base64.StdEncoding.DecodeString(result.Content)
		if err != nil {
			errs = append(errs, err)
		}
		// bytes, err := fmt.Printf("%q\n", data)
		debased := string(data[:])

		return debased, errs
	} else {
		errs = append(errs, errors.New("Unexpected response from Nelson server: HTTP status "+strconv.Itoa(r.StatusCode)))
		return "", errs
	}
}

/////////////////// CREATING BLUEPRINTS ///////////////////

/*
 * {
 *    "name": "use-nvidia-1080ti",
 *    "description": "only scheudle on nodes with nvida 1080ti hardware"
 *    "sha256": "1e34a423ebe1fafeda8277386ede3263b01357e490b124b69bc0bfb493e64140"
 *    "template": "<base64 encoded template>"
 * }
 */
type BlueprintResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Revision    string `json:"revision"`
	Sha256      string `json:"sha256"`
	Template    string `json:"template"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateBlueprintRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Sha256      string `json:"sha256"`
	Template    string `json:"template"`
}

func CreateBlueprint(req CreateBlueprintRequest, http *gorequest.SuperAgent, cfg *Config) (out BlueprintResponse, err []error) {

	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/blueprints"), cfg).Send(req).EndBytes()

	var result BlueprintResponse

	if errs != nil {
		return result, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Unexpectedly recieved a "+strconv.Itoa(r.StatusCode)+" reponse from the server."))
		return result, errs
	} else {
		if err := json.Unmarshal(body, &result); err != nil {
			errs = append(errs, errors.New("Unexpected response from Nelson server"))
		}
		return result, errs
	}
}

/////////////////// LISTING BLUEPRINTS ///////////////////

func ListBlueprints(http *gorequest.SuperAgent, cfg *Config) (list []BlueprintResponse, err []error) {
	uri := cfg.Endpoint + "/v1/blueprints"
	r, body, errs := AugmentRequest(http.Get(uri), cfg).EndBytes()

	if errs != nil {
		return nil, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return nil, errs
	} else {
		var list []BlueprintResponse
		if err := json.Unmarshal(body, &list); err != nil {
			panic(err)
		}
		return list, errs
	}
}

func PrintListBlueprints(bps []BlueprintResponse) {
	var tabulized = [][]string{}
	for _, r := range bps {
		name := r.Name + "@" + r.Revision
		tabulized = append(tabulized, []string{name, r.Description, r.Sha256, javaEpochToHumanizedTime(r.CreatedAt)})
	}

	RenderTableToStdout([]string{"Reference", "Description", "Sha256", "Created"}, tabulized)
}

/////////////////// INSPECTING BLUEPRINTS ///////////////////

func InspectBlueprint(namedRevision string, http *gorequest.SuperAgent, cfg *Config) (out BlueprintResponse, err []error) {
	r, body, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/blueprints/"+namedRevision), cfg).EndBytes()

	var result BlueprintResponse

	if errs != nil {
		return result, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Unexpectedly recieved a "+strconv.Itoa(r.StatusCode)+" reponse from the server."))
		return result, errs
	} else {
		if err := json.Unmarshal(body, &result); err != nil {
			errs = append(errs, errors.New("Unexpected response from Nelson server"))
		}
		return result, errs
	}
}
