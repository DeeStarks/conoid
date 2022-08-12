package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DeeStarks/conoid/cmd"
	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/domain/schemas"
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
		f, err := os.OpenFile(config.DEFAULT_DB, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return fmt.Errorf("error creating db file: %s; Error: %v", config.DEFAULT_DB, err)
		}
		f.Close()

		// Migrate db schema
		err = utils.Sqlite3ScriptMigrate(config.DEFAULT_DB, schemas.DefaultScript)
		if err != nil {
			return fmt.Errorf("error migrating schema: %v", err)
		}
	}
	return nil
}

func main() {
	// Setup dependencies during installation
	err := SetupDeps()
	if err != nil {
		log.Println(err)
		return
	}

	// Start process based on command argument
	cmd.Execute()
}
