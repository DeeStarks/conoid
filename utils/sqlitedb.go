package utils

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Migrate sqlite3 schema
func Sqlite3Migrate(pathToDB, pathToSchema string) error {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		return err
	}

	// Open the sql file
	schema, err := os.Open(pathToSchema)
	if err != nil {
		return err
	}

	// Read to buffer
	buf := make([]byte, 1024)
	n, err := schema.Read(buf)
	if err != nil {
		return err
	}

	// Execute
	script := string(buf[:n])
	_, err = db.Exec(script)
	if err != nil {
		return err
	}
	return nil
}
