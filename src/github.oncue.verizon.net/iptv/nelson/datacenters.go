package main

import (
  "fmt"
  "github.com/parnurzeal/gorequest"
)

type Region struct {
  Name string `json:"name"`
  // Namespaces []string
}

func listRegions(http *gorequest.SuperAgent, baseURL string){
  _, bytes, errs := http.
    SetDebug(true).
    Get(baseURL+"/v1/datacenters").
    EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  fmt.Println(len(bytes))
}