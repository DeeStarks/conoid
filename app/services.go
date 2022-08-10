package app

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	// "github.com/DeeStarks/conoid/app/tools"
	port "github.com/DeeStarks/conoid/domain/ports"
)

type (
	// A map to store running apps, the key represents incoming network address,
	// and the value as the server's address
	RunningServices map[string][]string

	Services struct {
		running   RunningServices
		nextPN    int // Next available port number
		defaultDB *sql.DB
	}

	IServices interface {
		ServeServices()                           // Retrieve all Services and serve
		GetRunningServices() RunningServices      // Get all running Services
		GetServiceServers(string) []string        // Get all servers' address that a service runs on
		ConnectToServer(string) (net.Conn, error) // Connect to a service running locally
		ServeStatic(string) int                   // Serve static Services, and return their port numbers
	}
)

func InitServices(defaultDB *sql.DB) IServices {
	return &Services{
		running:   RunningServices{},
		nextPN:    12000, // New servers will listen on port number :12000 and above for static apps
		defaultDB: defaultDB,
	}
}

// Retrieve all Services and serve
func (s *Services) ServeServices() {
	// Serve the welcome page
	welcomePort := s.ServeStatic("./assets/welcome/")
	// The welcome page will be served by default on port 80
	s.running["[::1]:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	s.running["localhost:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	s.running["localhost"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	s.running["127.0.0.1:80"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}
	s.running["127.0.0.1"] = []string{fmt.Sprintf("[::1]:%d", welcomePort)}

	// Serve
	dbPort := port.NewDomainPort(s.defaultDB)
	services, err := dbPort.ServiceProcesses().RetrieveRunning()
	if err != nil {
		log.Println("Could not serve apps:", err)
	}

	for _, service := range services {
		if service.Type == "static" {
			// Serve static
			portNo := s.ServeStatic(service.RootDirectory)
			addr := fmt.Sprintf("127.0.0.1:%d", portNo)
			_, err := dbPort.ServiceProcesses().Update(service.Name, map[string]interface{}{
				"listeners": addr,
			})
			if err != nil {
				log.Println("Could not update service state:", err)
			}
			s.running[service.RemoteServer] = []string{addr}
		} else if service.Type == "server" {
			servers := []string{}
			// Connect to all listening servers
			for _, addr := range service.Listeners {
				_, err := s.ConnectToServer(addr)
				if err != nil {
					log.Printf("Could not connect to server address: %s; Error: %v\n", addr, err)
					continue
				}
				// Append servers to listening servers
				servers = append(servers, addr)
			}
			s.running[service.RemoteServer] = servers

		}

		// Tunnelling
	}
}

// Get all running Services
func (s *Services) GetRunningServices() RunningServices {
	return s.running
}

// Get a service's port number using the remote address
func (s *Services) GetServiceServers(remoteAddr string) []string {
	return s.running[remoteAddr]
}

// Connect to a service running locally
func (s *Services) ConnectToServer(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Serve static Services, and return their port numbers
func (s *Services) ServeStatic(dir string) int {
	fs := http.FileServer(http.Dir(dir))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	// Get and listen on the next port number
	portNo := s.nextPN
	for {
		// Dial the port number to see if it's available
		_, err := net.Dial("tcp", fmt.Sprintf("[::]:%d", portNo))
		if err != nil {
			go http.ListenAndServe(fmt.Sprintf(":%d", portNo), mux)
			break
		}
		// If it's already in use, try the next port
		portNo++
	}
	s.nextPN = portNo
	return portNo
}
