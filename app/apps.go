package app

import (
	"fmt"
	"net"
	"net/http"
)

type (
	// A map to store running apps, the key represents incoming network address,
	// and the value as the server's address
	RunningApps map[string][]string

	Applications struct {
		running RunningApps
		nextPN  int // Next available port number
	}

	IApplications interface {
		RegisterApps()                            // Retrieve all applications and serve
		GetRunningApps() RunningApps              // Get all running applications
		GetAppServers(string) []string            // Get all servers' address that an application runs on
		ConnectToServer(string) (net.Conn, error) // Connect to an application running locally
		ServeStatic(string) int                   // Serve static applications, and return their port numbers
	}
)

func InitApplications() IApplications {
	return &Applications{
		running: RunningApps{},
		nextPN:  12000, // New servers will listen on port number :12000 and above
	}
}

// Retrieve all applications and serve
func (a *Applications) RegisterApps() {
	// Serve the welcome page
	welcomePort := a.ServeStatic("./assets/welcome/")
	// We'll view the welcome page on the default port :80
	// so we add the port to every possible host
	a.running["[::1]:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	a.running["localhost:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	a.running["localhost"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	a.running["127.0.0.1:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	a.running["127.0.0.1"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}

	// Serve apps

}

// Get all running applications
func (a *Applications) GetRunningApps() RunningApps {
	return a.running
}

// Get an application's port number using the remote address
func (a *Applications) GetAppServers(remoteAddr string) []string {
	return a.running[remoteAddr]
}

// Connect to an application running locally
func (a *Applications) ConnectToServer(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Serve static applications, and return their port numbers
func (a *Applications) ServeStatic(dir string) int {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	// Get and listen on the next port number
	portNo := a.nextPN
	for {
		// Dial the port number to see if it's available
		_, err := net.Dial("tcp", fmt.Sprintf("[::]:%d", portNo))
		if err != nil {
			go http.ListenAndServe(fmt.Sprintf(":%d", portNo), fs)
			break
		}
		// If it's already in use, try the next port
		portNo++
	}
	a.nextPN = portNo
	return portNo
}
