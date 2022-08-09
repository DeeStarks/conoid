package cli

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DeeStarks/conoid/utils"
	"github.com/aquasecurity/table"
	_ "github.com/mattn/go-sqlite3"
)

type AppCommand struct {
	defaultDB *sql.DB
}

func (cmd *CLICommands) Apps() *AppCommand {
	return &AppCommand{
		defaultDB: cmd.defaultDB,
	}
}

type AppProcess struct {
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

// List running applications
func (ac *AppCommand) ListRunning() {
	rows, err := ac.defaultDB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes WHERE status="running"
	`)
	if err != nil {
		log.Println("Error retrieving running apps:", err)
		return
	}

	// Parse result
	var processes []AppProcess
	for rows.Next() {
		var process AppProcess
		var listeners string

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &process.RootDirectory, &process.ClientAddress,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			log.Println("Error retrieving running apps:", err)
			return
		}

		// Listeners are stored in the db a string separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners, ",")
		// Append the process to list of processes
		processes = append(processes, process)
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetHeaders("PID", "NAME", "TYPE", "LISTENERS", "ROOT", "CLIENT", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())
		listeners := strings.Join(p.Listeners, ", ")
		tunnelled := "False"
		if p.Tunnelled {
			tunnelled = "True"
		}

		t.AddRow(
			string(p.Pid), p.Name, p.Type, listeners, p.RootDirectory,
			p.ClientAddress, tunnelled, created_at,
		)
	}
	t.Render()

}

// List all applications
func (ac *AppCommand) ListAll() {
	rows, err := ac.defaultDB.Query(`
	SELECT 
		pid, name, status, type, listeners, 
		root_directory, client_address, tunnelled, created_at 
	FROM processes
	`)
	if err != nil {
		log.Println("Error retrieving running apps:", err)
		return
	}

	// Parse result
	var processes []AppProcess
	for rows.Next() {
		var process AppProcess
		var listeners string

		err = rows.Scan(
			&process.Pid, &process.Name, &process.Status, &process.Type,
			&listeners, &process.RootDirectory, &process.ClientAddress,
			&process.Tunnelled, &process.CreatedAt,
		)
		if err != nil {
			log.Println("Error retrieving running apps:", err)
			return
		}

		// Listeners are stored in the db a string separated by comma
		// we'll split that into slice
		process.Listeners = strings.Split(listeners, ",")
		// Append the process to list of processes
		processes = append(processes, process)
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetHeaders("PID", "NAME", "STATUS", "TYPE", "LISTENERS", "ROOT", "CLIENT", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())
		listeners := strings.Join(p.Listeners, ", ")
		tunnelled := "False"
		if p.Tunnelled {
			tunnelled = "True"
		}

		t.AddRow(
			string(p.Pid), p.Name, p.Status, p.Type, listeners, p.RootDirectory,
			p.ClientAddress, tunnelled, created_at,
		)
	}
	t.Render()
}

// Add new application
func (ac *AppCommand) Add(filepath string) {
	// Ensure the file exists
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("Could not add new app; \"conoid.yml\": File not found")
		return
	}
	defer f.Close()

	conf, err := utils.DeserializeAppYAML(f)
	if err != nil {
		log.Println("Error deserializing configuration file:", err)
	}
	log.Println(conf)
}
