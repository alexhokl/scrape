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
