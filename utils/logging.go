package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/DeeStarks/conoid/config"
)

func Log(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log := fmt.Sprintf("[%v]: %s\r\n", time.Now().Format("2006-01-02 15:04:05"), msg)

	// Add to log file
	f, err := os.OpenFile(config.SERVICE_LOGS, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(log); err != nil {
		panic(err)
	}
}

func Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log := fmt.Sprintf("[%v]: %s\r\n", time.Now().Format("2006-01-02 15:04:05"), msg)

	// Add to log file
	f, err := os.OpenFile(config.SERVICE_LOGS, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(log); err != nil {
		panic(err)
	}
}
