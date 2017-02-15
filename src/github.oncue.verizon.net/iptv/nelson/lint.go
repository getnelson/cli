package main

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
)

/*
 * {
 *   "unit": "howdy-http",
 *   "resources": ["s3"],
 *   "template": "credentials.template"
 * }
 */
type LintTemplateRequest struct {
	Unit      string   `json:"unit"`
	Resources []string `json:"resources"`
	Template  string   `json:"template"`
}

type LintTemplateFailure struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func LintTemplate(req LintTemplateRequest, http *gorequest.SuperAgent, cfg *Config) (str string, err []error) {
	r, body, errs := AugmentRequest(http.Post(cfg.Endpoint+"/v1/validate-template"), cfg).Send(req).EndBytes()

	if r.StatusCode/100 == 2 {
		return "Template rendered successfully.\nRendered output discarded for security reasons.", errs
	} else if r.StatusCode == 400 || r.StatusCode == 504 {
		var fail LintTemplateFailure
		if err := json.Unmarshal(body, &fail); err != nil {
			resp := string(body[:])
			errs = append(errs, errors.New("Unexpected response from Nelson server"))
			return resp, errs
		}
		errs = append(errs, errors.New(fail.Message))
		return fail.Details, errs
	} else {
		resp := string(body[:])
		errs = append(errs, errors.New("Unexpected response from Nelson server"))
		return resp, errs
	}
}
