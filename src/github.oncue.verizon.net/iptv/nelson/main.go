package main

import (
  "os"
  "fmt"
  "strings"
  "gopkg.in/urfave/cli.v1"
  "github.com/parnurzeal/gorequest"
)

// https://github.com/urfave/cli#exit-code

func main() {
  app := cli.NewApp()
  http := gorequest.New()

  app.Commands = []cli.Command{
    {
      Name:    "login",
      Usage:   "login to nelson",
      Action:  func(c *cli.Context) error {
        user_token := strings.TrimSpace(c.Args().First())
        login(http, user_token)
        fmt.Println("Sucsessfully logged in!")
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
