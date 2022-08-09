package utils_test

import (
	"testing"
	"time"

	"github.com/DeeStarks/conoid/utils"
)

func TestTimeAgo(t *testing.T) {
	tests := []struct {
		past, present int64
		expected      string
	}{
		{
			past:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			expected: "1 year(s) ago.",
		},
		{
			past:     time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			expected: "2 year(s) ago.",
		},
		{
			past:     time.Date(2010, time.October, 10, 23, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			expected: "1 month(s) ago.",
		},
		{
			past:     time.Date(2010, time.November, 8, 23, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
			expected: "2 day(s) ago.",
		},
		{
			past:     time.Date(2010, time.November, 10, 10, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 15, 0, 0, 0, time.UTC).Unix(),
			expected: "5 hour(s) ago.",
		},
		{
			past:     time.Date(2010, time.November, 10, 15, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 15, 4, 0, 0, time.UTC).Unix(),
			expected: "4 minute(s) ago.",
		},
		{
			past:     time.Date(2010, time.November, 10, 15, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 15, 0, 15, 0, time.UTC).Unix(),
			expected: "15 second(s) ago.",
		},
		{
			past:     time.Date(2010, time.November, 10, 15, 0, 0, 0, time.UTC).Unix(),
			present:  time.Date(2010, time.November, 10, 15, 0, 0, 10, time.UTC).Unix(),
			expected: "0 second(s) ago.",
		},
	}

	for _, tc := range tests {
		timeago := utils.TimeAgo(tc.past, tc.present)
		if timeago != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, timeago)
		}
	}
}
