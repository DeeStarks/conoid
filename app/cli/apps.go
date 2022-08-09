package cli

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DeeStarks/conoid/domain/ports"
	"github.com/DeeStarks/conoid/utils"
	"github.com/aquasecurity/table"
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
		fmt.Println("Error retrieve apps:", err)
		return
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
	fmt.Println(validatedConf)
}
