package config

import (
	"os"
	"runtime"
)

const (
	// App info
	CURRENT_VERSION = "v0.0.2"

	// Network
	TCP_PORT       = 5000
	MAX_CONN_COUNT = 100
)

var (
	// File System
	DATA_ROOT  string // Data storage
	LOGS_ROOT  string // Log files
	SERVICE_LOGS  string // Log files
	DEFAULT_DB string // Database to use by default
)

func init() {
	switch runtime.GOOS {
	case "windows":
		DATA_ROOT = os.ExpandEnv(`C:\Program Files\Conoid\`)
		LOGS_ROOT = DATA_ROOT + `logs`
		SERVICE_LOGS = DATA_ROOT + `logs\services.log`
		DEFAULT_DB = DATA_ROOT + `Default.db`
	case "darwin":
		DATA_ROOT = os.ExpandEnv(`$HOME/Library/Application Support/conoid/`)
		LOGS_ROOT = DATA_ROOT + `logs`
		SERVICE_LOGS = DATA_ROOT + `logs/services.log`
		DEFAULT_DB = DATA_ROOT + `default.db`
	default:
		DATA_ROOT = os.ExpandEnv(`/var/lib/conoid/`)
		LOGS_ROOT = DATA_ROOT + `logs`
		SERVICE_LOGS = DATA_ROOT + `logs/services.log`
		DEFAULT_DB = DATA_ROOT + `default.db`
	}
}
