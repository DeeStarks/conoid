package config

import (
	"os"
	"runtime"
)

const (
	// Network
	TCP_PORT       = 5000
	MAX_CONN_COUNT = 100
)

var (
	// File System
	DATA_ROOT  string // Data storage
	DEFAULT_DB string // Database to use by default
)

func init() {
	switch runtime.GOOS {
	case "windows":
		DATA_ROOT = os.ExpandEnv(`C:\Program Files\Conoid\`)
		DEFAULT_DB = DATA_ROOT + `Default.db`
	case "darwin":
		DATA_ROOT = os.ExpandEnv(`$HOME/Library/Application Support/conoid/`)
		DEFAULT_DB = DATA_ROOT + `default.db`
	default:
		DATA_ROOT = os.ExpandEnv(`/var/lib/conoid/`)
		DEFAULT_DB = DATA_ROOT + `default.db`
	}
}
