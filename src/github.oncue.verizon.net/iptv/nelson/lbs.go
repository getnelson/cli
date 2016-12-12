package main

import (
  "errors"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
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
  Name string `json:"name"`
  MajorVersion int `json:"major_version"`
  Datacenter string `json:"datacenter"`
  Namespace string `json:"namespace"`
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
  Name string `json:"name"`
  Routes []LoadbalancerRoute `json:"routes"`
  Guid string `json:"guid"`
  DeployTime int64 `json:"deploy_time"`
  Datacenter string `json:"datacenter"`
  NamespaceRef string `json:"namespace"`
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
  BackendMajorVersion int `json:"backend_major_version"`
  BackendName string `json:"backend_name"`
  LBPort int `json:"lb_port"`
}

//////////////////////// LIST ////////////////////////

func ListLoadbalancers(delimitedDcs string, delimitedNamespaces string, delimitedStatuses string, http *gorequest.SuperAgent, cfg *Config) (list []Loadbalancer, err []error){
  uri := "/v1/loadbalancers?"
  // set the datacenters if specified
  if (isValidCommaDelimitedList(delimitedDcs)){
    uri = uri+"dc="+delimitedDcs+"&"
  }
  if (isValidCommaDelimitedList(delimitedNamespaces)){
    uri = uri+"ns="+delimitedNamespaces
  } else {
    uri = uri+"ns=dev,qa,prod"
  }

  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+uri), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    errs = append(errs, errors.New("bad response from Nelson server"))
    return nil, errs
  } else {
    var list []Loadbalancer
    if err := json.Unmarshal(bytes, &list); err != nil {
      panic(err)
    }
    return list, errs
  }
}

func PrintListLoadbalancers(lb []Loadbalancer){
  var tabulized = [][]string {}
  for _,l := range lb {
    routes := ""
    for i,r := range l.Routes {
      // 8443 ~> howdy-http@1->default
      routes = routes + strconv.Itoa(r.LBPort)+" ~> "+r.BackendName+"@"+strconv.Itoa(r.BackendMajorVersion)+"->"+r.BackendPortReference

      // if not the last element, lets bang on a comma
      if(i == len(l.Routes)){ routes = routes + ", " }
    }
    tabulized = append(tabulized,[]string{ l.Guid, l.Datacenter, l.NamespaceRef, l.Name, routes })
  }

  RenderTableToStdout([]string{ "GUID",  "Datacenter", "Namespace", "Name", "Routes"}, tabulized)
}

//////////////////////// REMOVE ////////////////////////

func RemoveLoadBalancer(guid string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error){
  r, body, errs := AugmentRequest(
    http.Delete(cfg.Endpoint+"/v1/loadbalancers/"+guid), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    resp := string(body[:])
    errs = append(errs, errors.New("Unexpected response from Nelson server"))
    return resp, errs
  } else {
    return "Requested removal of "+guid, errs
  }
}

//////////////////////// CREATE ////////////////////////

func CreateLoadBalancer(req LoadbalancerCreate, http *gorequest.SuperAgent, cfg *Config) (str string, err []error){
  r, body, errs := AugmentRequest(
    http.Post(cfg.Endpoint+"/v1/loadbalancers"), cfg).Send(req).EndBytes()

  if (r.StatusCode / 100 != 2){
    resp := string(body[:])
    errs = append(errs, errors.New("Unexpected response from Nelson server"))
    return resp, errs
  } else {
    return "Loadbalancer has been created.", errs
  }
}
