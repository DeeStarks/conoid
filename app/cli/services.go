package cli

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	port "github.com/DeeStarks/conoid/domain/ports"
	"github.com/DeeStarks/conoid/utils"
	"github.com/aquasecurity/table"
	"github.com/google/uuid"
)

type ServiceCommand struct {
	defaultDB *sql.DB
}

func (cmd *CLICommands) Services() *ServiceCommand {
	return &ServiceCommand{
		defaultDB: cmd.defaultDB,
	}
}

type ServiceProcess struct {
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

// List running services
func (ac *ServiceCommand) ListRunning() {
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveRunning()
	if err != nil {
		fmt.Println("Error retrieve running services:", err)
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

// List all services
func (ac *ServiceCommand) ListAll() {
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Error retrieving services:", err)
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

// Add new service
func (ac *ServiceCommand) Add(filepath string) {
	// Ensure the file exists
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Could not add new service; \"conoid.yml\": File not found")
		return
	}
	defer f.Close()

	conf, err := utils.DeserializeConf(f)
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

	// Ensure service doesn't already exist
	domainPort := port.NewDomainPort(ac.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Could not add service:", err)
		return
	}

	// Check processes' names
	for _, p := range processes {
		if p.Name == validatedConf.Name {
			fmt.Printf("Could not add service: an service with name \"%s\" already exists\n", validatedConf.Name)
			return
		}
	}

	// Convert the configurations to a map type
	// First, to json
	jsondata, err := json.Marshal(validatedConf)
	if err != nil {
		fmt.Println("Could not add service:", err)
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

	_, err = domainPort.ServiceProcesses().Create(mapConf)
	if err != nil {
		fmt.Println("Could not add service:", err)
		return
	}
	fmt.Println("Your service was successfully added. Restart conoid to start accepting request")
}
