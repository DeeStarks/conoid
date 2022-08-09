package config

import (
	"os"
	"runtime"
)

const (
	// App info
	VERSION = "0.0.1"

	// Network
	TCP_PORT = 80
)

var (
	// File System
	FS_ROOT    string
	PROCESS_DB string // DB for app processes
)

func init() {
	switch runtime.GOOS {
	case "windows":
		FS_ROOT = os.ExpandEnv(`C:\Program Files\Conoid\`)
		PROCESS_DB = FS_ROOT + `AppProcesses.db`
	case "darwin":
		FS_ROOT = os.ExpandEnv(`$HOME/Library/Conoid/`)
		PROCESS_DB = FS_ROOT + `AppProcesses.db`
	default:
		FS_ROOT = os.ExpandEnv(`/var/lib/conoid/`)
		PROCESS_DB = FS_ROOT + `app_processes.db`
	}
}
