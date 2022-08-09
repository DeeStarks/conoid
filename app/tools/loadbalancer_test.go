package tools_test

import (
	"testing"

	"github.com/DeeStarks/conoid/app/tools"
)

func TestLoadBalancer(t *testing.T) {
	var servers = []string{"google.com", "localhost", "stackoverflow.com", "yahoo.com"}
	var tests = []struct {
		expected string
	}{
		{
			expected: "google.com",
		},
		{
			expected: "localhost",
		},
		{
			expected: "stackoverflow.com",
		},
		{
			expected: "yahoo.com",
		},
		{
			expected: "google.com",
		},
		{
			expected: "localhost",
		},
	}

	lb := tools.NewLoadBalancer(servers)
	for _, tc := range tests {
		server := lb.GetNextServer()
		if server != tc.expected {
			t.Errorf("Expected %s; Got %s", tc.expected, server)
		}
	}
}
