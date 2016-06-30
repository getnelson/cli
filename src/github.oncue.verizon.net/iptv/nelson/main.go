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

  var userGithubToken string
  var disableTLS bool

  app.Commands = []cli.Command{
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
  }

  app.Run(os.Args)
}
