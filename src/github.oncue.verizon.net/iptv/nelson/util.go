package main

/// had to have a `util` file.
/// Because pragmatism.
/// Because irony.

import (
	"fmt"
	"github.com/briandowns/spinner"
	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"
)

func AugmentRequest(c *gorequest.SuperAgent, cfg *Config) *gorequest.SuperAgent {
	return c.
		AddCookie(cfg.GetAuthCookie()).
		Set("Content-type", "application/json").
		Set("User-Agent", UserAgentString(globalBuildVersion)).
		Timeout(10*time.Second).
		Retry(3, 1*time.Second, http.StatusBadGateway, http.StatusInternalServerError).
		SetCurlCommand(globalEnableCurl).
		SetDebug(globalEnableDebug)
}

func RenderTableToStdout(headers []string, data [][]string) {
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

func hostFromUri(str string) (error, string) {
	u, e := url.Parse(str)
	if e != nil {
		return e, ""
	}
	return e, u.Host
}

/*
 * return a java-style System.currentTimeMillis to match
 * the values being returned from the server.
 */
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func javaEpochToTime(long int64) time.Time {
	return time.Unix(0, long*int64(time.Millisecond))
}

func javaEpochToHumanizedTime(long int64) string {
	return humanize.Time(javaEpochToTime(long))
}

func JavaEpochToDateStr(long int64) string {
	t := javaEpochToTime(long)
	return t.Format(time.RFC3339)
}

func ProgressIndicator() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("green")
	return s
}

func PrintTerminalErrors(errs []error) {
	for i, j := 0, len(errs)-1; i < j; i, j = i+1, j-1 {
		errs[i], errs[j] = errs[j], errs[i]
	}

	for _, e := range errs {
		fmt.Println(e)
	}
}

func isValidGUID(in string) bool {
	match, _ := regexp.MatchString(`^[a-z0-9]{12,12}$`, in)
	return match
}

func isValidCommaDelimitedList(str string) bool {
	match, _ := regexp.MatchString(`^([a-z0-9/\\-]+,?)+$`, str)
	return match
}

func CurrentVersion() string {
	if len(globalBuildVersion) == 0 {
		return "dev"
	} else {
		return "v" + globalBuildVersion
	}
}

func UserAgentString(globalBuildVersion string) string {
	var name = "NelsonCLI"
	var version = getVersionForMode(globalBuildVersion)
	return name + "/" + version + " (" + runtime.GOOS + ")"
}

func getVersionForMode(globalBuildVersion string) string {
	if len(globalBuildVersion) == 0 {
		return "dev"
	} else {
		return "0.5." + globalBuildVersion
	}
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
