//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//:
//:   Licensed under the Apache License, Version 2.0 (the "License");
//:   you may not use this file except in compliance with the License.
//:   You may obtain a copy of the License at
//:
//:       http://www.apache.org/licenses/LICENSE-2.0
//:
//:   Unless required by applicable law or agreed to in writing, software
//:   distributed under the License is distributed on an "AS IS" BASIS,
//:   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//:   See the License for the specific language governing permissions and
//:   limitations under the License.
//:
//: ----------------------------------------------------------------------------
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

func TestConfigValidate(t *testing.T) {
	c := Config{
		Endpoint: "foo.bar.com",
		ConfigSession: ConfigSession{
			Token:     "xxx",
			ExpiresAt: (currentTimeMillis() - 60000), // now minus 1 minute
		},
	}

	if len(c.Validate()) != 1 {
		t.Error(1, len(c.Validate()))
	}
}
