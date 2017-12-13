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

type SessionResponse struct {
	User User `json:"user"`
}

type User struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GET /session
func WhoAmI(http *gorequest.SuperAgent, cfg *Config) (resp SessionResponse, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/session"), cfg).EndBytes()

	if errs != nil {
		return SessionResponse{}, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Bad response from Nelson server"))
		return SessionResponse{}, errs
	} else {
		var resp SessionResponse
		if err := json.Unmarshal(bytes, &resp); err != nil {
			panic(err)
		}
		return resp, errs
	}
}
