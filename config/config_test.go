package config

import (
	"testing"
)

// TestParseRevision is a test function that verifies the behavior of the parseRevision function.
func TestParseRevision(t *testing.T) {
	tests := []struct {
		name          string // Description of the test case
		revision      string // Input revision string
		expectedDate  string // Expected date output
		expectedHash  string // Expected hash output
		expectedDirty bool   // Expected isDirty output
	}{
		{"Valid revision", "2024.09.24-18d6874", "2024.09.24", "18d6874", false},
		{"Dirty revision", "2024.09.24-18d6874+dirty", "2024.09.24", "18d6874", true},
		{"Empty revision", "", "unknown", "unknown", false},
		{"Only hash", "18d6874", "unknown", "18d6874", false},
		{"Only hash dirty", "18d6874+dirty", "unknown", "18d6874", true},
		{"Invalid format", "2024.09.24-18d6874-extra", "unknown", "2024.09.24-18d6874-extra", false},
		{"Invalid format dirty", "2024.09.24-18d6874-extra+dirty", "unknown", "2024.09.24-18d6874-extra", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, gotHash, gotDirty := parseRevision(tt.revision)

			if gotDate != tt.expectedDate {
				t.Errorf("parseRevision() gotDate = %v, want %v", gotDate, tt.expectedDate)
			}
			if gotHash != tt.expectedHash {
				t.Errorf("parseRevision() gotHash = %v, want %v", gotHash, tt.expectedHash)
			}
			if gotDirty != tt.expectedDirty {
				t.Errorf("parseRevision() gotDirty = %v, want %v", gotDirty, tt.expectedDirty)
			}
		})
	}
}

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
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateURL(tt.urlStr, tt.urlType)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.String() != tt.expected {
					t.Errorf("validateURL() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
