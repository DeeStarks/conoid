package utils_test

import (
	"testing"

	"github.com/DeeStarks/conoid/utils"
)

func TestTruncateString(t *testing.T) {
	tests := []struct{
		s	string
		l 	int
		expected	string
	} {
		{
			s: "Hello, World!",
			l: 4,
			expected: "Hell...",
		},
		{
			s: "Hello, World!",
			l: 0,
			expected: "...",
		},
		{
			s: "Hello, World!",
			l: 12,
			expected: "Hello, World...",
		},
		{
			s: "Hello, World!",
			l: 13,
			expected: "Hello, World!",
		},
		{
			s: "Hello, World!",
			l: 14,
			expected: "Hello, World!",
		},
	}

	for _, tc := range tests {
		res := utils.TruncateString(tc.s, tc.l)
		if res != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, res)
		}
	}
}