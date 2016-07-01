package main

import (
  "fmt"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

type StackLog struct {
  Content []string `json:"content"`
  Offset int `json:"offset"`
}

// v1/deployments/:id/log
func GetDeploymentLog(id int, http *gorequest.SuperAgent, cfg *Config){
  idAsStr := strconv.Itoa(id)
  _, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/deployments/"+idAsStr+"/log"), cfg).EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var logs StackLog
  if err := json.Unmarshal(bytes, &logs); err != nil {
    panic(err)
  }

  fmt.Println("===>> logs for deployment "+ idAsStr)

  for _,l := range logs.Content {
    fmt.Println(l)
  }
}