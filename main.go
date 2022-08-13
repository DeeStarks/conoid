package main

import (
	"fmt"
	"os"

	"github.com/DeeStarks/conoid/cmd"
	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/utils"
)

// Create and setup dependencies
func SetupDeps() error {
	// 1. Setup the data root directory
	if _, err := os.Stat(config.DATA_ROOT); os.IsNotExist(err) {
		err := os.Mkdir(config.DATA_ROOT, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 2. Create the default database
	if _, err := os.Stat(config.DEFAULT_DB); os.IsNotExist(err) {
		// Create if doesn't exist
		f, err := os.OpenFile(config.DEFAULT_DB, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error creating db file: %s; Error: %v", config.DEFAULT_DB, err)
		}
		f.Write([]byte("{}")) // Initialize with an empty object
		f.Close()

	}

	// 3. Log root
	if _, err := os.Stat(config.LOGS_ROOT); os.IsNotExist(err) {
		// Create root if doesn't exist
		err := os.Mkdir(config.LOGS_ROOT, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 4. Service logs
	if _, err := os.Stat(config.SERVICE_LOGS); os.IsNotExist(err) {
		// Create if services log file
		f, err := os.OpenFile(config.SERVICE_LOGS, os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error creating log file: %s; Error: %v", config.SERVICE_LOGS, err)
		}
		f.Close()
	}
	return nil
}

func main() {
	// Setup dependencies during installation
	err := SetupDeps()
	if err != nil {
		utils.Log(err)
		return
	}

	// Start process based on command argument
	cmd.Execute()
}
