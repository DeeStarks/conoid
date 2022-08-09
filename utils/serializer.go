package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConf struct {
	Name          string   `yaml:"name"`      // Name of application
	Type          string   `yaml:"type"`      // "server" or "static"
	Listeners     []string `yaml:"listeners"` // Load will be automatically balanced across listners. Required if "Renderer" is "server"
	RootDirectory string   `yaml:"root"`      // Path to root directory. Required for "static" rendering
	ClientAddr    string   `yaml:"client"`    // Client address
	Tunnelled     bool     `yaml:"tunnelled"` // Share service to a remote network. This will be redundant if the "ClientAddr" is set
}

// Deserialize app yaml file
func DeserializeAppYAML(file *os.File) (AppConf, error) {
	conf := AppConf{}

	// Get file byte data
	buf := make([]byte, 512)
	n, err := file.Read(buf)

	if err != nil {
		log.Println("Error reading configuration file:", err)
	}

	// Unmarchal
	err = yaml.Unmarshal(buf[:n], &conf)
	if err != nil {
		log.Println("Error deserializing configuration file:", err)
	}
	return conf, nil
}
