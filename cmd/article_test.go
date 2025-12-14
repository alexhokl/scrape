package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateArticleOptions_ValidGuardianSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid guardian source, got: %v", err)
	}
}

func TestValidateArticleOptions_ValidMicrosoftSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "microsoft",
		url:    "https://learn.microsoft.com/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid microsoft source, got: %v", err)
	}
}

func TestValidateArticleOptions_ValidGoSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "go",
		url:    "https://go.dev/doc/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid go source, got: %v", err)
	}
}

func TestValidateArticleOptions_ValidTofuguSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "tofugu",
		url:    "https://www.tofugu.com/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err != nil {
		t.Errorf("expected no error for valid tofugu source, got: %v", err)
	}
}

func TestValidateArticleOptions_InvalidSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "unsupported",
		url:    "https://example.com",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

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

func TestValidateArticleOptions_EmptySource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "",
		url:    "https://example.com",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty source, got nil")
	}

	if !strings.Contains(err.Error(), "invalid source") {
		t.Errorf("expected error message to contain 'invalid source', got: %v", err)
	}
}

func TestValidateArticleOptions_InvalidFormat(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "json",
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}

	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected error message to contain 'invalid format', got: %v", err)
	}

	if !strings.Contains(err.Error(), "json") {
		t.Errorf("expected error message to contain the invalid format name, got: %v", err)
	}
}

func TestValidateArticleOptions_EmptyFormat(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "",
		source: "guardian",
		url:    "https://www.theguardian.com/some-article",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty format, got nil")
	}

	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected error message to contain 'invalid format', got: %v", err)
	}
}

func TestValidateArticleOptions_EmptyURL(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "markdown",
		source: "guardian",
		url:    "",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}

	if !strings.Contains(err.Error(), "url is required") {
		t.Errorf("expected error message to contain 'url is required', got: %v", err)
	}
}

func TestValidateArticleOptions_CaseSensitiveSource(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	testCases := []string{"Guardian", "GUARDIAN", "Microsoft", "MICROSOFT", "Go", "GO", "Tofugu", "TOFUGU"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			articleOpts = articleOptions{
				format: "markdown",
				source: tc,
				url:    "https://example.com",
			}

			err := validateArticleOptions(&cobra.Command{}, []string{})

			if err == nil {
				t.Errorf("expected error for case-sensitive source %q, got nil", tc)
			}

			if !strings.Contains(err.Error(), "invalid source") {
				t.Errorf("expected error message to contain 'invalid source', got: %v", err)
			}
		})
	}
}

func TestValidateArticleOptions_CaseSensitiveFormat(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	testCases := []string{"Markdown", "MARKDOWN", "MarkDown"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			articleOpts = articleOptions{
				format: tc,
				source: "guardian",
				url:    "https://example.com",
			}

			err := validateArticleOptions(&cobra.Command{}, []string{})

			if err == nil {
				t.Errorf("expected error for case-sensitive format %q, got nil", tc)
			}

			if !strings.Contains(err.Error(), "invalid format") {
				t.Errorf("expected error message to contain 'invalid format', got: %v", err)
			}
		})
	}
}

func TestValidateArticleOptions_ValidationOrder(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	// When both format and source are invalid, format should be validated first
	articleOpts = articleOptions{
		format: "invalid-format",
		source: "invalid-source",
		url:    "https://example.com",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when both format and source are invalid, got nil")
	}

	// Should fail on format validation first
	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected error message to contain 'invalid format' (validation order), got: %v", err)
	}
}

func TestValidateArticleOptions_AllFieldsEmpty(t *testing.T) {
	// Save original opts and restore after test
	originalOpts := articleOpts
	defer func() { articleOpts = originalOpts }()

	articleOpts = articleOptions{
		format: "",
		source: "",
		url:    "",
	}

	err := validateArticleOptions(&cobra.Command{}, []string{})

	if err == nil {
		t.Fatal("expected error when all fields are empty, got nil")
	}

	// Should fail on format validation first
	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected error message to contain 'invalid format', got: %v", err)
	}
}

func TestArticleOptions_StructFields(t *testing.T) {
	opts := articleOptions{
		format: "markdown",
		source: "guardian",
		url:    "https://www.theguardian.com/test",
	}

	if opts.format != "markdown" {
		t.Errorf("expected format to be 'markdown', got %q", opts.format)
	}

	if opts.source != "guardian" {
		t.Errorf("expected source to be 'guardian', got %q", opts.source)
	}

	if opts.url != "https://www.theguardian.com/test" {
		t.Errorf("expected url to be 'https://www.theguardian.com/test', got %q", opts.url)
	}
}
