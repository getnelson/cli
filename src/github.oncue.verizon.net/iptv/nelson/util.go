package main

/// had to have a `util` file.
/// Because pragmatism.
/// Because irony.

import (
  "os"
  "fmt"
  "time"
  "github.com/parnurzeal/gorequest"
  "github.com/olekukonko/tablewriter"
  "github.com/briandowns/spinner"
)

func AugmentRequest(c *gorequest.SuperAgent, cfg *Config) *gorequest.SuperAgent {
  return c.
    AddCookie(cfg.GetAuthCookie()).
    Set("Content-type","application/json").
    Timeout(15*time.Second).
    SetCurlCommand(false).
    SetDebug(globalEnableDebug)
}

func RenderTableToStdout(headers []string, data [][]string){
  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader(headers)
  table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
  table.SetHeaderLine(false)
  table.SetRowLine(false)
  table.SetColWidth(100)
  table.SetColumnSeparator("")
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetAlignment(tablewriter.ALIGN_LEFT)
  table.AppendBulk(data) // Add Bulk Data
  table.Render()
}

func JavaEpochToDateStr(long int64) string {
  t := time.Unix(0, long*int64(time.Millisecond))
  return t.Format(time.RFC3339)
}

func ProgressIndicator() *spinner.Spinner {
  s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
  s.Color("green")
  return s
}

func PrintTerminalErrors(errs []error){
  for i, j := 0, len(errs)-1; i < j; i, j = i+1, j-1 {
    errs[i], errs[j] = errs[j], errs[i]
  }

  for _,e := range errs {
    fmt.Println(e)
  }
}