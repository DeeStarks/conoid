package cli

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DeeStarks/conoid/domain/ports"
	"github.com/DeeStarks/conoid/utils"
	"github.com/aquasecurity/table"
	"github.com/google/uuid"
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
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.AppProcesses().RetrieveRunning()
	if err != nil {
		fmt.Println("Error retrieve running apps:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetHeaders("PID", "NAME", "TYPE", "LISTENERS", "ROOT", "CLIENT", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())
		listeners := strings.Join(p.Listeners, ", ")

		// Handle booleans
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
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.AppProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Error retrieving apps:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetHeaders("PID", "NAME", "STATUS", "TYPE", "LISTENERS", "ROOT", "CLIENT", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())
		listeners := strings.Join(p.Listeners, ", ")

		// Handle booleans
		tunnelled := "False"
		if p.Tunnelled {
			tunnelled = "True"
		}
		status := "Stopped"
		if p.Status {
			status = "Running"
		}

		t.AddRow(
			string(p.Pid), p.Name, status, p.Type, listeners, p.RootDirectory,
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
		fmt.Println("Could not add new app; \"conoid.yml\": File not found")
		return
	}
	defer f.Close()

	conf, err := utils.DeserializeAppYAML(f)
	if err != nil {
		fmt.Println("Error deserializing configuration file:", err)
		return
	}

	// Vallidate configuration
	validatedConf, err := utils.ValidateConf(conf)
	if err != nil {
		fmt.Println("Invalid configuration file:", err)
		return
	}

	// Ensure app doesn't already exist
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.AppProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Could not add app:", err)
		return
	}

	// Check processes' names
	for _, p := range processes {
		if p.Name == validatedConf.Name {
			fmt.Printf("Could not add app: an app with name \"%s\" already exists\n", validatedConf.Name)
			return
		}
	}

	// Convert the configurations to a map type
	// First, to json
	jsondata, err := json.Marshal(validatedConf)
	if err != nil {
		fmt.Println("Could not add app:", err)
		return
	}
	// Now, to map
	var mapConf map[string]interface{}
	json.Unmarshal(jsondata, &mapConf)

	// Change the "tunnelled" field from boolean type to int
	if mapConf["tunnelled"].(bool) {
		mapConf["tunnelled"] = 1
	} else {
		mapConf["tunnelled"] = 0
	}

	// Join the "listeners" slice and save in the database as string type
	if s, ok := mapConf["listeners"].([]interface{}); ok {
		ss := make([]string, len(s))
		for i, v := range s {
			ss[i] = fmt.Sprintf("%v", v)
		}
		mapConf["listeners"] = strings.Join(ss, ", ")
	}

	// Generate pid, status, and created at
	mapConf["pid"] = strings.ReplaceAll(uuid.New().String(), "-", "") // Stripping hyphens
	mapConf["status"] = 1
	mapConf["created_at"] = time.Now().Unix()
	
	_, err = domainPort.AppProcesses().Create(mapConf)
	if err != nil {
		fmt.Println("Could not add app:", err)
		return
	}
	fmt.Println("Your app was successfully added. Restart conoid to start accepting request")
}
