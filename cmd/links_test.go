package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateLinksOptions_ValidGuardianSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	linksOpts = linksOptions{
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateLinksOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid guardian source, got: %v", err)
	}
}

func TestValidateLinksOptions_InvalidSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	linksOpts = linksOptions{
		source: "unsupported",
		url:    "https://example.com",
	}

	err := validateLinksOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for invalid source, got nil")
	}

	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}

	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected error message to contain the invalid source name, got: %v", err)
	}
}

func TestValidateLinksOptions_EmptySource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	linksOpts = linksOptions{
		source: "",
		url:    "https://example.com",
	}

	err := validateLinksOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestValidateLinksOptions_EmptyURL(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	linksOpts = linksOptions{
		source: "guardian",
		url:    "",
	}

	err := validateLinksOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}

	if !strings.Contains(err.Error(), "url is required") {
		t.Errorf("expected error message to contain 'url is required', got: %v", err)
	}
}

func TestValidateLinksOptions_CaseSensitiveSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	testCases := []string{"Guardian", "GUARDIAN", "GuArDiAn"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			linksOpts = linksOptions{
				source: tc,
				url:    "https://example.com",
			}

			err := validateLinksOptions(&cobra.Command{}, []string{})

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if !strings.Contains(err.Error(), "invalid source") {
				t.Errorf("expected error message to contain 'invalid source', got: %v", err)
			}
		})
	}
}

func TestValidateLinksOptions_BothSourceAndURLEmpty(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := linksOpts
	defer func() { linksOpts = originalOpts }()

	linksOpts = linksOptions{
		source: "",
		url:    "",
	}

	err := validateLinksOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when both source and URL are empty, got nil")
	}

	// Should fail on source validation first
	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestLinksOptions_StructFields(t *testing.T) {
	opts := linksOptions{
		source: "guardian",
		url:    "https://www.theguardian.com/test",
	}

	if opts.source != "guardian" {
		t.Errorf("expected source to be 'guardian', got %q", opts.source)
	}

	if opts.url != "https://www.theguardian.com/test" {
		t.Errorf("expected url to be 'https://www.theguardian.com/test', got %q", opts.url)
	}
}
