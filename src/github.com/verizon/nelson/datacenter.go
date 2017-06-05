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

type Datacenter struct {
	Name       string      `json:"name"`
	Namespaces []Namespace `json:"namespaces"`
}
type Namespace struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func ListDatacenters(http *gorequest.SuperAgent, cfg *Config) (list []Datacenter, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/datacenters"), cfg).EndBytes()

	if errs != nil {
		return nil, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Bad response from Nelson server"))
		return nil, errs
	} else {
		var list []Datacenter
		if err := json.Unmarshal(bytes, &list); err != nil {
			panic(err)
		}
		return list, errs
	}
}

func PrintListDatacenters(datacenters []Datacenter) {
	var tabulized = [][]string{}
	for _, r := range datacenters {
		namespace := ""
		for i, ns := range r.Namespaces {
			if i == 0 {
				namespace = ns.Name
			} else {
				namespace = namespace + ", " + ns.Name
			}
		}
		tabulized = append(tabulized, []string{r.Name, namespace})
	}

	RenderTableToStdout([]string{"Datacenter", "Namespaces"}, tabulized)
}
