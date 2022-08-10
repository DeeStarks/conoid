package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type AppConf struct {
	Name          string   `yaml:"name" json:"name"`                             // Name of application
	Type          string   `yaml:"type" json:"type"`                             // "server" or "static"
	Listeners     []string `yaml:"listeners" json:"listeners,omitempty"`         // Load will be automatically balanced across listners. Required if "Renderer" is "server"
	RootDirectory string   `yaml:"root" json:"root_directory,omitempty"`         // Path to root directory. Required for "static" rendering
	RemoteServer  string   `yaml:"remote-server" json:"remote_server,omitempty"` // Address to accept and respond to requests
	Tunnelled     bool     `yaml:"tunnelled" json:"tunnelled"`                   // Share service to a remote network. This will be redundant if the "ClientAddr" is set
}

// Deserialize app yaml file
func DeserializeConf(file *os.File) (AppConf, error) {
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

// Validate and fix errors in configuration
func ValidateConf(conf AppConf) (AppConf, error) {
	// Set name to lowercase letters
	conf.Name = strings.ToLower(conf.Name)

	// Validate type
	conf.Type = strings.ToLower(conf.Type)
	if conf.Type != "server" && conf.Type != "static" {
		return AppConf{}, errors.New("expected type of \"server\" or \"static\"")
	}

	if conf.Type == "server" {
		if len(conf.Listeners) == 0 {
			return AppConf{}, errors.New("type of \"server\" requires at least a \"listener\" or a server address")
		}

		// Set the static root as empty
		conf.RootDirectory = ""
	}

	// RootDirectory
	if conf.Type == "static" {
		if conf.RootDirectory == "" {
			return AppConf{}, errors.New("expected \"root\" directory for a \"static\" type")
		}

		// Get directory's absolute path
		wd, _ := os.Getwd()
		root := filepath.Join(wd, conf.RootDirectory)
		conf.RootDirectory = root

		// Empty listeners
		conf.Listeners = []string{}
	}

	// Set "Tunnelled" to false if "ClientAddr" is passed
	if len(conf.RemoteServer) > 0 {
		conf.Tunnelled = false
	}
	return conf, nil
}
