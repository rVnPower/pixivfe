package config

import (
	"testing"
)

func TestGetRandomUserAgent(t *testing.T) {
	// Test that the function returns a non-empty string
	ua := GetRandomUserAgent()
	if ua == "" {
		t.Error("GetRandomUserAgent returned an empty string")
	}

	// Test that the returned user agent is in the list
	found := false
	for _, builtinUA := range BuiltinUserAgentList {
		if ua == builtinUA {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetRandomUserAgent returned a user agent not in BuiltinUserAgentList")
	}

	// Test that BuiltinUserAgentList is not empty
	if len(BuiltinUserAgentList) == 0 {
		t.Error("BuiltinUserAgentList is empty")
	}
}
