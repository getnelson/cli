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
	"encoding/base64"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var globalEnableDebug bool
var globalEnableCurl bool
var globalBuildVersion string

func main() {
	year, _, _ := time.Now().Date()
	app := cli.NewApp()
	app.Name = "nelson-cli"
	app.Version = CurrentVersion()
	app.Copyright = "© " + strconv.Itoa(year) + " Verizon Labs"
	app.Usage = "remote control for the Nelson deployment system"
	app.EnableBashCompletion = true

	http := gorequest.New()
	pi := ProgressIndicator()

	// switches for the cli
	var userGithubToken string
	var disableTLS bool
	var selectedDatacenter string
	var selectedNamespace string
	var selectedLoadbalancer string
	var selectedStatus string
	var selectedUnitPrefix string
	var selectedVersion string
	var selectedPort int64
	var selectedServiceType string
	var selectedManifest string
	var selectedTemplate string
	var stackHash string
	var description string
	var selectedNoGrace bool
	var repository string
	var owner string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug mode on the network requests",
			Destination: &globalEnableDebug,
		},
		cli.BoolFlag{
			Name:        "debug-curl",
			Usage:       "Print the curl command analog for the current request",
			Destination: &globalEnableCurl,
		},
	}

	app.Commands = []cli.Command{
		////////////////////////////// LOGIN //////////////////////////////////
		{
			Name:  "login",
			Usage: "login to nelson",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "token, t",
					Value:       "",
					Usage:       "your github personal access token",
					EnvVar:      "GITHUB_TOKEN",
					Destination: &userGithubToken,
				},
				cli.BoolFlag{
					Name:        "disable-tls",
					Destination: &disableTLS,
				},
			},
			Action: func(c *cli.Context) error {
				host := strings.TrimSpace(c.Args().First())
				if len(host) <= 0 {
					host = os.Getenv("NELSON_ADDR")
					if len(host) <= 0 {
						return cli.NewExitError("Either supply a host explicitly, or set $NELSON_ADDR", 1)
					}
				}

				if len(userGithubToken) <= 0 {
					return cli.NewExitError("You must set your GITHUB_TOKEN environment variable or specify a token using -t", 1)
				}

				// fmt.Println("token: ", userGithubToken)
				// fmt.Println("host: ", host)
				pi.Start()
				Login(http, userGithubToken, host, disableTLS)
				pi.Stop()
				fmt.Println("Successfully logged in to " + host)
				return nil
			},
		},
		////////////////////////////// DATACENTER //////////////////////////////////
		{
			Name:    "datacenters",
			Aliases: []string{"dcs", "datacenter"},
			Usage:   "Set of commands for working with Nelson datacenters",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "List all the available datacenters",
					Action: func(c *cli.Context) error {
						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						r, e := ListDatacenters(http, cfg)
						pi.Stop()
						if e != nil {
							return cli.NewExitError("Unable to list datacenters.", 1)
						} else {
							PrintListDatacenters(r)
						}
						return nil
					},
				},
			},
		},
		////////////////////////////// REPOS //////////////////////////////////
		{
			Name:    "repos",
			Aliases: []string{"repo"},
			Usage:   "Commands for enabling and disabling project repositories",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "List enabled/disabled statuses for project repositories",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "owner, o",
							Value:       "",
							Usage:       "Organization or user that owns the GitHub repository",
							Destination: &owner,
						},
					},
					Action: func(c *cli.Context) error {
						if len(owner) > 0 {
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							r, e := ListRepos(owner, http, cfg)
							pi.Stop()
							if e != nil {
								return cli.NewExitError("Unable to list project statuses. Sorry!", 1)
							} else {
								PrintListRepos(r)
								return nil
							}
						} else {
							return cli.NewExitError("You must supply a --owner or -o argument to specify the repository owner", 1)
						}
					},
				},
				{
					Name:  "enable",
					Usage: "Enable a project repository for use with Nelson",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "repository, repo, r",
							Value:       "",
							Usage:       "GitHub repository to be enabled",
							Destination: &repository,
						},
						cli.StringFlag{
							Name:        "owner, o",
							Value:       "",
							Usage:       "Organization or user that owns the GitHub repository",
							Destination: &owner,
						},
					},
					Action: func(c *cli.Context) error {
						if len(owner) > 0 {
							if len(repository) > 0 {
								req := EnableRepoRequest{
									Owner: owner,
									Repo:  repository,
								}
								pi.Start()
								cfg := LoadDefaultConfigOrExit(http)
								r, e := Enable(req, http, cfg)
								pi.Stop()
								if e != nil {
									return cli.NewExitError("Unable to enable project "+req.Owner+"/"+req.Repo+". Response was:\n"+r, 1)
								} else {
									fmt.Println(r)
								}
							} else {
								return cli.NewExitError("You must supply a --repository or --repo or -r argument to specify the repository", 1)
							}
						} else {
							return cli.NewExitError("You must supply a --owner or -o argument to specify the repository owner", 1)
						}
						return nil
					},
				},
				{
					Name:  "disable",
					Usage: "Disable a project repository for use with Nelson",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "repository, repo, r",
							Value:       "",
							Usage:       "GitHub repository to be disabled",
							Destination: &repository,
						},
						cli.StringFlag{
							Name:        "owner, o",
							Value:       "",
							Usage:       "Organization or user that owns the GitHub repository",
							Destination: &owner,
						},
					},
					Action: func(c *cli.Context) error {
						if len(owner) > 0 {
							if len(repository) > 0 {
								req := EnableRepoRequest{
									Owner: owner,
									Repo:  repository,
								}
								pi.Start()
								cfg := LoadDefaultConfigOrExit(http)
								r, e := Disable(req, http, cfg)
								pi.Stop()
								if e != nil {
									return cli.NewExitError("Unable to disable project "+req.Owner+"/"+req.Repo+". Response was:\n"+r, 1)
								} else {
									fmt.Println(r)
								}
							} else {
								return cli.NewExitError("You must supply a --repository or --repo or -r argument to specify the repository", 1)
							}
						} else {
							return cli.NewExitError("You must supply a --owner or -o argument to specify the repository owner", 1)
						}
						return nil
					},
				},
			},
		},
		////////////////////////////// UNITS //////////////////////////////////
		{
			Name:    "units",
			Aliases: []string{"unit"},
			Usage:   "Set of commands to obtain details about logical deployment units",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the available units",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "datacenters, d",
							Value:       "",
							Usage:       "Restrict list of units to a particular datacenter",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespaces, ns",
							Value:       "",
							Usage:       "Restrict list of units to a particular namespace",
							Destination: &selectedNamespace,
						},
						cli.StringFlag{
							Name:        "statuses, s",
							Value:       "",
							Usage:       "Restrict list of units to a particular status. Defaults to 'ready,warming,manual'",
							Destination: &selectedStatus,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 {
							if !isValidCommaDelimitedList(selectedDatacenter) {
								return cli.NewExitError("You supplied an argument for 'datacenters' but it was not a valid comma-delimited list.", 1)
							}
						}
						if len(selectedNamespace) > 0 {
							if !isValidCommaDelimitedList(selectedNamespace) {
								return cli.NewExitError("You supplied an argument for 'namespaces' but it was not a valid comma-delimited list.", 1)
							}
						} else {
							return cli.NewExitError("You must supply --namespaces or -ns argument to specify the namesapce(s) as a comma delimted form. i.e. dev,qa,prod or just dev", 1)
						}
						if len(selectedStatus) > 0 {
							if !isValidCommaDelimitedList(selectedStatus) {
								return cli.NewExitError("You supplied an argument for 'statuses' but it was not a valid comma-delimited list.", 1)
							}
						}

						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						us, errs := ListUnits(selectedDatacenter, selectedNamespace, selectedStatus, http, cfg)
						pi.Stop()
						if errs != nil {
							return cli.NewExitError("Unable to list units", 1)
						} else {
							PrintListUnits(us)
						}
						return nil
					},
				},
				{
					Name:  "commit",
					Usage: "Commit a unit@version combination to a specific target namespace.",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "unit, u",
							Value:       "",
							Usage:       "The specific unit you want to deprecate.",
							Destination: &selectedUnitPrefix,
						},
						cli.StringFlag{
							Name:        "version, v",
							Value:       "",
							Usage:       "The feature version series you want to deprecate. For example 1.2.3 or 5.3.12",
							Destination: &selectedVersion,
						},
						cli.StringFlag{
							Name:        "target, t",
							Value:       "",
							Usage:       "The target namespace you want to commit this unit too.",
							Destination: &selectedNamespace,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedUnitPrefix) > 0 && len(selectedVersion) > 0 {
							match, _ := regexp.MatchString("(\\d+)\\.(\\d+).(\\d+)", selectedVersion)
							if match == true {
								req := CommitRequest{
									UnitName: selectedUnitPrefix,
									Version:  selectedVersion,
									Target:   selectedNamespace,
								}

								pi.Start()
								cfg := LoadDefaultConfigOrExit(http)
								r, e := CommitUnit(req, http, cfg)
								pi.Stop()

								unitWithVersion := selectedUnitPrefix + "@" + selectedVersion

								if e != nil {
									return cli.NewExitError("Unable to commit "+unitWithVersion+" to '"+selectedNamespace+"'. Response was:\n"+r, 1)
								} else {
									fmt.Println("===>> Commited " + unitWithVersion + " to '" + selectedNamespace + "'.")
								}
							} else {
								return cli.NewExitError("You must supply a version of the format XXX.XXX.XXX, e.g. 2.3.4, 4.56.6, 1.7.9", 1)
							}
						} else {
							return cli.NewExitError("Required --unit, --version or --target inputs were not valid", 1)
						}
						return nil
					},
				},
				{
					Name:  "inspect",
					Usage: "Display details about a logical unit",
					Action: func(c *cli.Context) error {
						fmt.Println("Currently not implemented.")
						return nil
					},
				},
				{
					Name:  "deprecate",
					Usage: "Deprecate a unit/version combination",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:        "no-grace, n",
							Usage:       "expire this unit immedietly rather than allowing the usual grace period",
							Destination: &selectedNoGrace,
						},
						cli.StringFlag{
							Name:        "unit, u",
							Value:       "",
							Usage:       "The unit you want to deprecate",
							Destination: &selectedUnitPrefix,
						},
						cli.StringFlag{
							Name:        "version, v",
							Value:       "",
							Usage:       "The feature version series you want to deprecate",
							Destination: &selectedVersion,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedUnitPrefix) > 0 && len(selectedVersion) > 0 {
							match, _ := regexp.MatchString("^(\\d+)\\.(\\d+)$", selectedVersion)
							if match == true {
								splitVersion := strings.Split(selectedVersion, ".")
								mjr, _ := strconv.Atoi(splitVersion[0])
								min, _ := strconv.Atoi(splitVersion[1])
								ver := FeatureVersion{
									Major: mjr,
									Minor: min,
								}
								req := DeprecationExpiryRequest{
									ServiceType: selectedUnitPrefix,
									Version:     ver,
								}
								pi.Start()
								cfg := LoadDefaultConfigOrExit(http)
								r, e := Deprecate(req, http, cfg)
								pi.Stop()

								if e != nil {
									return cli.NewExitError("Unable to deprecate unit+version series. Response was:\n"+r, 1)
								} else {
									if selectedNoGrace == true {
										_, e2 := Expire(req, http, cfg)
										if e2 != nil {
											return cli.NewExitError("Unable to deprecate unit+version series. Response was:\n"+r, 1)
										} else {
											fmt.Println("===>> Deprecated and expired " + selectedUnitPrefix + " " + selectedVersion)
										}
									} else {
										fmt.Println("===>> Deprecated " + selectedUnitPrefix + " " + selectedVersion)
									}
								}
							} else {
								return cli.NewExitError("You must supply a feature version of the format XXX.XXX, e.g. 2.3, 4.56, 1.7", 1)
							}
						} else {
							return cli.NewExitError("Required --unit and/or --version inputs were not valid", 1)
						}
						return nil
					},
				},
			},
		},
		////////////////////////////// STACK //////////////////////////////////
		{
			Name:    "stacks",
			Aliases: []string{"stack"},
			Usage:   "Set of commands to obtain details about deployed stacks",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the available stacks",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "unit, u",
							Value:       "",
							Usage:       "Only list stacks that match a specified unit prefix",
							Destination: &selectedUnitPrefix,
						},
						cli.StringFlag{
							Name:        "datacenters, d",
							Value:       "",
							Usage:       "Restrict list of units to a particular datacenter",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespaces, ns",
							Value:       "",
							Usage:       "Restrict list of units to a particular namespace",
							Destination: &selectedNamespace,
						},
						cli.StringFlag{
							Name:        "statuses, s",
							Value:       "",
							Usage:       "Restrict list of units to a particular status. Defaults to 'ready,manual'",
							Destination: &selectedStatus,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 {
							if !isValidCommaDelimitedList(selectedDatacenter) {
								return cli.NewExitError("You supplied an argument for 'datacenters' but it was not a valid comma-delimited list.", 1)
							}
						}
						if len(selectedNamespace) > 0 {
							if !isValidCommaDelimitedList(selectedNamespace) {
								return cli.NewExitError("You supplied an argument for 'namespaces' but it was not a valid comma-delimited list.", 1)
							}
						} else {
							return cli.NewExitError("You must supply --namespaces or -ns argument to specify the namesapce(s) as a comma delimted form. i.e. dev,qa,prod or just dev", 1)
						}
						if len(selectedStatus) > 0 {
							if !isValidCommaDelimitedList(selectedStatus) {
								return cli.NewExitError("You supplied an argument for 'statuses' but it was not a valid comma-delimited list.", 1)
							}
						}

						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						r, e := ListStacks(selectedDatacenter, selectedNamespace, selectedStatus, http, cfg)
						pi.Stop()
						if e != nil {
							PrintTerminalErrors(e)
							return cli.NewExitError("Unable to list stacks.", 1)
						} else {
							PrintListStacks(r)
						}
						return nil
					},
				},
				{
					Name:  "inspect",
					Usage: "Display the current status and details about a specific stack",
					Action: func(c *cli.Context) error {
						guid := c.Args().First()
						if len(guid) > 0 && isValidGUID(guid) {
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							r, e := InspectStack(guid, http, cfg)
							pi.Stop()
							if e != nil {
								PrintTerminalErrors(e)
								return cli.NewExitError("Unable to inspect stacks '"+guid+"'.", 1)
							} else {
								PrintInspectStack(r)
							}
						} else {
							return cli.NewExitError("You must supply a valid GUID for the stack you want to inspect.", 1)
						}
						return nil
					},
				},
				{
					Name:  "runtime",
					Usage: "Display the runtime status for a particular stack",
					Action: func(c *cli.Context) error {
						guid := c.Args().First()
						if isValidGUID(guid) {
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							r, e := GetStackRuntime(guid, http, cfg)
							pi.Stop()
							if e != nil {
								PrintTerminalErrors(e)
								return cli.NewExitError("Unable to fetch runtime status.", 1)
							} else {
								PrintStackRuntime(r)
							}
						} else {
							return cli.NewExitError("You must specify a valid GUID reference in order to display runtime status of a stack.", 1)
						}
						return nil
					},
				},
				{
					Name:  "redeploy",
					Usage: "Trigger a redeployment for a specific stack",
					Action: func(c *cli.Context) error {
						guid := c.Args().First()
						if isValidGUID(guid) {
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							r, e := Redeploy(guid, http, cfg)
							pi.Stop()

							if e != nil {
								return cli.NewExitError("Unable to request a redeploy. Response was:\n"+r, 1)
							} else {
								fmt.Println("===>> " + r)
							}
						} else {
							return cli.NewExitError("You must specify a valid GUID reference in order to redeploy a stack.", 1)
						}
						return nil
					},
				},
				{
					Name:  "manual",
					Usage: "Register a manual deployment",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "datacenter, dc",
							Value:       "",
							Usage:       "The datacenter for the service",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespace, ns",
							Value:       "",
							Usage:       "The namespace for the service",
							Destination: &selectedNamespace,
						},
						cli.StringFlag{
							Name:        "service-type, st",
							Value:       "",
							Usage:       "The service type for the service",
							Destination: &selectedServiceType,
						},
						cli.StringFlag{
							Name:        "version, v",
							Value:       "",
							Usage:       "The version for the service",
							Destination: &selectedVersion,
						},
						cli.StringFlag{
							Name:        "hash",
							Value:       "",
							Usage:       "The hash for the stack",
							Destination: &stackHash,
						},
						cli.StringFlag{
							Name:        "description, d",
							Value:       "",
							Usage:       "Description for the service",
							Destination: &description,
						},
						cli.Int64Flag{
							Name:        "port",
							Value:       0,
							Usage:       "The exposed port for the service",
							Destination: &selectedPort,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 && len(selectedNamespace) > 0 &&
							len(selectedServiceType) > 0 &&
							len(selectedVersion) > 0 &&
							len(stackHash) > 0 &&
							len(description) > 0 &&
							selectedPort > 0 {
							req := ManualDeploymentRequest{
								Datacenter:  selectedDatacenter,
								Namespace:   selectedNamespace,
								ServiceType: selectedServiceType,
								Version:     selectedVersion,
								Hash:        stackHash,
								Description: description,
								Port:        selectedPort,
							}
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							res, e := RegisterManualDeployment(req, http, cfg)
							pi.Stop()
							if e != nil {
								return cli.NewExitError("Unable to register manual deployment.", 1)
							} else {
								fmt.Println(res)
							}
						} else {
							return cli.NewExitError("You must specify the following switches: \n\t--datacenter <string> \n\t--namespace <string> \n\t--service-type <string> \n\t--version <string> \n\t--hash <string> \n\t--description <string> \n\t--port <int>", 1)
						}
						return nil
					},
				},
				{
					Name:    "fs",
					Aliases: []string{"logs"},
					Usage:   "Fetch the deployment log for a given stack",
					Action: func(c *cli.Context) error {
						guid := c.Args().First()
						if len(guid) > 0 && isValidGUID(guid) {
							cfg := LoadDefaultConfigOrExit(http)
							GetDeploymentLog(guid, http, cfg)
						} else {
							return cli.NewExitError("You must specify a valid GUID for the stack you wish to view logs for.", 1)
						}
						return nil
					},
				},
			},
		},
		////////////////////////////// SYSTEM //////////////////////////////////
		{
			Name:  "system",
			Usage: "A set of operations to query Nelson to see what options are available",
			Subcommands: []cli.Command{
				{
					Name:  "cleanup-policies",
					Usage: "list the available cleanup policies",
					Action: func(c *cli.Context) error {
						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						policies, e := ListCleanupPolicies(http, cfg)
						pi.Stop()
						if e != nil {
							PrintTerminalErrors(e)
							return cli.NewExitError("Unable to list the cleanup policies at this time.", 1)
						} else {
							PrintCleanupPolicies(policies)
						}
						return nil
					},
				},
				{
					Name:  "version",
					Usage: "Ask for info about the current Nelson build",
					Action: func(c *cli.Context) error {
						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						sr, e := WhoAreYou(http, cfg)
						pi.Stop()
						if e != nil {
							PrintTerminalErrors(e)
							return cli.NewExitError("Unable to fetch build info for Nelson.", 1)
						} else {
							fmt.Println(sr.Banner)
							fmt.Println(" " + cfg.Endpoint)
						}
						return nil
					},
				},
			},
		},
		////////////////////////////// WHOAMI //////////////////////////////////
		{
			Name:  "whoami",
			Usage: "Ask nelson who you are currently logged in as",
			Action: func(c *cli.Context) error {
				pi.Start()
				cfg := LoadDefaultConfigOrExit(http)
				sr, e := WhoAmI(http, cfg)
				pi.Stop()
				if e != nil {
					PrintTerminalErrors(e)
					return cli.NewExitError("Unable to determine who is currently logged into Nelson.", 1)
				} else {
					fmt.Println("===>> Currently logged in to " + sr.User.Name + " @ " + cfg.Endpoint)
				}
				return nil
			},
		},
		/////////////////////////// LOADBALANCERS //////////////////////////////
		{
			Name:    "loadbalancers",
			Aliases: []string{"lbs", "loadbalancer"},
			Usage:   "Set of commands to obtain details about available load balancers",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the available loadbalancers",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "datacenters, d",
							Value:       "",
							Usage:       "Restrict list of loadbalancers to a particular datacenter",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespaces, ns",
							Value:       "",
							Usage:       "Restrict list of loadbalancers to a particular namespace",
							Destination: &selectedNamespace,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 {
							if !isValidCommaDelimitedList(selectedDatacenter) {
								return cli.NewExitError("You supplied an argument for 'datacenters' but it was not a valid comma-delimited list.", 1)
							}
						}
						if len(selectedNamespace) > 0 {
							if !isValidCommaDelimitedList(selectedNamespace) {
								return cli.NewExitError("You supplied an argument for 'namespaces' but it was not a valid comma-delimited list.", 1)
							}
						} else {
							return cli.NewExitError("You must supply --namespaces or -ns argument to specify the namesapce(s) as a comma delimted form. i.e. dev,qa,prod or just dev", 1)
						}

						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						us, errs := ListLoadbalancers(selectedDatacenter, selectedNamespace, selectedStatus, http, cfg)
						pi.Stop()
						if errs != nil {
							return cli.NewExitError("Unable to list load balancers right now. Sorry!", 1)
						} else {
							PrintListLoadbalancers(us)
						}
						return nil
					},
				},
				{
					Name:  "down",
					Usage: "remove the specified load balancer",
					Action: func(c *cli.Context) error {
						guid := c.Args().First()
						if len(guid) > 0 && isValidGUID(guid) {
							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							r, e := RemoveLoadBalancer(guid, http, cfg)
							pi.Stop()
							if e != nil {
								PrintTerminalErrors(e)
								return cli.NewExitError("Unable to remove loadbalancer '"+guid+"'.", 1)
							} else {
								fmt.Println("==>>> " + r)
							}
						} else {
							return cli.NewExitError("You must supply a valid GUID for the loadbalancer you want to remove.", 1)
						}
						return nil
					},
				},
				{
					Name:  "up",
					Usage: "start a load balancer",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "datacenter, dc",
							Value:       "",
							Usage:       "The datacenter for the service",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespace, ns",
							Value:       "",
							Usage:       "The namespace for the service",
							Destination: &selectedNamespace,
						},
						cli.StringFlag{
							Name:        "name, n",
							Value:       "",
							Usage:       "The name for the loadbalancer - should include your unit name e.g. howdy-http-lb",
							Destination: &selectedUnitPrefix,
						},
						cli.StringFlag{
							Name:        "major-version, mv",
							Value:       "",
							Usage:       "The major version ",
							Destination: &selectedVersion,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 &&
							len(selectedNamespace) > 0 &&
							len(selectedUnitPrefix) > 0 &&
							len(selectedVersion) > 0 {

							mjver, err := strconv.ParseInt(selectedVersion, 10, 64)
							if err != nil {
								return cli.NewExitError("The specified major version does not look like an integer value.", 1)
							}

							req := LoadbalancerCreate{
								Name:         selectedUnitPrefix,
								MajorVersion: int(mjver),
								Datacenter:   selectedDatacenter,
								Namespace:    selectedNamespace,
							}

							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							res, e := CreateLoadBalancer(req, http, cfg)
							pi.Stop()
							if e != nil {
								PrintTerminalErrors(e)
								return cli.NewExitError("Unable to launch the specified loadbalancer.", 1)
							} else {
								fmt.Println(res)
							}
						} else {
							return cli.NewExitError("You must specify the following switches: \n\t--datacenter <string> \n\t--namespace <string> \n\t--major-version <int> \n\t--name <string>", 1)
						}
						return nil
					},
				},
				{
					Name:  "inspect",
					Usage: "inspect the specified loadbalancer",
					Action: func(c *cli.Context) error {
						selectedLoadbalancer = c.Args().First()
						if len(selectedLoadbalancer) == 0 {
							return cli.NewExitError("you must specify a loadbalancer guid as the first argument", 1)
						}
						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						lb, e := InspectLoadBalancer(selectedLoadbalancer, http, cfg)
						pi.Stop()
						if e != nil {
							PrintTerminalErrors(e)
							return cli.NewExitError("Unable to inspect loadbalancer right now, Sorry!", 1)
						} else {
							PrintInspectLoadbalancer(lb)
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "namespaces",
			Aliases: []string{"ns", "namespace"},
			Usage:   "Set of commands to obtain details about available namespaces",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create namespace",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "datacenter, dc",
							Value:       "",
							Usage:       "The datacenter for the namespace",
							Destination: &selectedDatacenter,
						},
						cli.StringFlag{
							Name:        "namespace, ns",
							Value:       "",
							Usage:       "The namespace to create",
							Destination: &selectedNamespace,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedDatacenter) > 0 &&
							len(selectedNamespace) > 0 {

							req := NamespaceRequest{
								Namespace: selectedNamespace,
							}

							pi.Start()
							cfg := LoadDefaultConfigOrExit(http)
							res, e := CreateNamespace(req, selectedDatacenter, http, cfg)
							pi.Stop()
							if e != nil {
								PrintTerminalErrors(e)
								return cli.NewExitError("Unable to create the specified namespace.", 1)
							} else {
								fmt.Println(res)
							}
						} else {
							return cli.NewExitError("You must specify the following switches: \n\t--datacenter <string> \n\t--namespace <string>", 1)
						}
						return nil
					},
				},
			},
		},
		/////////////////////////// LINT //////////////////////////////
		{
			Name:  "lint",
			Usage: "Set of commands to lint aspects of your deployment",
			Subcommands: []cli.Command{
				{
					Name:  "manifest",
					Usage: "Test whether a Nelson manifest file is valid",
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "unit, u",
							Usage: "Units to be deployed",
						},
						cli.StringFlag{
							Name:        "manifest, m",
							Value:       "",
							Usage:       "The Nelson manifest file to validate",
							Destination: &selectedManifest,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedManifest) <= 0 {
							selectedManifest = ".nelson.yml"
						}
						manifest, err := ioutil.ReadFile(selectedManifest)
						if err != nil {
							return cli.NewExitError("Could not read "+selectedManifest, 1)
						}
						manifestBase64 := base64.StdEncoding.EncodeToString(manifest)
						var unitNames []string = c.StringSlice("unit")
						var manifestUnits []ManifestUnit = []ManifestUnit{}
						for i := 0; i < len(unitNames); i++ {
							var n string = unitNames[i]
							manifestUnits = append(
								manifestUnits,
								ManifestUnit{
									Name: n,
									Kind: n,
								},
							)
						}
						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						req := LintManifestRequest{
							Units:    manifestUnits,
							Manifest: manifestBase64,
						}
						msg, errs := LintManifest(req, http, cfg)
						pi.Stop()
						if errs != nil {
							PrintTerminalErrors(errs)
							fmt.Println(msg)
							return cli.NewExitError("Manifest validation failed.", 1)
						} else {
							fmt.Println(msg)
						}
						return nil
					},
				},
				{
					Name:  "template",
					Usage: "Test whether a template will render in your container",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "unit, u",
							Value:       "",
							Usage:       "The unit name that owns the template (e.g., howdyhttp)",
							Destination: &selectedUnitPrefix,
						},
						cli.StringSliceFlag{
							Name:  "resource, r",
							Usage: "resources required by this template (e.g., s3); repeatable",
						},
						cli.StringFlag{
							Name:        "template, t",
							Value:       "",
							Usage:       "The file name containing the template to lint",
							Destination: &selectedTemplate,
						},
					},
					Action: func(c *cli.Context) error {
						if len(selectedUnitPrefix) <= 0 {
							return cli.NewExitError("You must specify a unit name for the template to be linted.", 1)
						}

						if len(selectedTemplate) <= 0 {
							return cli.NewExitError("You must specify a template file to lint.", 1)
						}
						template, err := ioutil.ReadFile(selectedTemplate)
						if err != nil {
							return cli.NewExitError("Could not read "+selectedTemplate, 1)
						}
						templateBase64 := base64.StdEncoding.EncodeToString(template)

						pi.Start()
						cfg := LoadDefaultConfigOrExit(http)
						req := LintTemplateRequest{
							Unit:      selectedUnitPrefix,
							Resources: c.StringSlice("resource"),
							Template:  templateBase64,
						}
						msg, errs := LintTemplate(req, http, cfg)
						pi.Stop()
						if errs != nil {
							PrintTerminalErrors(errs)
							fmt.Println(msg)
							return cli.NewExitError("Template linting failed.", 1)
						} else {
							fmt.Println(msg)
						}
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
