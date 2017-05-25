package main

import (
	"errors"
	"github.com/parnurzeal/gorequest"
)

type NamespaceRequest struct {
	Namespace string `json:"namespace"`
}

func CreateNamespace(req NamespaceRequest, dc string, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {

	r, body, errs := AugmentRequest(
		http.Post(cfg.Endpoint+"/v1/datacenters/"+dc+"/namespaces"), cfg).Send(req).EndBytes()

	if errs != nil {
		return "", errs
	}

	if r.StatusCode/100 != 2 {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	} else {
		return "namespace(s) has been created.", errs
	}
}
