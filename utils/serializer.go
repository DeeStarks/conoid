package utils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type AppConf struct {
	Name          string   `json:"name"`         // Name of application
	Type          string   `json:"type"`         // "server" or "static"
	Listeners     []string `json:"listeners"`    // Required if "Renderer" is "server"
	RootDirectory string   `json:"root"`         // Path to root directory. Required for "static" rendering
	Tunnelled     bool     `json:"is_tunnelled"` // Share service to the internet
	Server        string   `json:"server"`       // Leave empty. If tunnelled, the domain name will be automatically added
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

		// Validate the addresses
		for i, l := range conf.Listeners {
			pUrl, err := url.Parse(l)
			if err != nil {
				return AppConf{}, fmt.Errorf("invalid address: %s", err)
			}
			conf.Listeners[i] = pUrl.Host
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
