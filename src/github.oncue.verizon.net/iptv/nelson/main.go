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
      Usage:       "Set of commands for working with Nelson regions",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "List all the available regions",
          Action: func(c *cli.Context) error {
            ListRegions(http, LoadDefaultConfig())
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "Dispaly details about a specified region",
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
      Usage:       "Set of commands to obtain details about logical deployment units",
      Subcommands: []cli.Command{
        {
          Name:  "list",
          Usage: "list the available regions",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name:   "region, r",
              Value:  "",
              Usage:  "Restrict list of units to a particular region",
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
          Usage: "Display details about a logical unit",
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
          },
          Action: func(c *cli.Context) error {
            fmt.Println("Not Implemented")
            return nil
          },
        },
        {
          Name:  "inspect",
          Usage: "Display the current status and details about a specific stack",
          Action: func(c *cli.Context) error {
            fmt.Println("inspecting stack: ", c.Args().First())
            return nil
          },
        },
        {
          Name:  "fs",
          Usage: "Fetch the deployment log for a given stack",
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
    ////////////////////////////// WHOAMI //////////////////////////////////
    {
      Name:    "whoami",
      Usage:   "Ask nelson who you are currently logged in as",
      Action:  func(c *cli.Context) error {
        WhoAmI(http, LoadDefaultConfig())
        return nil
      },
    },
  }

  app.Run(os.Args)
}
