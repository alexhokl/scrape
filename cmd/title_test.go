package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateTitleOptions_ValidGuardianSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid guardian source, got: %v", err)
	}
}

func TestValidateTitleOptions_ValidMicrosoftSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "microsoft",
		url:    "https://learn.microsoft.com/some-article",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid microsoft source, got: %v", err)
	}
}

func TestValidateTitleOptions_ValidGoSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "go",
		url:    "https://go.dev/doc/some-article",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid go source, got: %v", err)
	}
}

func TestValidateTitleOptions_ValidTofuguSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "tofugu",
		url:    "https://www.tofugu.com/some-article",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid tofugu source, got: %v", err)
	}
}

func TestValidateTitleOptions_InvalidSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "unsupported",
		url:    "https://example.com",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

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

func TestValidateTitleOptions_EmptySource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "",
		url:    "https://example.com",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestValidateTitleOptions_EmptyURL(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "guardian",
		url:    "",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}

	if !strings.Contains(err.Error(), "url is required") {
		t.Errorf("expected error message to contain 'url is required', got: %v", err)
	}
}

func TestValidateTitleOptions_CaseSensitiveSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	testCases := []string{"Guardian", "GUARDIAN", "Microsoft", "MICROSOFT", "Go", "GO", "Tofugu", "TOFUGU"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			titleOpts = titleOptions{
				source: tc,
				url:    "https://example.com",
			}

			err := validateTitleOptions(&cobra.Command{}, []string{})

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if !strings.Contains(err.Error(), "invalid source") {
				t.Errorf("expected error message to contain 'invalid source', got: %v", err)
			}
		})
	}
}

func TestValidateTitleOptions_ValidationOrder(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	// When both source is invalid and url is empty, source should be validated first
	titleOpts = titleOptions{
		source: "invalid-source",
		url:    "",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when both source is invalid and url is empty, got nil")
	}

	// Should fail on source validation first
	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source' (validation order), got: %v", err)
	}
}

func TestValidateTitleOptions_AllFieldsEmpty(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := titleOpts
	defer func() { titleOpts = originalOpts }()

	titleOpts = titleOptions{
		source: "",
		url:    "",
	}

	err := validateTitleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when all fields are empty, got nil")
	}

	// Should fail on source validation first
	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestTitleOptions_StructFields(t *testing.T) {
	opts := titleOptions{
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
