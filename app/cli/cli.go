package cli

import (
	"database/sql"
	"log"

	"github.com/DeeStarks/conoid/config"
	_ "github.com/mattn/go-sqlite3"
)

type ICLICommands interface {
	Apps() *AppCommand
}

type CLICommands struct {
	defaultDB *sql.DB
}

// Accept and process CLI commands.
func NewCLICommands() ICLICommands {
	// Connect to the default db
	defaultDB, err := sql.Open("sqlite3", config.DEFAULT_DB)
	if err != nil {
		log.Panicln("Could not connect DB:", err)
	}
	return &CLICommands{
		defaultDB: defaultDB,
	}
}
