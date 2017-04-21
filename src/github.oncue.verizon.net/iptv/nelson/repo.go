package main

import (
	"errors"
	"github.com/parnurzeal/gorequest"
)

type EnableRepoRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func Enable(req EnableRepoRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	uri := cfg.Endpoint + "/v1/repos/" + req.Owner + "/" + req.Repo + "/hook"
	r, body, errs := AugmentRequest(http.Post(uri), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "The project " + req.Owner + "/" + req.Repo + " has been enabled.", errs
	}
}

func Disable(req EnableRepoRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {

	uri := cfg.Endpoint + "/v1/repos/" + req.Owner + "/" + req.Repo + "/hook"
	r, body, errs := AugmentRequest(http.Delete(uri), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "The project " + req.Owner + "/" + req.Repo + " has been disabled.", errs
	}
}
