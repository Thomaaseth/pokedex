package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "learninggo",
			expected: []string{"learninggo"},
		},
		{
			input:    "tomorrow is gonna be sunny",
			expected: []string{"tomorrow", "is", "gonna", "be", "sunny"},
		},
		{
			input:    "testing thecleaning up",
			expected: []string{"testing", "thecleaning", "up"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) returned word %q at position %d, expected %q", c.input, actual[i], i, c.expected[i])
			}

		}

	}
}
