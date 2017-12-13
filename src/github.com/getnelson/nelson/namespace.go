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
	"errors"
	"github.com/parnurzeal/gorequest"
)

type NamespaceRequest struct {
	Namespace string `json:"namespace"`
}

func CreateNamespace(req NamespaceRequest, dc string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {

	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/datacenters/"+dc+"/namespaces"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "namespace(s) has been created.", errs
	}
}
