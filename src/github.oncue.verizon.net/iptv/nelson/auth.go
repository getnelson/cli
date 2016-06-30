package main

import (
  "os"
  "fmt"
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

func login(client *gorequest.SuperAgent, githubToken string, nelsonHost string, disableTLS bool) bool {
  baseURL := createEndpointURL(nelsonHost, !disableTLS)

  fmt.Println(">>> ", baseURL)

  sess := createSession(client, githubToken, baseURL)
  createConfigFile(sess, baseURL)
  return true
}

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
  fmt.Println("~~~~~~~~~~~~ ", url)
  _, bytes, errs := client.
      SetDebug(true).
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

type Config struct {
  Endpoint string `yaml:endpoint`
  ConfigSession `yaml:"session"`
}

type ConfigSession struct {
  Token string `yaml:"token"`
  ExpiresAt int64 `yaml:"expires_at"`
}

// returns Unit, no error handling. YOLO
func createConfigFile(s Session, url string) {
  targetDir := os.Getenv("HOME") + "/.nelson"
  os.Mkdir(targetDir, 0644)
  path := targetDir + "/config.yml"

  temp := &Config {
      Endpoint: url,
      ConfigSession: ConfigSession {
        Token: s.SessionToken,
        ExpiresAt: s.ExpiresAt,
      },
    }

  d, err := yaml.Marshal(&temp)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  yamlConfig := "---\n"+string(d)

  err = ioutil.WriteFile(path, []byte(yamlConfig), 0644)
  if err != nil {
    panic(err)
  }
}

// func loadConfigFile(path string) Config {

// }