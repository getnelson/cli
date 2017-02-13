package main

import (
	"strings"
	"testing"
)

func TestUserAgentString(t *testing.T) {
	// Test for development, where global build version doesn't exist.
	var globalBuildVersion1 = ""
	var expectedDevelUserAgentString1 = "NelsonCLI/dev"
	var result1 = UserAgentString(globalBuildVersion1)
	if !strings.Contains(result1, expectedDevelUserAgentString1+" (") {
		t.Error("devel user agent string is incorrect: \n" + result1 + "\n" + expectedDevelUserAgentString1)
	}

	// Test for released version, where global build version exists.
	var globalBuildVersion2 = "123"
	var expectedDevelUserAgentString2 = "NelsonCLI/0.2." + globalBuildVersion2
	var result2 = UserAgentString(globalBuildVersion2)
	if !strings.Contains(result2, expectedDevelUserAgentString2+" (") {
		t.Error("devel user agent string is incorrect: \n" + result2 + "\n" + expectedDevelUserAgentString2)
	}

}
