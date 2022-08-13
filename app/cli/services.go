package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	port "github.com/DeeStarks/conoid/domain/ports"
	"github.com/DeeStarks/conoid/utils"
	"github.com/aquasecurity/table"
)

type ServiceCommand struct{}

func (cmd *CLICommands) Services() *ServiceCommand {
	return &ServiceCommand{}
}

type ServiceProcess struct {
	Pid           string
	Name          string
	Status        string
	Type          string
	Listeners     []interface{}
	RootDirectory string
	RemoteServer  string
	Tunnelled     bool
	CreatedAt     int64
}

// List running services
func (c *ServiceCommand) ListRunning() {
	domainPort := port.NewDomainPort()
	processes, err := domainPort.ServiceProcesses().RetrieveRunning()
	if err != nil {
		fmt.Println("Error retrieving running services:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)

	t.SetHeaders("NAME", "TYPE", "LISTENING ON", "ROOT", "REMOTE ADDRESS", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())

		assertListeners := make([]string, len(p.Listeners))
		for i, l := range p.Listeners {
			assertListeners[i] = l.(string)
		}
		listeners := strings.Join(assertListeners, ", ")

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
	domainPort := port.NewDomainPort()
	processes, err := domainPort.ServiceProcesses().RetrieveAll()
	if err != nil {
		fmt.Println("Error retrieving services:", err)
		return
	}

	// Draw table to list processes
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)

	t.SetHeaders("NAME", "STATUS", "TYPE", "LISTENING ON", "ROOT", "REMOTE ADDRESS", "TUNNELLED", "CREATED")
	for _, p := range processes {
		created_at := utils.TimeAgo(p.CreatedAt, time.Now().Unix())

		assertListeners := make([]string, len(p.Listeners))
		for i, l := range p.Listeners {
			assertListeners[i] = l.(string)
		}
		listeners := strings.Join(assertListeners, ", ")

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
func (c *ServiceCommand) Add(conf utils.AppConf, update bool) {
	// Vallidate configuration
	validatedConf, err := utils.ValidateConf(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Ensure service doesn't already exist
	domainPort := port.NewDomainPort()
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
				fmt.Printf("A service already exists with the name \"%s\"; Use \"conoid update --name %s [flags]\" to modify service\n", validatedConf.Name, validatedConf.Name)
				return
			}

			// Store data for update
			serviceToUpdate = p.Name
			break
		}
	}

	// The service to be updated exists
	if update && serviceToUpdate == "" {
		fmt.Printf("No service named \"%s\". Do you mean to \"--add\"?\n", validatedConf.Name)
		return
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

	if !update {
		mapConf["status"] = true
		mapConf["created_at"] = time.Now().Unix()

		_, err = domainPort.ServiceProcesses().Create(mapConf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Your service was successfully added. Restart conoid to start accepting requests")
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
	domainPort := port.NewDomainPort()
	service, err := domainPort.ServiceProcesses().Get(name)
	if err != nil {
		fmt.Println(err)
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

	// Type assertion on the service listeners: interface{} -> string
	assertListeners := make([]string, len(service.Listeners))
	for i, l := range service.Listeners {
		assertListeners[i] = l.(string)
	}
	listeners := strings.Join(assertListeners, ", ")

	// Show table
	t := table.New(os.Stdout)
	t.SetDividers(table.UnicodeDividers)
	t.AddRow("NAME", service.Name)
	t.AddRow("STATUS", status)
	t.AddRow("TYPE", strings.Title(service.Type)+" rendering")
	t.AddRow("SERVING FROM", listeners)
	if service.Type == "static" {
		t.AddRow("DOCUMENT ROOT", service.RootDirectory)
	}
	t.AddRow("REMOTE ADDRESS", service.RemoteServer)
	t.AddRow("TUNNELLED", tunnelled)
	t.AddRow("CREATED", utils.TimeAgo(service.CreatedAt, time.Now().Unix()))
	t.Render()
}

// Start a stopped servive
func (c *ServiceCommand) Start(name string) {
	// Retrieve from db
	domainPort := port.NewDomainPort()
	service, err := domainPort.ServiceProcesses().Get(name)
	if err != nil {
		fmt.Printf("No such service: \"%s\"\n", name)
		return
	}

	// Check if service is already running
	if service.Status {
		fmt.Printf("Can't restart a running service: \"%s\"\n", service.Name)
		return
	}

	// Change status
	_, err = domainPort.ServiceProcesses().Update(name, map[string]interface{}{
		"status": true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s: service restarted. Restart conoid server to synchronize update\n", name)
}

// Stop a running servive
func (c *ServiceCommand) Stop(name string) {
	// Retrieve from db
	domainPort := port.NewDomainPort()
	service, err := domainPort.ServiceProcesses().Get(name)
	if err != nil {
		fmt.Printf("No such service: \"%s\"\n", name)
		return
	}

	// Check if service is already running
	if !service.Status {
		fmt.Println("Service stopped")
		return
	}

	// Change status
	_, err = domainPort.ServiceProcesses().Update(name, map[string]interface{}{
		"status": false,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Service stopped")
}
