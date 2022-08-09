package repository

import (
	"database/sql"
	"log"
	"strings"
)

type AppProcessModel struct {
	Pid           string
	Name          string
	Status        string
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
func (p AppProcess) RetrieveRunning() []AppProcessModel {
	rows, err := p.DB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes WHERE status="running"
	`)
	if err != nil {
		log.Println("Error retrieving running apps:", err)
		return nil
	}

	// Parse result
	var processes []AppProcessModel
	for rows.Next() {
		var process AppProcessModel
		var listeners string

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &process.RootDirectory, &process.ClientAddress,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			log.Println("Error retrieving running services:", err)
			return nil
		}

		// Listeners are stored in the db as strings separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners, ",")
		// Append the process to list of processes
		processes = append(processes, process)
	}
	return processes
}

// Retrive all services
func (p AppProcess) RetrieveAll() []AppProcessModel {
	rows, err := p.DB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes
	`)
	if err != nil {
		log.Println("Error retrieving services:", err)
		return nil
	}

	// Parse result
	var processes []AppProcessModel
	for rows.Next() {
		var process AppProcessModel
		var listeners string

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &process.RootDirectory, &process.ClientAddress,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			log.Println("Error retrieving running apps:", err)
			return nil
		}

		// Listeners are stored in the db as strings separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners, ",")
		// Append the process to list of processes
		processes = append(processes, process)
	}
	return processes
}
