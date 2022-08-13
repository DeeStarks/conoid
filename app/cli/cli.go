package cli

import (
	"database/sql"

	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/utils"
	_ "github.com/mattn/go-sqlite3"
)

type ICLICommands interface {
	Services() *ServiceCommand
}

type CLICommands struct {
	defaultDB *sql.DB
}

// Accept and process CLI commands.
func NewCLICommands() ICLICommands {
	// Connect to the default db
	defaultDB, err := sql.Open("sqlite3", config.DEFAULT_DB)
	if err != nil {
		utils.Log("Could not connect DB:", err)
	}
	return &CLICommands{
		defaultDB: defaultDB,
	}
}
