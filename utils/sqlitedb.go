package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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

// -> $1, $2...$n
func GeneratePlaceholders(n int) string {
	holders := make([]string, n)
	for i := 0; i < n; i++ {
		holders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(holders, ", ")
}

// -> id = $1...
func GenerateSetConditions(keys []string) string {
	conditions := make([]string, len(keys))
	for i, k := range keys {
		conditions[i] = fmt.Sprintf("%s = $%d", k, i+1)
	}
	return strings.Join(conditions, ", ")
}
