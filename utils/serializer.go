package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConf struct {
	Name          string   `yaml:"name"`     // Name of application
	Renderer      string   `yaml:"renderer"` // "server" or "static"
	Servers       []string `yaml:"servers"`  // Load will be automatically balanced across servers. Required if "Renderer" is "server"
	RootDirectory string   `yaml:"root"`     // Path to root directory. Required for "static" rendering
	ClientAddr    string   `yaml:"client"`   // Client address
	Tunnel        bool     `yaml:"tunnel"`   // Auto-share service to a remote network. This will be redundant if the "RemoteClient" is set
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
