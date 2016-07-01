package main

/// had to have a `util` file.
/// Because pragmatism.
/// Because irony.

import (
  "os"
  "time"
  "github.com/parnurzeal/gorequest"
  "github.com/olekukonko/tablewriter"
)

func AugmentRequest(c *gorequest.SuperAgent, cfg *Config) *gorequest.SuperAgent {
  return c.
    AddCookie(cfg.GetAuthCookie()).
    Set("Content-type","application/json").
    Timeout(2*time.Second)
}

func RenderTableToStdout(headers []string, data [][]string){
  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{ "Region", "Namespaces" })
  table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
  table.SetHeaderLine(false)
  table.SetRowLine(false)
  table.SetColumnSeparator("")
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetAlignment(tablewriter.ALIGN_LEFT)
  table.AppendBulk(data) // Add Bulk Data
  table.Render()
}
