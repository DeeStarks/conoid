package config

import (
	"os"
	"runtime"
)

const (
	// App info
	VERSION = "0.0.1"

	// Network
	TCP_PORT       = 80
	MAX_CONN_COUNT = 100
)

var (
	// File System
	FS_ROOT    string
	DEFAULT_DB string // Database to use by default
)

func init() {
	switch runtime.GOOS {
	case "windows":
		FS_ROOT = os.ExpandEnv(`C:\Program Files\Conoid\`)
		DEFAULT_DB = FS_ROOT + `Default.db`
	case "darwin":
		FS_ROOT = os.ExpandEnv(`$HOME/Library/Conoid/`)
		DEFAULT_DB = FS_ROOT + `Default.db`
	default:
		FS_ROOT = os.ExpandEnv(`/var/lib/conoid/`)
		DEFAULT_DB = FS_ROOT + `default.db`
	}
}
