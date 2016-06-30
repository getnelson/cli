package main

import (
  "os"
  // "fmt"
  "time"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
  "github.com/olekukonko/tablewriter"
)

type Region struct {
  Name string `json:"name"`
  Namespaces []Namespace `json:"namespaces"`
}
type Namespace struct {
  Id int `json:"id"`
  Name string `json:"name"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func ListRegions(http *gorequest.SuperAgent, cfg *Config){
  _, bytes, errs := http.
    Get(cfg.Endpoint+"/v1/datacenters").
    AddCookie(cfg.GetAuthCookie()).
    Set("Content-type","application/json").
    Timeout(10*time.Second).
    EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var regions []Region
  if err := json.Unmarshal(bytes, &regions); err != nil {
    panic(err)
  }

  var tabulized = [][]string {}
  for _,r := range regions {
    namespace := ""
    for i,ns := range r.Namespaces {
      if(i == 0){
        namespace = ns.Name
      } else {
        namespace = namespace+", "+ns.Name
      }
    }
    tabulized = append(tabulized,[]string{ r.Name, namespace })
  }

  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{ "Region", "Namespaces" })
  table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
  table.SetHeaderLine(false)
  table.SetRowLine(false)
  table.SetColumnSeparator("")
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetAlignment(tablewriter.ALIGN_LEFT)
  table.AppendBulk(tabulized) // Add Bulk Data
  table.Render()
}