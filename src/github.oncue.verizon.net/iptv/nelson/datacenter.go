package main

import (
  // "time"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

type Datacenter struct {
  Name string `json:"name"`
  Namespaces []Namespace `json:"namespaces"`
}
type Namespace struct {
  Id int `json:"id"`
  Name string `json:"name"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func ListRegions(http *gorequest.SuperAgent, cfg *Config){

  _, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/datacenters"), cfg).EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var datacenters []Datacenter
  if err := json.Unmarshal(bytes, &datacenters); err != nil {
    panic(err)
  }

  var tabulized = [][]string {}
  for _,r := range datacenters {
    namespace := ""
    for i,ns := range r.Namespaces {
      if(i == 0){
        namespace = ns.Name
      } else {
        namespace = namespace+", "+ns.Name
      }
    }
    tabulized = append(tabulized,[]string{ r.Name, namespace })
  }

  RenderTableToStdout([]string{ "Region", "Namespaces" }, tabulized)
}