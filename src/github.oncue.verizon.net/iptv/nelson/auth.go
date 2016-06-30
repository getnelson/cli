package main

import (
  "os"
  "log"
  "github.com/parnurzeal/gorequest"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
)

type CreateSessionRequest struct {
  AccessToken string `json:"access_token"`
}
// { "session_token": "xxx", "expires_at": 12345 }
type Session struct {
  SessionToken string `json:"session_token"`
  ExpiresAt int64 `json:"expires_at"`
}

func login(client *gorequest.SuperAgent, githubToken string, nelsonHost string) bool {
  sess := createSession(client, githubToken, nelsonHost)
  createConfigFile(sess)
  return true
}

/* TODO: any error handling here... would be nice */
func createSession(client *gorequest.SuperAgent, githubToken string, nelsonHost string) Session {
  ver := CreateSessionRequest { AccessToken: githubToken }

  _, bytes, errs := client.
      SetDebug(true).
      Post("https://"+nelsonHost+"/auth/github").
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

type YamlConfig struct {
  YamlSession `yaml:"session"`
}

type YamlSession struct {
  Token string `yaml:"token"`
  ExpiresAt int64 `yaml:"expires_at"`
}

// returns Unit, no error handling. YOLO
func createConfigFile(s Session) {
  targetDir := os.Getenv("HOME") + "/.nelson"
  os.Mkdir(targetDir, 0644)
  path := targetDir + "/config.yml"

  temp := &YamlConfig {
      YamlSession: YamlSession {
        Token: s.SessionToken,
        ExpiresAt: s.ExpiresAt,
      },
    }

  d, err := yaml.Marshal(&temp)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  err = ioutil.WriteFile(path, []byte(string(d)), 0644)
  if err != nil {
    panic(err)
  }
}
