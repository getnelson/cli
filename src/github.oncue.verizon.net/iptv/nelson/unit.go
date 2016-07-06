package main

import (
  "errors"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

/*
 * {
 *   "nuid": "3fbc7381a664",
 *   "namespace": "devel",
 *   "service_type": "heydiddlyho-http",
 *   "version": {
 *     "major": 0,
 *     "minor": 33
 *   }
 * }
 */
type UnitSummary struct {
  Guid string `json:"nuid"`
  NamespaceRef string `json:"namespace"`
  ServiceType string `json:"service_type"`
  Version FeatureVersion `json:"version"`
}

type FeatureVersion struct {
  Major int `json:"major"`
  Minor int `json:"minor"`
}

/////////////////// LIST ///////////////////

func ListUnits(dc string, http *gorequest.SuperAgent, cfg *Config) (list []UnitSummary, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/datacenters/"+dc+"/units?status=active,manual"), cfg).SetDebug(false).EndBytes()

  if (r.StatusCode / 100 != 2){
    errs = append(errs, errors.New("bad response from Nelson server"))
    return nil, errs
  } else {
    var list []UnitSummary
    if err := json.Unmarshal(bytes, &list); err != nil {
      panic(err)
    }
    return list, errs
  }
}

func PrintListUnits(units []UnitSummary){
  var tabulized = [][]string {}
  for _,u := range units {
    tabulized = append(tabulized,[]string{ u.Guid, u.NamespaceRef, u.ServiceType, strconv.Itoa(u.Version.Major)+"."+strconv.Itoa(u.Version.Minor) })
  }

  RenderTableToStdout([]string{ "GUID", "Namespace", "Unit", "Version" }, tabulized)
}
