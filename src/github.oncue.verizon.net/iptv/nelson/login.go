package main

import (
  // "os"
  // "fmt"
  // "log"
  "github.com/parnurzeal/gorequest"
  // "io/ioutil"
  "encoding/json"
  // "gopkg.in/yaml.v2"
)

type CreateSessionRequest struct {
  AccessToken string `json:"access_token"`
}
// { "session_token": "xxx", "expires_at": 12345 }
type Session struct {
  SessionToken string `json:"session_token"`
  ExpiresAt int64 `json:"expires_at"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func Login(client *gorequest.SuperAgent, githubToken string, nelsonHost string, disableTLS bool) bool {
  baseURL := createEndpointURL(nelsonHost, !disableTLS)
  sess := createSession(client, githubToken, baseURL)
  writeConfigFile(sess, baseURL, defaultConfigPath())
  return true
}

///////////////////////////// INTERNALS ////////////////////////////////

func createEndpointURL(host string, useTLS bool) string {
  u := "://"+host
  if(useTLS){
    return "https"+u
  } else {
    return "http"+u
  }
}

/* TODO: any error handling here... would be nice */
func createSession(client *gorequest.SuperAgent, githubToken string, baseURL string) Session {
  ver := CreateSessionRequest { AccessToken: githubToken }
  url := baseURL+"/auth/github"
  _, bytes, errs := client.
      Post(url).
      Send(ver).
      EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var result Session
  if err := json.Unmarshal(bytes, &result); err != nil {
    //return err, nil
    panic(err)
  }

  return result
}