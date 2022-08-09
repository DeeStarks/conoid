package utils

import (
	"fmt"
)

var (
	SECONDS int64 = 1
	MINUTES int64 = 60
	HOURS         = MINUTES * 60
	DAYS          = HOURS * 24
	WEEKS         = DAYS * 7
	MONTHS        = WEEKS * 4
	YEARS         = MONTHS * 12
)

// Expects unix timestamp as parameter values
func TimeAgo(past, present int64) string {
	diff := present - past

	if diff/YEARS >= 1 {
		return fmt.Sprintf("%d year(s) ago.", diff/YEARS)
	} else if diff/MONTHS >= 1 {
		return fmt.Sprintf("%d month(s) ago.", diff/MONTHS)
	} else if diff/WEEKS >= 1 {
		return fmt.Sprintf("%d week(s) ago.", diff/WEEKS)
	} else if diff/DAYS >= 1 {
		return fmt.Sprintf("%d day(s) ago.", diff/DAYS)
	} else if diff/HOURS >= 1 {
		return fmt.Sprintf("%d hour(s) ago.", diff/HOURS)
	} else if diff/MINUTES >= 1 {
		return fmt.Sprintf("%d minute(s) ago.", diff/MINUTES)
	} else {
		return fmt.Sprintf("%d second(s) ago.", diff/SECONDS)
	}
}
