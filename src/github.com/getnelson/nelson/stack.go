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
	"net/url"
	"strconv"

	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
)

/////////////////// MANUAL DEPLOYMENT ///////////////////

/*
 * {
 *   "datacenter": "perryman",
 *   "namespace": "stage",
 *   "serviceType": "cassandra",
 *   "version": "1.2.3,
 *   "hash": "abcd1234",
 *   "description": "a cassandra for great good",
 *   "port": 1234
 * }
 */
type ManualDeploymentRequest struct {
	Datacenter  string `json:"datacenter"`
	Namespace   string `json:"namespace"`
	ServiceType string `json:"service_type"`
	Version     string `json:"version"`
	Hash        string `json:"hash"`
	Port        int64  `json:"port"`
	Description string `json:"description"`
}

func RegisterManualDeployment(
	req ManualDeploymentRequest,
	http *gorequest.SuperAgent, cfg *Config) (string, []error) {

	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/deployments"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "Manual stack has been registered.", errs
	}
}

/////////////////// REVERSE ///////////////////
type ReverseTrafficShiftErr struct {
	Message string `json:"message"`
}

func ReverseTrafficShift(
	guid string,
	http *gorequest.SuperAgent, cfg *Config) (string, []error) {

	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/deployments/"+guid+"/trafficshift/reverse"), cfg).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Bad response from Nelson server. Status code: "+strconv.Itoa(r.StatusCode)))

		var er ReverseTrafficShiftErr
		if err := json.Unmarshal(body, &er); err == nil {
			errs = append(errs, errors.New(er.Message))
		}

		return resp, errs
	} else {
		return "Traffic shift reversed.", errs
	}

}

/////////////////// INSPECT ///////////////////

/*
 * {
 *   "timestamp": "2016-06-28T20:02:34.449Z",
 *   "message": "instructing perryman's chronos to handle job container",
 *   "status": "deploying"
 * }
 */
type StackStatus struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}

/*
 *   "dependencies": {
 *     "outbound": [
 *      ...
 *     ],
 *     "inbound": []
 *   },
 */
type StackDependencies struct {
	Outbound []Stack `json:"outbound"`
	Inbound  []Stack `json:"inbound"`
}

/*
 * {
 *   "workflow": "pulsar",
 *   "guid": "e4184c271bb9",
 *   "statuses": [
 *     {
 *       "timestamp": "2016-07-14T22:30:22.358Z",
 *       "message": "inventory-inventory deployed to perryman",
 *       "status": "ready"
 *     },
 *     ...
 *   ],
 *   "stack_name": "inventory-inventory--2-0-11--8gufie2b",
 *   "deployed_at": 1468535384221,
 *   "unit": "inventory-inventory",
 *   "plan": "service",
 *   "expiration": 1469928212871,
 *   "dependencies": {
 *     "outbound": [
 *       {
 *         "workflow": "manual",
 *         "guid": "1a69395e919d",
 *         "stack_name": "dev-iptv-cass-dev--4-8-4--mtq2odqzndc0mg",
 *         "deployed_at": 1468518896093,
 *         "unit": "dev-iptv-cass-dev",
 *         "plan": "service"
 *       }
 *     ],
 *     "inbound": []
 *   },
 *   "namespace": "dev"
 * }
 */
type StackSummary struct {
	Workflow     string            `json:"workflow"`
	Guid         string            `json:"guid"`
	StackName    string            `json:"stack_name"`
	DeployedAt   int64             `json:"deployed_at"`
	UnitName     string            `json:"unit"`
	Plan         string            `json:"plan"`
	NamespaceRef string            `json:"namespace"`
	Expiration   int64             `json:"expiration"`
	Statuses     []StackStatus     `json:"statuses"`
	Dependencies StackDependencies `json:"dependencies"`
	Resources    []string          `json:"resources"`
}

func InspectStack(guid string, http *gorequest.SuperAgent, cfg *Config) (result StackSummary, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/deployments/"+guid), cfg).EndBytes()

	if errs != nil {
		return StackSummary{}, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Bad response from Nelson server"))
		// zomg, constructor just nulls *all* the fields, because who cares about correctness!
		return StackSummary{}, errs
	} else {
		var result StackSummary
		if err := json.Unmarshal(bytes, &result); err != nil {
			panic(err)
		}
		return result, errs
	}
}

