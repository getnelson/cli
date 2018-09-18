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
