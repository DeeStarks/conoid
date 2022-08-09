package cli

import (
	"log"
	"os"

	"github.com/DeeStarks/conoid/utils"
)

type AppCommand struct{}

func (cmd *CLICommands) Apps() *AppCommand {
	return &AppCommand{}
}

// List running applications
func (ac *AppCommand) ListRunning() {
	log.Println("Running processes")
}

// List all applications
func (ac *AppCommand) ListAll() {
	log.Println("All processes")
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