func PrintInspectStack(s StackSummary) {
	//>>>>>>>>>>> status history
	var tabulized = [][]string{}
	tabulized = append(tabulized, []string{"GUID:", s.Guid})
	tabulized = append(tabulized, []string{"STACK NAME:", s.StackName})
	tabulized = append(tabulized, []string{"NAMESPACE:", s.NamespaceRef})
	tabulized = append(tabulized, []string{"PLAN:", s.Plan})
	tabulized = append(tabulized, []string{"WORKFLOW:", s.Workflow})
	tabulized = append(tabulized, []string{"DEPLOYED AT:", JavaEpochToDateStr(s.DeployedAt)})
	tabulized = append(tabulized, []string{"EXPIRES AT:", JavaEpochToDateStr(s.Expiration)})

	if len(s.Resources) != 0 {
		resources := ""
		for i, r := range s.Resources {
			if i == 0 {
				resources = r
			} else {
				resources = resources + ", " + r
			}
		}
		tabulized = append(tabulized, []string{"RESOURCES:", resources})
	}

	fmt.Println("===>> Stack Information")
	RenderTableToStdout([]string{"Parameter", "Value"}, tabulized)

	//>>>>>>>>>>> dependency information

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if (len(s.Dependencies.Outbound) + len(s.Dependencies.Inbound)) != 0 {
		fmt.Println("") // give us a new line for spacing
		fmt.Println("===>> Dependencies")
		var dependencies = [][]string{}

		if len(s.Dependencies.Outbound) != 0 {
			for _, w := range s.Dependencies.Outbound {
				dependencies = append(dependencies, []string{w.Guid, w.StackName, w.Type, JavaEpochToDateStr(w.DeployedAt), yellow("OUTBOUND"), strconv.FormatInt(w.Weight, 10)})
			}
		}
		if len(s.Dependencies.Inbound) != 0 {
			for _, w := range s.Dependencies.Inbound {
				dependencies = append(dependencies, []string{w.Guid, w.StackName, w.Type, JavaEpochToDateStr(w.DeployedAt), green("INBOUND"), strconv.FormatInt(w.Weight, 10)})
			}
		}
		RenderTableToStdout([]string{"GUID", "Stack", "Type", "Deployed At", "Direction", "Weight"}, dependencies)
	}

	//>>>>>>>>>>> status history
	fmt.Println("") // give us a new line for spacing
	fmt.Println("===>> Status History")
	var statuslines = [][]string{}
	for _, o := range s.Statuses {
		statuslines = append(statuslines, []string{o.Status, o.Timestamp, o.Message})
	}
	RenderTableToStdout([]string{"Status", "Timestamp", "Message"}, statuslines)
}

/////////////////// REDEPLOYMENT ///////////////////

func Redeploy(guid string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/deployments/"+guid+"/redeploy"), cfg).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("bad response from Nelson server"))
		return resp, errs
	} else {
		return "Redeployment requested.", errs
	}
}

/////////////////// LISTING STACKS ///////////////////

/**
 * {
 *   "workflow": "quasar",
 *   "guid": "67e04d28d6ab",
 *   "stack_name": "blobstore-testsuite--0-1-55--kbqg9nff",
 *   "deployed_at": 1467225866870,
 *   "unit": "blobstore-testsuite",
 *   "plan": "fooo",
 *   "namespace": "dev"
 * }
 */
type Stack struct {
	Workflow     string `json:"workflow"`
	Guid         string `json:"guid"`
	StackName    string `json:"stack_name"`
	DeployedAt   int64  `json:"deployed_at"`
	UnitName     string `json:"unit"`
	Plan         string `json:"plan"`
	Type         string `json:"type,omitempty"`
	NamespaceRef string `json:"namespace,omitempty"`
	Status       string `json:"status"`
	Weight       int64  `json:"weight,omitempty"`
}

func ListStacks(delimitedDcs string, delimitedNamespaces string, delimitedStatuses string, unit string, http *gorequest.SuperAgent, cfg *Config) (list []Stack, err []error) {
	uri := "/v1/deployments?"
	qs := url.Values{}
	// set the datacenters if specified
	if isValidCommaDelimitedList(delimitedDcs) {
		qs.Set("dc", delimitedDcs)
	}
	if isValidCommaDelimitedList(delimitedStatuses) {
		qs.Set("status", delimitedStatuses)
	} else {
		// if the user didnt specify statuses, they probally want all the stacks except historical terminated ones.
		qs.Set("status", "pending,deploying,warming,ready,deprecated,failed")
	}
	if isValidCommaDelimitedList(delimitedNamespaces) {
		qs.Set("ns", delimitedNamespaces)
	} else {
		qs.Set("ns", "dev,qa,prod")
	}
	if unit != "" {
		qs.Set("unit", unit)
	}
	uri = uri + qs.Encode()

	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+uri), cfg).EndBytes()

	if errs != nil {
		return nil, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("Bad response from Nelson server"))
		return nil, errs
	} else {
		var list []Stack
		if err := json.Unmarshal(bytes, &list); err != nil {
			panic(err)
		}
		return list, errs
	}
}

