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
	RemoteServer  string
	Tunnelled     bool
	CreatedAt     int64
}

// List running services
func (c *ServiceCommand) ListRunning() {
	domainPort := port.NewDomainPort(c.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveRunning()
	if err != nil {
		fmt.Println("Error retrieve running services:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)

	t.SetHeaders("NAME", "TYPE", "LISTENERS", "ROOT", "REMOTE SERVER", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())
		listeners := strings.Join(p.Listeners, ", ")

		// Handle booleans
		tunnelled := "False"
		if p.Tunnelled {
			tunnelled = "True"
		}

		t.AddRow(
			p.Name, p.Type, utils.TruncateString(listeners, 20),
			utils.TruncateString(p.RootDirectory, 20),
			p.RemoteServer, tunnelled, created_at,
		)
	}
	t.Render()

}

// List all services
func (c *ServiceCommand) ListAll() {
	domainPort := port.NewDomainPort(c.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Error retrieving services:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)

	t.SetHeaders("NAME", "STATUS", "TYPE", "LISTENERS", "ROOT", "REMOTE SERVER", "TUNNELLED", "CREATED")
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
			p.Name, status, p.Type, utils.TruncateString(listeners, 20),
			utils.TruncateString(p.RootDirectory, 20),
			p.RemoteServer, tunnelled, created_at,
		)
	}
	t.Render()
}

// Add new service
func (c *ServiceCommand) Add(filepath string, update bool) {
	// Ensure the file exists
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("No file named: \"conoid.yml\"")
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
	domainPort := port.NewDomainPort(c.defaultDB)
	processes, err := domainPort.ServiceProcesses().RetrieveAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check processes' names
	var serviceToUpdate string
	for _, p := range processes {
		if p.Name == validatedConf.Name {
			if !update {
				fmt.Printf("A service already exists with the name \"%s\"; Use the \"--update\" flag to modify service\n", validatedConf.Name)
				return
			}

			// Store data for update
			serviceToUpdate = p.Name
			break
		}
	}

	// Convert the configurations to a map type
	// First, to json
	jsondata, err := json.Marshal(validatedConf)
	if err != nil {
		fmt.Println(err)
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

	if !update {
		// Generate pid, status, and created at
		mapConf["pid"] = strings.ReplaceAll(uuid.New().String(), "-", "") // Stripping hyphens
		mapConf["status"] = 1
		mapConf["created_at"] = time.Now().Unix()

		_, err = domainPort.ServiceProcesses().Create(mapConf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Your service was successfully added. Restart conoid to start accepting request")
	} else {
		// Updating service
		_, err := domainPort.ServiceProcesses().Update(serviceToUpdate, mapConf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\"%s\" service was updated. Restart conoid to synchronize update\n", serviceToUpdate)
	}
}

// Retrieve details of a servive
func (c *ServiceCommand) Get(name string) {
	// Retrieve from db
	domainPort := port.NewDomainPort(c.defaultDB)
	service, err := domainPort.ServiceProcesses().Get(name)
	if err != nil {
		fmt.Printf("Could not retrieve service: \"%s\"\n", name)
		return
	}

	// Convert booleans to string type
	tunnelled := "False"
	if service.Tunnelled {
		tunnelled = "True"
	}
	status := "Stopped"
	if service.Status {
		status = "Running"
	}

	// Show table
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)
	t.AddRow("NAME", service.Name)
	t.AddRow("STATUS", status)
	t.AddRow("TYPE", strings.Title(service.Type)+" rendering")
	t.AddRow("LISTENING ON", strings.Join(service.Listeners, ", "))
	if service.Type == "static" {
		t.AddRow("DOCUMENT ROOT", service.RootDirectory)
	}
	t.AddRow("REMOTE SERVER", service.RemoteServer)
	t.AddRow("TUNNELLED", tunnelled)
	t.AddRow("CREATED", utils.TimeAgo(service.CreatedAt, time.Now().Unix()))
	t.Render()
}
