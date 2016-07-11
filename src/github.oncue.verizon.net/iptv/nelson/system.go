package main

import (
  "errors"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

/*
 * {
 *   "description": "retains the latest version",
 *   "policy": "retain-latest"
 * }
 */
type CleanupPolicy struct {
  Description string `json:"description"`
  Policy string `json:"policy"`
}

func ListCleanupPolicies(http *gorequest.SuperAgent, cfg *Config) (list []CleanupPolicy, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/cleanup-policies"), cfg).SetDebug(false).EndBytes()

  if r != nil {
    if (r.StatusCode / 100 != 2){
      codeAsStr := strconv.Itoa(r.StatusCode)
      errs = append(errs, errors.New("Unxpected response from Nelson server ["+codeAsStr+"]"))
      return nil, errs
    } else {
      var list []CleanupPolicy
      if err := json.Unmarshal(bytes, &list); err != nil {
        panic(err)
      }
      return list, errs
    }
  } else {
    return nil, errs
  }
}

func PrintCleanupPolicies(policies []CleanupPolicy){
  var tabulized = [][]string {}
  for _,s := range policies {
    tabulized = append(tabulized,[]string{ s.Policy, s.Description })
  }
  RenderTableToStdout([]string{ "Policy", "Description" }, tabulized)
}