func PrintListStacks(stacks []Stack) {
	var tabulized = [][]string{}
	for _, s := range stacks {
		tabulized = append(tabulized, []string{s.Guid, s.NamespaceRef, truncateString(s.StackName, 55), s.Status, s.Plan, s.Workflow, javaEpochToHumanizedTime(s.DeployedAt)})
	}

	RenderTableToStdout([]string{"GUID", "Namespace", "Stack", "Status", "Plan", "Workflow", "Deployed At"}, tabulized)
}

/////////////////// DEPLOYMENT LOG ///////////////////

type StackLog struct {
	Content []string `json:"content"`
	Offset  int      `json:"offset"`
}

// v1/deployments/:id/log
func GetDeploymentLog(guid string, http *gorequest.SuperAgent, cfg *Config) {
	_, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/deployments/"+guid+"/log"), cfg).EndBytes()

	if len(errs) > 0 {
		panic(errs)
	}

	var logs StackLog
	if err := json.Unmarshal(bytes, &logs); err != nil {
		panic(err)
	}

	fmt.Println("===>> logs for stack " + guid)

	for _, l := range logs.Content {
		fmt.Println(l)
	}
}

/////////////////// RUNTIME INSPECT ///////////////////

/*
 * {
 *   "consul_health": [{
 *     "check_id": "64c0a2b5bd2972c9521fb7313b00db7dad58c04c",
 *     "node": "ip-10-113-128-190",
 *     "status": "passing",
 *     "name": "service: default howdy-http--1-0-344--9uuu1mp2 check"
 *   }],
 *   "scheduler": {
 *     "failed": 0,
 *     "completed": 0,
 *     "pending": 0,
 *     "running": 1
 *   },
 *   "current_status": "ready",
 *   "expires_at" : 10101333
 * }
 */
type StackRuntime struct {
	CurrentStatus string                `json:"current_status"`
	ExpiresAt     int64                 `json:"expires_at"`
	Scheduler     StackRuntimeScheduler `json:"scheduler"`
	ConsulHealth  []StackRuntimeHealth  `json:"consul_health"`
}

type StackRuntimeHealth struct {
	CheckId string `json:"check_id"`
	Node    string `json:"node"`
	Status  string `json:"status"`
	Name    string `json:"name"`
}

type StackRuntimeScheduler struct {
	Failed    int `json:"failed"`
	Completed int `json:"completed"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
}

func GetStackRuntime(guid string, http *gorequest.SuperAgent, cfg *Config) (runtime StackRuntime, err []error) {
	r, bytes, errs := AugmentRequest(
		http.Get(cfg.Endpoint+"/v1/deployments/"+guid+"/runtime"), cfg).EndBytes()

	if errs != nil {
		return StackRuntime{}, errs
	}

	if r.StatusCode/100 != 2 {
		errs = append(errs, errors.New("bad response from Nelson server"))
		return runtime, errs
	} else {
		if err := json.Unmarshal(bytes, &runtime); err != nil {
			panic(err)
		}
		return runtime, errs
	}
}

func PrintStackRuntime(r StackRuntime) {

	fmt.Println("")
	fmt.Println("==>> Stack Status")

	var tabulized = [][]string{}
	tabulized = append(tabulized, []string{"STATUS:", r.CurrentStatus})
	tabulized = append(tabulized, []string{"EXPIRES:", JavaEpochToDateStr(r.ExpiresAt)})
	RenderTableToStdout([]string{"Parameter", "Value"}, tabulized)
	fmt.Println("")

	// >>>>>>>>> Scheduler Summary
	tabulized = [][]string{}
	tabulized = append(tabulized, []string{"PENDING:", strconv.Itoa(r.Scheduler.Pending)})
	tabulized = append(tabulized, []string{"RUNNING:", strconv.Itoa(r.Scheduler.Running)})
	tabulized = append(tabulized, []string{"COMPLETED:", strconv.Itoa(r.Scheduler.Completed)})
	tabulized = append(tabulized, []string{"FAILED:", strconv.Itoa(r.Scheduler.Failed)})

	fmt.Println("==>> Scheduler")
	RenderTableToStdout([]string{"Parameter", "Value"}, tabulized)

	// >>>>>>>>> Consul Health
	tabulized = [][]string{}
	for _, s := range r.ConsulHealth {
		tabulized = append(tabulized, []string{s.CheckId, s.Node, s.Status, s.Name})
	}

	fmt.Println("")
	fmt.Println("==>> Health Checks")
	RenderTableToStdout([]string{"ID", "Node", "Status", "Details"}, tabulized)

}
