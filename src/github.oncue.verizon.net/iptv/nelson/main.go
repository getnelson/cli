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

func main() {
  year, _, _ := time.Now().Date()
  app := cli.NewApp()
  app.Name = "nelson-cli"
  app.Version = "v0.1"
  app.Copyright = "© "+strconv.Itoa(year)+" Verizon Labs"
  app.Usage = "remote control for the Nelson deployment system"

  http := gorequest.New()

  // switches for the cli
  var userGithubToken string
  var disableTLS bool
  var selectedRegion string
  var selectedUnitPrefix string

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
        Login(http, userGithubToken, host, disableTLS)
        fmt.Println("Sucsessfully logged in to " + host)
        return nil
      },
    },
    ////////////////////////////// REGION //////////////////////////////////
    {
      Name:        "region",
      Usage:       "control nelson regions",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available regions",
          Action: func(c *cli.Context) error {
            ListRegions(http, LoadDefaultConfig())
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "show details about a specified region",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting region: ", c.Args().First())
            return nil
          },
        },
      },
    },
    ////////////////////////////// UNITS //////////////////////////////////
    {
      Name:        "unit",
      Usage:       "control nelson units",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available regions",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "region, r",
              Value:  "",
              Usage:  "only list units in a particular region",
              Destination: &selectedRegion,
            },
          },
          Action: func(c *cli.Context) error {
            fmt.Println("Not Implemented")
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "show details about a specified unit",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting region: ", c.Args().First())
            return nil
          },
        },
      },
    },
    ////////////////////////////// STACK //////////////////////////////////
    {
      Name:        "stack",
      Usage:       "get specific information about a given stack",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available stacks",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "unit, u",
              Value:  "",
              Usage:  "only list stakcs for a specified unit prefix",
              Destination: &selectedUnitPrefix,
            },
          },
          Action: func(c *cli.Context) error {
            fmt.Println("Not Implemented")
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "show current status and details about a specified stack",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting stack: ", c.Args().First())
            return nil
          },
        },
        {
          Name:  "fs",
          Usage: "get the log for a given deployment",
          Action: func(c *cli.Context) error {

            i64, err := strconv.ParseInt(c.Args().First(), 10, 16)
            if err != nil {
              return cli.NewExitError("The supplied argument was not a parsable integer", 1)
            }

            GetDeploymentLog(int(i64), http, LoadDefaultConfig())
            return nil
          },
        },
      },
    },
  }

  app.Run(os.Args)
}
