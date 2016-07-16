package main

import (
  "fmt"
  "errors"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

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
  Message string `json:"message"`
  Status string `json:"status"`
}

/*
 * {
 *   "workflow": "quasar",
 *   "guid": "b8ff485a0306",
 *   "stack_name": "howdy-batch--0-38-131--35s2gfqc",
 *   "deployed_at": 1467144100546,
 *   "unit": "howdy-batch",
 *   "type": "job",
 *   "expiration": 1468353700590,
 *   "statuses": [
 *     {
 *       "timestamp": "2016-06-28T20:02:34.449Z",
 *       "message": "instructing perryman's chronos to handle job container",
 *       "status": "deploying"
 *     }
 *   ]
 * }
 */
type StackSummary struct {
  Workflow string `json:"workflow"`
  Guid string `json:"guid"`
  StackName string `json:"stack_name"`
  DeployedAt int64 `json:"deployed_at"`
  UnitName string `json:"unit"`
  Type string `json:"type"`
  NamespaceRef string `json:"namespace"`
  Expiration int64 `json:"deployed_at"`
  Statuses []StackStatus `json:"statuses"`
}

func InspectStack(guid string, http *gorequest.SuperAgent, cfg *Config) (result StackSummary, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/deployments/"+guid), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    errs = append(errs, errors.New("Bad response from Nelson server"))
    // zomg, constructor just nulls *all* the fields, because who cares about correctness!
    return StackSummary { }, errs
  } else {
    var result StackSummary
    if err := json.Unmarshal(bytes, &result); err != nil {
      panic(err)
    }
    return result, errs
  }
}

func PrintInspectStack(s StackSummary){
  //// stack information
  var tabulized = [][]string {}
  tabulized = append(tabulized, []string{ "GUID:", s.Guid })
  tabulized = append(tabulized, []string{ "STACK NAME:", s.StackName })
  tabulized = append(tabulized, []string{ "NAMESPACE:", s.NamespaceRef })
  tabulized = append(tabulized, []string{ "WORKFLOW:", s.Workflow })
  tabulized = append(tabulized, []string{ "TYPE:", s.Type })
  tabulized = append(tabulized, []string{ "DEPLOYED AT:", JavaEpochToDateStr(s.DeployedAt) })
  tabulized = append(tabulized, []string{ "EXPIRES AT:", JavaEpochToDateStr(s.Expiration) })

  fmt.Println("===>> Stack Information")
  RenderTableToStdout([]string{ "Paramater", "Value" }, tabulized)

  //// status history
  fmt.Println("") // give us a new line for spacing
  fmt.Println("===>> Status History")

  var statuslines = [][]string {}
  for _,o := range s.Statuses {
    statuslines = append(statuslines,[]string{ o.Status, o.Timestamp, o.Message })
  }

  RenderTableToStdout([]string{ "Status", "Timestamp", "Message" }, statuslines)
}

/////////////////// DEPRECATION ///////////////////

/*
 * {
 *   "service_type": "heydiddlyho-http",
 *   "version":{
 *     "major":1,
 *     "minor":33
 *   }
 * }
 */
type DeprecationRequest struct {
  ServiceType string `json:"service_type"`
  Version FeatureVersion `json:"version"`
}

func Deprecate(req DeprecationRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error){
  r, body, errs := AugmentRequest(
    http.Post(cfg.Endpoint+"/v1/deployments/deprecate"), cfg).Send(req).EndBytes()

  if (r.StatusCode / 100 != 2){
    resp := string(body[:])
    errs = append(errs, errors.New("Unexpected response from Nelson server"))
    return resp, errs
  } else {
    return "Requested deprecation of "+req.ServiceType+" "+strconv.Itoa(req.Version.Major)+"."+strconv.Itoa(req.Version.Minor), errs
  }
}

/////////////////// REDEPLOYMENT ///////////////////

func Redeploy(guid string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error){
  r, body, errs := AugmentRequest(
    http.Post(cfg.Endpoint+"/v1/deployments/"+guid+"/redeploy"), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    resp := string(body[:])
    errs = append(errs, errors.New("bad response from Nelson server"))
    return resp, errs
  } else {
    return "Redeployment requested.", errs
  }
}

/////////////////// LISTING STACKs ///////////////////

/**
 * {
 *   "workflow": "quasar",
 *   "guid": "67e04d28d6ab",
 *   "stack_name": "blobstore-testsuite--0-1-55--kbqg9nff",
 *   "deployed_at": 1467225866870,
 *   "unit": "blobstore-testsuite",
 *   "type": "job",
 *   "namespace": "devel"
 * }
 */
type Stack struct {
  Workflow string `json:"workflow"`
  Guid string `json:"guid"`
  StackName string `json:"stack_name"`
  DeployedAt int64 `json:"deployed_at"`
  UnitName string `json:"unit"`
  Type string `json:"type"`
  NamespaceRef string `json:"namespace"`
}

func ListStacks(dc string, http *gorequest.SuperAgent, cfg *Config) (list []Stack, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/datacenters/"+dc+"/deployments"), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
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

func PrintListStacks(stacks []Stack){
  var tabulized = [][]string {}
  for _,s := range stacks {
    tabulized = append(tabulized,[]string{ s.Guid, s.NamespaceRef, s.StackName, s.Type, s.Workflow, JavaEpochToDateStr(s.DeployedAt) })
  }

  RenderTableToStdout([]string{ "GUID", "Namespace", "Stack", "Type", "Workflow", "Deployed At" }, tabulized)
}

/////////////////// DEPLOYMENT LOG ///////////////////

type StackLog struct {
  Content []string `json:"content"`
  Offset int `json:"offset"`
}

// v1/deployments/:id/log
func GetDeploymentLog(guid string, http *gorequest.SuperAgent, cfg *Config){
  _, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/deployments/"+guid+"/log"), cfg).EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var logs StackLog
  if err := json.Unmarshal(bytes, &logs); err != nil {
    panic(err)
  }

  fmt.Println("===>> logs for stack "+ guid)

  for _,l := range logs.Content {
    fmt.Println(l)
  }
}