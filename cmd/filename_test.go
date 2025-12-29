package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateFilenameOptions_ValidGuardianSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid guardian source, got: %v", err)
	}
}

func TestValidateFilenameOptions_ValidMicrosoftSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "microsoft",
		url:    "https://learn.microsoft.com/some-article",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid microsoft source, got: %v", err)
	}
}

func TestValidateFilenameOptions_ValidGoSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "go",
		url:    "https://go.dev/doc/some-article",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid go source, got: %v", err)
	}
}

func TestValidateFilenameOptions_ValidTofuguSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "tofugu",
		url:    "https://www.tofugu.com/some-article",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid tofugu source, got: %v", err)
	}
}

func TestValidateFilenameOptions_InvalidSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "unsupported",
		url:    "https://example.com",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

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

func TestValidateFilenameOptions_EmptySource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "",
		url:    "https://example.com",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestValidateFilenameOptions_EmptyURL(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "guardian",
		url:    "",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}

	if !strings.Contains(err.Error(), "url is required") {
		t.Errorf("expected error message to contain 'url is required', got: %v", err)
	}
}

func TestValidateFilenameOptions_CaseSensitiveSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	testCases := []string{"Guardian", "GUARDIAN", "Microsoft", "MICROSOFT", "Go", "GO", "Tofugu", "TOFUGU"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			filenameOpts = filenameOptions{
				source: tc,
				url:    "https://example.com",
			}

			err := validateFilenameOptions(&cobra.Command{}, []string{})

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if !strings.Contains(err.Error(), "invalid source") {
				t.Errorf("expected error message to contain 'invalid source', got: %v", err)
			}
		})
	}
}

func TestValidateFilenameOptions_ValidationOrder(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	// When both source is invalid and url is empty, source should be validated first
	filenameOpts = filenameOptions{
		source: "invalid-source",
		url:    "",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when both source is invalid and url is empty, got nil")
	}

	// Should fail on source validation first
	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source' (validation order), got: %v", err)
	}
}

func TestValidateFilenameOptions_AllFieldsEmpty(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := filenameOpts
	defer func() { filenameOpts = originalOpts }()

	filenameOpts = filenameOptions{
		source: "",
		url:    "",
	}

	err := validateFilenameOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when all fields are empty, got nil")
	}

	// Should fail on source validation first
	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestFilenameOptions_StructFields(t *testing.T) {
	opts := filenameOptions{
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
