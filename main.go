package main

import (
	"log"
	"os"

	"github.com/DeeStarks/conoid/cmd"
	"github.com/DeeStarks/conoid/config"
)

// Setup filesystem during new installation
func SetupFs() {
	if _, err := os.Stat(config.FS_ROOT); os.IsNotExist(err) {
		err := os.Mkdir(config.FS_ROOT, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Create dependency files
	depFiles := []string{
		config.PROCESS_DB,
	}
	for _, dep := range depFiles {
		f, err := os.OpenFile(dep, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Error creating dependency file: %s; Error: %v\n", dep, err)
			return
		}
		f.Close()
	}
}

func main() {
	// Setup filesystem
	SetupFs()

	// Start process based on command argument
	cmd.Execute()
}
