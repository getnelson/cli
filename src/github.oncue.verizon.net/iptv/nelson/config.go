package main

import (
  "os"
  // "fmt"
  "log"
  // "github.com/parnurzeal/gorequest"
  "io/ioutil"
  // "encoding/json"
  "gopkg.in/yaml.v2"
)

/////////////////////////////// CONFIG YAML ////////////////////////////////////

type Config struct {
  Endpoint string `yaml:endpoint`
  ConfigSession `yaml:"session"`
}

type ConfigSession struct {
  Token string `yaml:"token"`
  ExpiresAt int64 `yaml:"expires_at"`
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
  os.Mkdir(targetDir, 0644)
  return targetDir + "/config.yml"
}

// returns Unit, no error handling. YOLO
func writeConfigFile(s Session, url string, configPath string) {
  yamlConfig := generateConfigYaml(s,url)

  err := ioutil.WriteFile(configPath, []byte(yamlConfig), 0644)
  if err != nil {
    panic(err)
  }
}

func readConfigFile(configPath string) *Config {
  b, err := ioutil.ReadFile(configPath)
  if err != nil {
      panic(err)
  }
  return parseConfigYaml(b)
}

// func loadConfigFile(path string) { //Config
//   b, err := ioutil.ReadFile("input.txt")
//   if err != nil {
//     panic(err)
//   }
// }