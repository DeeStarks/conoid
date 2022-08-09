package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/DeeStarks/conoid/utils"
)

type AppProcessModel struct {
	Pid           string
	Name          string
	Status        bool
	Type          string
	Listeners     []string
	RootDirectory string
	ClientAddress string
	Tunnelled     bool
	CreatedAt     int64
}

type AppProcess struct {
	DB *sql.DB
}

// Retrive all running services
func (p AppProcess) RetrieveRunning() ([]AppProcessModel, error) {
	rows, err := p.DB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes WHERE status=1
	`)
	if err != nil {
		return nil, err
	}

	// Parse result
	var processes []AppProcessModel
	for rows.Next() {
		var process AppProcessModel

		// Handle null values
		var listeners, root_directory, client_address sql.NullString

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &root_directory, &client_address,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		process.RootDirectory = root_directory.String
		process.ClientAddress = client_address.String
		// Listeners are stored in the db as strings separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners.String, ", ")

		// Append the process to list of processes
		processes = append(processes, process)
	}
	return processes, nil
}

// Retrive all services
func (p AppProcess) RetrieveAll() ([]AppProcessModel, error) {
	rows, err := p.DB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes
	`)
	if err != nil {
		return nil, err
	}

	// Parse result
	var processes []AppProcessModel
	for rows.Next() {
		var process AppProcessModel

		// Handle null values
		var listeners, root_directory, client_address sql.NullString

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &root_directory, &client_address,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		process.RootDirectory = root_directory.String
		process.ClientAddress = client_address.String
		// Listeners are stored in the db as strings separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners.String, ", ")

		// Append the process to list of processes
		processes = append(processes, process)
	}
	return processes, nil
}

func (p AppProcess) Create(data map[string]interface{}) (AppProcessModel, error) {
	var process AppProcessModel

	// Get keys and values
	refKeys := reflect.ValueOf(data).MapKeys()
	keys := make([]string, len(refKeys))
	values := make([]interface{}, len(refKeys))
	for i, k := range refKeys {
		keys[i] = k.String()
		values[i] = data[k.String()]
	}

	// Execute query
	query := fmt.Sprintf(`
		INSERT INTO processes ( %s ) VALUES ( %s )
		RETURNING pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	`, strings.Join(keys, ", "), utils.GeneratePlaceholders(len(keys)))

	// Handle null values
	var listeners, root_directory, client_address sql.NullString
	err := p.DB.QueryRow(query, values...).Scan(
		&process.Pid, &process.Name, &process.Status, &process.Type,
		&listeners, &root_directory, &client_address,
		&process.Tunnelled, &process.CreatedAt,
	)
	if err != nil {
		return AppProcessModel{}, err
	}

	process.RootDirectory = root_directory.String
	process.ClientAddress = client_address.String
	// Listeners are stored in the db as strings separated by comma
	// we'll split that into slice
	process.Listeners = strings.Split(listeners.String, ", ")

	// Return result
	return process, nil
}
