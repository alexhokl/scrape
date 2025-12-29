package scraper

import "testing"

func TestRemoveExtraSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no extra spaces",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "multiple spaces between words",
			input:    "hello    world",
			expected: "hello world",
		},
		{
			name:     "leading spaces",
			input:    "   hello",
			expected: " hello",
		},
		{
			name:     "trailing spaces",
			input:    "hello   ",
			expected: "hello ",
		},
		{
			name:     "newlines are removed",
			input:    "hello\nworld",
			expected: "helloworld",
		},
		{
			name:     "carriage returns are removed",
			input:    "hello\rworld",
			expected: "helloworld",
		},
		{
			name:     "mixed newlines and spaces",
			input:    "hello  \n  world",
			expected: "hello world",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: " ",
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: "",
		},
		{
			name:     "tabs are preserved",
			input:    "hello\tworld",
			expected: "hello\tworld",
		},
		{
			name:     "unicode characters",
			input:    "héllo  wörld",
			expected: "héllo wörld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeExtraSpaces(tt.input)
			if result != tt.expected {
				t.Errorf("removeExtraSpaces(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTrimSpacesAndLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no trimming needed",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "leading spaces",
			input:    "   hello",
			expected: "hello",
		},
		{
			name:     "trailing spaces",
			input:    "hello   ",
			expected: "hello",
		},
		{
			name:     "leading and trailing spaces",
			input:    "   hello   ",
			expected: "hello",
		},
		{
			name:     "leading newlines",
			input:    "\n\nhello",
			expected: "hello",
		},
		{
			name:     "trailing newlines",
			input:    "hello\n\n",
			expected: "hello",
		},
		{
			name:     "leading and trailing newlines",
			input:    "\n\nhello\n\n",
			expected: "hello",
		},
		{
			name:     "mixed spaces and newlines at ends",
			input:    "  \n hello \n  ",
			expected: "hello",
		},
		{
			name:     "internal spaces preserved",
			input:    "  hello   world  ",
			expected: "hello   world",
		},
		{
			name:     "internal newlines preserved",
			input:    "\nhello\nworld\n",
			expected: "hello\nworld",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: "",
		},
		{
			name:     "space newline space pattern",
			input:    " \n ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimSpacesAndLineBreaks(tt.input)
			if result != tt.expected {
				t.Errorf("trimSpacesAndLineBreaks(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateFileNameFromTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple title",
			input:    "Hello World",
			expected: "hello_world",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello_world",
		},
		{
			name:     "all uppercase",
			input:    "HELLO WORLD",
			expected: "hello_world",
		},
		{
			name:     "single word",
			input:    "Introduction",
			expected: "introduction",
		},
		{
			name:     "multiple spaces become single underscore",
			input:    "Hello  World",
			expected: "hello_world",
		},
		{
			name:     "three spaces become single underscore",
			input:    "Hello   World",
			expected: "hello__world",
		},
		{
			name:     "preserves numbers",
			input:    "Version 2.0 Release",
			expected: "version_2.0_release",
		},
		{
			name:     "preserves special characters",
			input:    "Go: The Language",
			expected: "go:_the_language",
		},
		{
			name:     "preserves unicode characters",
			input:    "Go 言語入門",
			expected: "go_言語入門",
		},
		{
			name:     "leading space",
			input:    " Hello World",
			expected: "_hello_world",
		},
		{
			name:     "trailing space",
			input:    "Hello World ",
			expected: "hello_world_",
		},
		{
			name:     "mixed case with numbers",
			input:    "Azure SDK 2.0 Release Notes",
			expected: "azure_sdk_2.0_release_notes",
		},
		{
			name:     "hyphens preserved",
			input:    "Step-by-Step Guide",
			expected: "step-by-step_guide",
		},
		{
			name:     "underscores preserved",
			input:    "file_name example",
			expected: "file_name_example",
		},
		{
			name:     "dots preserved",
			input:    "v1.2.3 Release",
			expected: "v1.2.3_release",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFileNameFromTitle(tt.input)
			if result != tt.expected {
				t.Errorf("generateFileNameFromTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
