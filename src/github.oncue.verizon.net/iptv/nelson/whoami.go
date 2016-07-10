package main

import (
  "errors"
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
func WhoAmI(http *gorequest.SuperAgent, cfg *Config) (resp SessionResponse, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/session"), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    errs = append(errs, errors.New("Bad response from Nelson server"))
    return SessionResponse{}, errs
  } else {
    var resp SessionResponse
    if err := json.Unmarshal(bytes, &resp); err != nil {
      panic(err)
    }
    return resp, errs
  }
}