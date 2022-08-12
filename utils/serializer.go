package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type AppConf struct {
	Name          string   `json:"name"`                     // Name of application
	Type          string   `json:"type"`                     // "server" or "static"
	Listeners     []string `json:"listeners,omitempty"`      // Load will be automatically balanced across listners. Required if "Renderer" is "server"
	RootDirectory string   `json:"root_directory,omitempty"` // Path to root directory. Required for "static" rendering
	Tunnelled     bool     `json:"tunnelled"`                // Share service to the internet
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
			return AppConf{}, errors.New("type of \"server\" requires \"--listener\" or \"-l\" (e.g. \"[options] --listener localhost:8000\")")
		}

		// Set the static root as empty
		conf.RootDirectory = ""
	}

	// RootDirectory
	if conf.Type == "static" {
		if conf.RootDirectory == "" {
			return AppConf{}, errors.New("expected \"--directory\" or \"-d\" for \"static\" types")
		}

		// Get directory's absolute path
		wd, _ := os.Getwd()
		root := filepath.Join(wd, conf.RootDirectory)
		conf.RootDirectory = root

		// Empty listeners
		conf.Listeners = []string{}
	}
	return conf, nil
}
