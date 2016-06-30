package main

import (
  "os"
  // "fmt"
  "log"
  "time"
  // "github.com/parnurzeal/gorequest"
  "io/ioutil"
  // "encoding/json"
  "gopkg.in/yaml.v2"
  "net/http"
)

///////////////////////////// CLI ENTRYPOINT //////////////////////////////////

func LoadDefaultConfig() *Config {
  return readConfigFile(defaultConfigPath())
}

////////////////////////////// CONFIG YAML ///////////////////////////////////

type Config struct {
  Endpoint string `yaml:endpoint`
  ConfigSession `yaml:"session"`
}

type ConfigSession struct {
  Token string `yaml:"token"`
  ExpiresAt int64 `yaml:"expires_at"`
}

func (c *Config) GetAuthCookie() *http.Cookie {
  expire := time.Now().AddDate(0, 0, 1)
  cookie := &http.Cookie {
    Name: "nelson.session",
    Value: "'"+c.ConfigSession.Token+"'",
    Path: "/",
    Domain: "nelson-beta.oncue.verizon.net",
    Expires: expire,
    RawExpires: expire.Format(time.UnixDate),
    MaxAge: 86400,
    Secure: true,
    HttpOnly: false,
  }

  return cookie

  // return &http.Cookie {
  //   Name: "nelson.session",
  //   Value: c.ConfigSession.Token,
  //   MaxAge: int(c.ConfigSession.ExpiresAt),
  //   Secure: true,
  //   HttpOnly: false,
  // }
}

func generateConfigYaml(s Session, url string) string {
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

  return "---\n"+string(d)
}

func parseConfigYaml(yamlAsBytes []byte) *Config {
  temp := &Config {}
  err := yaml.Unmarshal(yamlAsBytes, &temp)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  return temp
}

/////////////////////////////// CONFIG I/O ////////////////////////////////////

func defaultConfigPath() string {
  targetDir := os.Getenv("HOME")+"/.nelson"
  os.Mkdir(targetDir, 0775)
  return targetDir + "/config.yml"
}

// returns Unit, no error handling. YOLO
func writeConfigFile(s Session, url string, configPath string) {
  yamlConfig := generateConfigYaml(s,url)

  err := ioutil.WriteFile(configPath, []byte(yamlConfig), 0755)
  if err != nil {
    panic(err)
  }
}

func readConfigFile(configPath string) *Config {
  if _, err := os.Stat(configPath); os.IsNotExist(err) {
    panic("No config file existed at "+configPath+". You need to `nelson login` before running other commands.")
  }

  b, err := ioutil.ReadFile(configPath)
  if err != nil {
      panic(err)
  }
  return parseConfigYaml(b)
}
