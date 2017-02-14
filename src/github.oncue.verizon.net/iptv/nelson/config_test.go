package main

import (
	"os"
	"testing"
)

func TestGenerateConfigYaml(t *testing.T) {
	sess := Session{SessionToken: "abc", ExpiresAt: 1234}
	cfg := generateConfigYaml(sess, "http://foo.com")
	fixture := `---
endpoint: http://foo.com
session:
  token: abc
  expires_at: 1234
`

	if cfg != fixture {
		t.Error("Exected \n"+fixture+"\nbut got:\n", cfg)
	}
}

func TestRoundTripConfigFile(t *testing.T) {
	host := "http://foo.com"
	sess := Session{SessionToken: "abc", ExpiresAt: 1234}
	path := "/tmp/nelson-cli-config-test.yml"
	expected := Config{
		Endpoint: host,
		ConfigSession: ConfigSession{
			Token:     sess.SessionToken,
			ExpiresAt: sess.ExpiresAt,
		},
	}

	writeConfigFile(sess, host, path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Expected a file to exist at ", path)
	}

	_, loadedCfg := readConfigFile(path)

	if expected != *loadedCfg {
		t.Error(expected, loadedCfg)
	}
}
