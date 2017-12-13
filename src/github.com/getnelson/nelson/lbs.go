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
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"time"
)

/*
 * {
 *   "name": "howdy-lb",
 *   "major_version": 1,
 *   "datacenter": "us-east-1",
 *   "namespace": "dev"
 * }
 */
type LoadbalancerCreate struct {
	Name         string `json:"name"`
	MajorVersion int    `json:"major_version"`
	Datacenter   string `json:"datacenter"`
	Namespace    string `json:"namespace"`
}

/*
 * {
 *   "name": "howdy-lb--1--974u8r6v",
 *   "routes": [
 *     ...
 *   ],
 *   "guid": "b74b8209468b",
 *   "deploy_time": 1481065235649,
 *   "datacenter": "us-east-1",
 *   "namespace": "dev"
 * }
 */
type Loadbalancer struct {
	Name         string              `json:"name"`
	Routes       []LoadbalancerRoute `json:"routes"`
	Guid         string              `json:"guid"`
	DeployTime   int                 `json:"deploy_time"`
	Datacenter   string              `json:"datacenter"`
	Namespace    string              `json:"namespace"`
	Address      string              `json:"address"`
	Version      int                 `json:"major_version"`
	Dependencies DependencyArray     `json:"dependencies"`
}

/*
 * {
 *   "backend_port_reference": "default",
 *   "backend_major_version": 1,
 *   "backend_name": "howdy-http",
 *   "lb_port": 8444
 * }
 */
type LoadbalancerRoute struct {
	BackendPortReference string `json:"backend_port_reference"`
	BackendName          string `json:"backend_name"`
	LBPort               int    `json:"lb_port"`
}

type DependencyArray struct {
	Outbound []LoadbalancerDependencyOutbound `json:"outbound"`
}

type LoadbalancerDependencyOutbound struct {
	DeployTime int    `json:"deployed_at"`
	Type       string `json:"type"`
	StackName  string `json:"stack_name"`
	Guid       string `json:"guid"`
}

//////////////////////// LIST ////////////////////////

func ListLoadbalancers(delimitedDcs string, delimitedNamespaces string, delimitedStatuses string, http *gorequest.SuperAgent, cfg *Config) (list []Loadbalancer, err []error) {
	uri := "/v1/loadbalancers?"
	// set the datacenters if specified
	if isValidCommaDelimitedList(delimitedDcs) {
		uri = uri + "dc=" + delimitedDcs + "&"
	}
	if isValidCommaDelimitedList(delimitedNamespaces) {
		uri = uri + "ns=" + delimitedNamespaces
	} else {
		uri = uri + "ns=dev,qa,prod"
	}

	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+uri), cfg).EndBytes()

	if r != nil {
		if r.StatusCode/100 != 2 {
			errs = append(errs, errors.New("bad response from Nelson server"))
			return nil, errs
		} else {
			var list []Loadbalancer
			if err := json.Unmarshal(bytes, &list); err != nil {
				panic(err)
			}
			return list, errs
		}
	} else {
		errs = append(errs, errors.New("No response from the Nelson server, aborting."))
		return nil, errs
	}
}

func PrintListLoadbalancers(lb []Loadbalancer) {
	var tabulized = [][]string{}
	for _, l := range lb {
		routes := ""
		for i, r := range l.Routes {
			// 8443 ~> howdy-http->default
			routes = routes + strconv.Itoa(r.LBPort) + " ~> " + r.BackendName + "->" + r.BackendPortReference

			// if not the last element, lets bang on a comma
			if i == len(l.Routes) {
				routes = routes + ", "
			}
		}
		tabulized = append(tabulized, []string{l.Guid, l.Datacenter, l.Namespace, l.Name, routes, l.Address})
	}

	RenderTableToStdout([]string{"GUID", "Datacenter", "Namespace", "Name", "Routes", "Address"}, tabulized)
}

func InspectLoadBalancer(guid string, http *gorequest.SuperAgent, cfg *Config) (lb Loadbalancer, err []error) {
	uri := "/v1/loadbalancers/" + guid
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+uri), cfg).EndBytes()

	if errs != nil {
		return Loadbalancer{}, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("bad response from Nelson server"))
		_ = bytes
		return Loadbalancer{}, errs
	} else {
		var lb Loadbalancer
		if err := json.Unmarshal(bytes, &lb); err != nil {
			panic(err)
		}
		return lb, errs
	}
	return Loadbalancer{}, errs
}

func PrintInspectLoadbalancer(lb Loadbalancer) {
	var tabulized = [][]string{}
	tabulized = append(tabulized, []string{"GUID:", lb.Guid})
	tabulized = append(tabulized, []string{"NAME:", lb.Name})
	tabulized = append(tabulized, []string{"ADDRESS:", lb.Address})
	tabulized = append(tabulized, []string{"NAMESPACE:", lb.Namespace})
	tabulized = append(tabulized, []string{"DATECENTER:", lb.Datacenter})
	tabulized = append(tabulized, []string{"VERSION:", strconv.FormatInt(int64(lb.Version), 10)})
	tabulized = append(tabulized, []string{"TIMESTAMP:", time.Unix(int64(lb.DeployTime)/1000, 0).Format(time.RFC3339)})
	fmt.Println("===>> Loadbalancer Information")
	RenderTableToStdout([]string{"Parameter", "Value"}, tabulized)
	if len(lb.Routes) != 0 {
		fmt.Println("")
		fmt.Println("===>> Routes")
		var routes = [][]string{}

		var w LoadbalancerRoute
		for _, w = range lb.Routes {
			routes = append(routes, []string{w.BackendPortReference, w.BackendName, strconv.FormatInt(int64(w.LBPort), 10)})
		}
		RenderTableToStdout([]string{"Reference", "Unit", "Port"}, routes)
	}

	if len(lb.Dependencies.Outbound) != 0 {
		fmt.Println("")
		fmt.Println("===>> Outbound Dependencies")
		var dependencies = [][]string{}

		var r LoadbalancerDependencyOutbound
		for _, r = range lb.Dependencies.Outbound {
			dependencies = append(dependencies, []string{time.Unix(int64(r.DeployTime)/1000, 0).Format(time.RFC3339), r.Type, r.StackName, r.Guid})
		}
		RenderTableToStdout([]string{"Timestamp", "Type", "Stack-Name", "GUID"}, dependencies)
	}
}

//////////////////////// REMOVE ////////////////////////

func RemoveLoadBalancer(guid string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(
		http.Delete(cfg.Endpoint+"/v1/loadbalancers/"+guid), cfg).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "Requested removal of " + guid, errs
	}
}

//////////////////////// CREATE ////////////////////////

func CreateLoadBalancer(req LoadbalancerCreate, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/loadbalancers"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "Loadbalancer has been created.", errs
	}
}
