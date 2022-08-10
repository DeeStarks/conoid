package main

import (
	"log"
	"os"

	"github.com/DeeStarks/conoid/cmd"
	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/domain/schemas"
	"github.com/DeeStarks/conoid/utils"
)

// Create and setup dependencies
func SetupDeps() {
	// 1. Setup the data root directory
	if _, err := os.Stat(config.FS_ROOT); os.IsNotExist(err) {
		err := os.Mkdir(config.FS_ROOT, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// 2. Create the default database
	if _, err := os.Stat(config.DEFAULT_DB); os.IsNotExist(err) {
		// Create if doesn't exist
		f, err := os.OpenFile(config.DEFAULT_DB, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Error creating db file: %s; Error: %v\n", config.DEFAULT_DB, err)
			return
		}
		f.Close()

		// Migrate db schema
		err = utils.Sqlite3ScriptMigrate(config.DEFAULT_DB, schemas.DefaultScript)
		if err != nil {
			log.Println("Error migrating schema:", err)
			return
		}
	}
}

func main() {
	// Setup dependencies during installation
	SetupDeps()

	// Start process based on command argument
	cmd.Execute()
}
