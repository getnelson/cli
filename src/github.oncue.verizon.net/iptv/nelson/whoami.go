package main

import (
  "fmt"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

type SessionResponse struct {
  User User `json:"user"`
}

type User struct {
  Login string `json:"login"`
  Name string `json:"name"`
  Avatar string `json:"avatar"`
}

// GET /session
func WhoAmI(http *gorequest.SuperAgent, cfg *Config){
  _, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/session"), cfg).EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var resp SessionResponse
  if err := json.Unmarshal(bytes, &resp); err != nil {
    panic(err)
  }

  fmt.Println("===>> Currently logged in as "+resp.User.Name)
}