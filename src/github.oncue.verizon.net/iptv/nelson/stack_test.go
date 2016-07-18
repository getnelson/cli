package main

import (
  "encoding/json"
  "testing"
  // "fmt"
)

func TestStackJsonUnmarshaling(t *testing.T) {
  fixture := `
    {
      "workflow": "manual",
      "guid": "1a69395e919d",
      "stack_name": "dev-iptv-cass-dev--4-8-4--mtq2odqzndc0mg",
      "deployed_at": 1468518896093,
      "unit": "dev-iptv-cass-dev",
      "type": "service"
    }`

  var result Stack
  err := json.Unmarshal([]byte(fixture), &result)
  if result.Workflow != "manual" {
    t.Error("Should be able to parse JSON, but got error:\n", err)
  }
}

func TestStackSummaryJsonUnmarshaling(t *testing.T) {
  fixture := `
  {
  "workflow": "pulsar",
  "guid": "e4184c271bb9",
  "statuses": [
    {
      "timestamp": "2016-07-14T22:30:22.358Z",
      "message": "inventory-inventory deployed to perryman",
      "status": "active"
    },
    {
      "timestamp": "2016-07-14T22:30:22.310Z",
      "message": "writing alert definitions to perryman's consul",
      "status": "deploying"
    },
    {
      "timestamp": "2016-07-14T22:30:22.264Z",
      "message": "waiting for the application to become ready",
      "status": "deploying"
    },
    {
      "timestamp": "2016-07-14T22:30:21.421Z",
      "message": "instructing perryman's marathon to launch container",
      "status": "deploying"
    },
    {
      "timestamp": "2016-07-14T22:29:44.370Z",
      "message": "replicating docker.iptv.oncue.com/units/inventory-inventory-2.0:2.0.11 to remote registry",
      "status": "deploying"
    },
    {
      "timestamp": "2016-07-14T22:29:44.323Z",
      "message": "",
      "status": "pending"
    }
  ],
  "stack_name": "inventory-inventory--2-0-11--8gufie2b",
  "deployed_at": 1468535384221,
  "unit": "inventory-inventory",
  "type": "service",
  "expiration": 1469928212871,
  "dependencies": {
    "outbound": [
      {
        "workflow": "manual",
        "guid": "1a69395e919d",
        "stack_name": "dev-iptv-cass-dev--4-8-4--mtq2odqzndc0mg",
        "deployed_at": 1468518896093,
        "unit": "dev-iptv-cass-dev",
        "type": "service"
      }
    ],
    "inbound": []
  },
  "namespace": "devel"
}`

  var result StackSummary
  err := json.Unmarshal([]byte(fixture), &result)

  t.Log(result.Dependencies.Outbound)

  if len(result.Dependencies.Outbound) != 1 {
    t.Error("Should have had one outbound dependency, but got error:\n", err)
  }
}