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
        login(http, userGithubToken, host)
        fmt.Println("Sucsessfully logged in to " + host)
        return nil
      },
    },
    {
      Name:    "complete",
      Aliases: []string{"c"},
      Usage:   "complete a task on the list",
      Action:  func(c *cli.Context) error {
        fmt.Println("completed task: ", c.Args().First())
        return nil
      },
    },
    {
      Name:        "template",
      Aliases:     []string{"t"},
      Usage:       "options for task templates",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new template",
          Action: func(c *cli.Context) error {
            fmt.Println("new task template: ", c.Args().First())
            return nil
          },
        },
        {
          Name:  "remove",
          Usage: "remove an existing template",
          Action: func(c *cli.Context) error {
            fmt.Println("removed task template: ", c.Args().First())
            return nil
          },
        },
      },
    },
  }

  app.Run(os.Args)
}
