package utils_test

import (
	"testing"

	"github.com/DeeStarks/conoid/utils"
)

func TestGeneratePlaceholders(t *testing.T) {
	tests := []struct{
		n int
		expected string
	} {
		{n: 4, expected: "$1, $2, $3, $4"},
		{n: 0, expected: ""},
		{n: 1, expected: "$1"},
	}

	for _, tc := range tests {
		res := utils.GeneratePlaceholders(tc.n)
		if res != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, res)
		}
	}
}