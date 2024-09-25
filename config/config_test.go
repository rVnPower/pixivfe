package config

import (
	"testing"
)

// TestValidateURL is a test function that verifies the behavior of the validateURL function.
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string // Description of the test case
		urlStr   string // Input URL string to be validated
		urlType  string // Type of URL (used for error messaging)
		wantErr  bool   // Whether an error is expected
		expected string // Expected output URL string
	}{
		{"Valid URL", "https://example.com", "Test", false, "https://example.com"},
		{"Valid URL with path", "https://example.com/path", "Test", false, "https://example.com/path"},
		{"Missing scheme", "example.com", "Test", true, ""},
		{"Missing host", "https://", "Test", true, ""},
		{"Trailing slash", "https://example.com/", "Test", true, ""},
		{"Empty URL", "", "Test", true, ""},
		{"URL with query params", "https://example.com/path?q=test", "Test", false, "https://example.com/path?q=test"},
		{"URL with fragment", "https://example.com/path#fragment", "Test", false, "https://example.com/path#fragment"},
	}

	for _, tt := range tests {
		// Run each test case as a subtest
		t.Run(tt.name, func(t *testing.T) {
			// Call the validateURL function with test input
			got, err := validateURL(tt.urlStr, tt.urlType)

			// Check if the error result matches the expected error state
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error is expected, compare the output URL string with the expected result
			if !tt.wantErr {
				if got.String() != tt.expected {
					t.Errorf("validateURL() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
