package main

import (
  "os"
  // "fmt"
  "github.com/parnurzeal/gorequest"
  "io/ioutil"
  "encoding/json"
)

type CreateSessionRequest struct {
  AccessToken string `json:"access_token"`
}
// { "session_token": "xxx", "expires_at": 12345 }
type Session struct {
  SessionToken string `json:"session_token"`
  ExpiresAt int64 `json:"expires_at"`
}

func login(client *gorequest.SuperAgent, githubToken string) bool {
  sess := createSession(client, githubToken)

  createAuthFile(sess)

  return true
}

/* TODO: any error handling here... would be nice */
func createSession(client *gorequest.SuperAgent, githubToken string) Session {
  ver := CreateSessionRequest { AccessToken: githubToken }

  _, bytes, errs := client.
      SetDebug(true).
      Post("https://nelson-beta.oncue.verizon.net/auth/github").
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

// returns Unit, no error handling. YOLO
func createAuthFile(s Session) {
  targetDir := os.Getenv("HOME") + "/.nelson"
  os.Mkdir(targetDir, 0644)
  err := ioutil.WriteFile(targetDir + "/config.yml", []byte(s.SessionToken), 0644)
  if err != nil {
    panic(err)
  }
}
