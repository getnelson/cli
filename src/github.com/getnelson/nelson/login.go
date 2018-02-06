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
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"strings"
)

type CreateSessionRequest struct {
	AccessToken string `json:"access_token"`
}

// { "session_token": "xxx", "expires_at": 12345 }
type Session struct {
	SessionToken string `json:"session_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func Login(client *gorequest.SuperAgent, githubToken string, nelsonHost string, disableTLS bool) []error {
	baseURL := createEndpointURL(nelsonHost, !disableTLS)
	e, sess := createSession(client, githubToken, baseURL)
	if e != nil {
		return []error{e}
	}
	writeConfigFile(sess, baseURL, defaultConfigPath()) // TIM: side-effect, discarding errors seems wrong
	return nil
}

///////////////////////////// INTERNALS ////////////////////////////////

func createEndpointURL(host string, useTLS bool) string {
	u := "://" + host
	if useTLS {
		return "https" + u
	} else {
		return "http" + u
	}
}

/* TODO: any error handling here... would be nice */
func createSession(client *gorequest.SuperAgent, githubToken string, baseURL string) (error, Session) {
	ver := CreateSessionRequest{AccessToken: githubToken}
	url := baseURL + "/auth/github"
	_, bytes, errs := client.
		Post(url).
		Set("User-Agent", UserAgentString(globalBuildVersion)).
		Send(ver).
		SetCurlCommand(globalEnableCurl).
		SetDebug(globalEnableDebug).
		EndBytes()

	if len(errs) > 0 {
		errStrs := make([]string, len(errs))
		for i, e := range errs {
			errStrs[i] = e.Error()
		}
		panic(strings.Join(errStrs, "\n"))
	}

	var result Session
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err, Session{}
	}

	return nil, result
}
