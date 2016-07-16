package main

import (
  "os"
  "fmt"
  "time"
  "strings"
  "strconv"
  "gopkg.in/urfave/cli.v1"
  "github.com/parnurzeal/gorequest"
)

var globalEnableDebug bool

func main() {
  year, _, _ := time.Now().Date()
  app := cli.NewApp()
  app.Name = "nelson-cli"
  app.Version = "v0.2"
  app.Copyright = "© "+strconv.Itoa(year)+" Verizon Labs"
  app.Usage = "remote control for the Nelson deployment system"

  http := gorequest.New()
  pi   := ProgressIndicator()

  // switches for the cli
  var userGithubToken string
  var disableTLS bool
  var selectedDatacenter string
  var selectedNamespace string
  var selectedStatus string
  var selectedUnitPrefix string
  var selectedVersion string

  app.Flags = []cli.Flag {
    cli.BoolFlag{
      Name: "debug",
      Usage: "Enable debug mode on the network requests",
      Destination: &globalEnableDebug,
    },
  }

  app.Commands = []cli.Command {
    ////////////////////////////// LOGIN //////////////////////////////////
    {
      Name:    "login",
      Usage:   "login to nelson",
      Flags: []cli.Flag {
        cli.StringFlag{
          Name:   "token, t",
          Value:  "",
          Usage:  "your github personal access token",
          EnvVar: "GITHUB_TOKEN",
          Destination: &userGithubToken,
        },
        cli.BoolFlag {
          Name:  "disable-tls",
          Destination: &disableTLS,
        },
      },
      Action:  func(c *cli.Context) error {
        host := strings.TrimSpace(c.Args().First())
        if len(host) <= 0 {
          host = os.Getenv("NELSON_ADDR")
          if len(host) <= 0 {
            return cli.NewExitError("Either supply a host explicitly, or set $NELSON_ADDR", 1)
          }
        }

        if len(userGithubToken) <= 0 {
          return cli.NewExitError("You must specifiy a `--token` or a `-t` to login.", 1)
        }

        // fmt.Println("token: ", userGithubToken)
        // fmt.Println("host: ", host)
        pi.Start()
        Login(http, userGithubToken, host, disableTLS)
        pi.Stop()
        fmt.Println("Sucsessfully logged in to " + host)
        return nil
      },
    },
    ////////////////////////////// DATACENTER //////////////////////////////////
    {
      Name:        "datacenters",
      Aliases:     []string{"dcs"},
      Usage:       "Set of commands for working with Nelson datacenters",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "List all the available datacenters",
          Action: func(c *cli.Context) error {
            pi.Start()
            r, e := ListDatacenters(http, LoadDefaultConfig())
            pi.Stop()
            if e != nil {
              return cli.NewExitError("Unable to list datacenters.", 1)
            } else {
              PrintListDatacenters(r)
            }
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "Dispaly details about a specified datacenter",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting datacenter: ", c.Args().First())
            return nil
          },
        },
      },
    },
    ////////////////////////////// UNITS //////////////////////////////////
    {
      Name:        "units",
      Usage:       "Set of commands to obtain details about logical deployment units",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available units",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "datacenter, d",
              Value:  "",
              Usage:  "Restrict list of units to a particular datacenter",
              Destination: &selectedDatacenter,
            },
            cli.StringFlag{
              Name:   "namespace, ns",
              Value:  "",
              Usage:  "Restrict list of units to a particular namespace",
              Destination: &selectedNamespace,
            },
            cli.StringFlag{
              Name:   "status, s",
              Value:  "",
              Usage:  "Restrict list of units to a particular status. Defaults to 'active,manual'",
              Destination: &selectedStatus,
            },
          },
          Action: func(c *cli.Context) error {
            if len(selectedDatacenter) > 0 {
              pi.Start()
              us, errs := ListUnits(selectedDatacenter, http, LoadDefaultConfig())
              pi.Stop()
              if(errs != nil){
                return cli.NewExitError("Unable to list units", 1)
              } else {
                PrintListUnits(us)
              }
            } else {
              return cli.NewExitError("Missing --datacenter flag; cannot list units for all datacenters in one request.", 1)
            }
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "Display details about a logical unit",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting unit: ", c.Args().First())
            return nil
          },
        },
      },
    },
    ////////////////////////////// STACK //////////////////////////////////
    {
      Name:        "stacks",
      Usage:       "Set of commands to obtain details about deployed stacks",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available stacks",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "unit, u",
              Value:  "",
              Usage:  "Only list stacks that match a specified unit prefix",
              Destination: &selectedUnitPrefix,
            },
            cli.StringFlag{
              Name:   "datacenter, d",
              Value:  "",
              Usage:  "Only list stacks that reside in a given datacenter",
              Destination: &selectedDatacenter,
            },
          },
          Action: func(c *cli.Context) error {
            if len(selectedDatacenter) > 0 {
              pi.Start()
              r, e := ListStacks(selectedDatacenter, http, LoadDefaultConfig())
              pi.Stop()
              if e != nil {
                PrintTerminalErrors(e)
                return cli.NewExitError("Unable to list stacks.", 1)
              } else {
                PrintListStacks(r)
              }

            } else {
              return cli.NewExitError("You must suppled --datacenter in order to list stacks", 1)
            }
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "Display the current status and details about a specific stack",
          Action: func(c *cli.Context) error {
            guid := c.Args().First()
            if len(guid) > 0 {
              pi.Start()
              r, e := InspectStack(guid, http, LoadDefaultConfig())
              pi.Stop()
              if e != nil {
                PrintTerminalErrors(e)
                return cli.NewExitError("Unable to inspect stacks '"+guid+"'.", 1)
              } else {
                PrintInspectStack(r)
              }
            } else {
              return cli.NewExitError("You must supply a GUID of the stack you want to inspect.", 1)
            }
            return nil
          },
        },
        {
          Name:  "redeploy",
          Usage: "Trigger a redeployment for a specific stack",
          Action: func(c *cli.Context) error {
            guid := c.Args().First()

            r,e := Redeploy(guid, http, LoadDefaultConfig())

            if e != nil {
              return cli.NewExitError("Unable to request a redeploy. Response was:\n"+r, 1)
            } else {
              fmt.Println("===>> "+r)
            }
            return nil
          },
        },
        {
          Name:  "deprecate",
          Usage: "Deprecate a unit/version combination (and all patch series)",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "unit, u",
              Value:  "",
              Usage:  "The unit you want to deprecate",
              Destination: &selectedUnitPrefix,
            },
            cli.StringFlag{
              Name:   "version, v",
              Value:  "",
              Usage:  "The feature version series you want to deprecate",
              Destination: &selectedVersion,
            },
          },
          Action: func(c *cli.Context) error {
            if len(selectedUnitPrefix) > 0 && len(selectedVersion) > 0 {
              splitVersion := strings.Split(selectedVersion, ".")
              mjr, _ := strconv.Atoi(splitVersion[0])
              min, _ := strconv.Atoi(splitVersion[1])
              ver := FeatureVersion {
                Major: mjr,
                Minor: min,
              }
              req := DeprecationRequest {
                ServiceType: selectedUnitPrefix,
                Version: ver,
              }
              r,e := Deprecate(req, http, LoadDefaultConfig())
              if e != nil {
                return cli.NewExitError("Unable to deprecate unit+version series. Response was:\n"+r, 1)
              } else {
                fmt.Println("===>> Deprecated "+selectedUnitPrefix+" "+selectedVersion)
              }
            } else {
              return cli.NewExitError("Required --unit and/or --version inputs were not valid", 1)
            }
            return nil
          },
        },
        {
          Name:  "fs",
          Usage: "Fetch the deployment log for a given stack",
          Action: func(c *cli.Context) error {
            guid := c.Args().First()
            if len(guid) > 0 {
              GetDeploymentLog(guid, http, LoadDefaultConfig())
            } else {
              return cli.NewExitError("You must specify which stack you would like to redeploy.", 1)
            }
            return nil
          },
        },
      },
    },
    ////////////////////////////// SYSTEM //////////////////////////////////
    {
      Name:    "system",
      Usage:   "A set of operations to query Nelson to see what options are available",
      Subcommands: []cli.Command{
        {
          Name:  "cleanup-policies",
          Usage: "list the available cleanup policies",
          Action: func(c *cli.Context) error {
            pi.Start()
            policies, e := ListCleanupPolicies(http, LoadDefaultConfig())
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
      },
    },
    ////////////////////////////// WHOAMI //////////////////////////////////
    {
      Name:    "whoami",
      Usage:   "Ask nelson who you are currently logged in as",
      Action:  func(c *cli.Context) error {
        pi.Start()
        cfg := LoadDefaultConfig()
        sr, e := WhoAmI(http, cfg)
        pi.Stop()
        if e != nil {
          return cli.NewExitError("Unable to list stacks.", 1)
        } else {
          fmt.Println("===>> Currently logged in as "+sr.User.Name+" @ "+cfg.Endpoint)
        }
        return nil
      },
    },
  }

  app.Run(os.Args)
}
