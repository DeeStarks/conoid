package app

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/DeeStarks/conoid/app/tools"
	port "github.com/DeeStarks/conoid/domain/ports"
	"github.com/DeeStarks/conoid/utils"
)

type (
	// A map to store running services, the key represents local server,
	// and the value as the remote server's address
	RunningServices map[string][]string

	Services struct {
		running   RunningServices
		defaultDB *sql.DB
	}

	IServices interface {
		ServeServices(string, string, chan<- net.Conn) // Retrieve all Services and serve
		GetRunningServices() RunningServices           // Get all running Services
		GetServiceServers(string) []string             // Get all servers' address that a service runs on
		ConnectToServer(string) (net.Conn, error)      // Connect to a service running locally
		ServeStatic(string, int) (string, string)      // Serve static Services, and return their port numbers
	}
)

func InitServices(defaultDB *sql.DB) IServices {
	return &Services{
		running:   RunningServices{},
		defaultDB: defaultDB,
	}
}

// Retrieve all Services and serve
func (s *Services) ServeServices(conoidHost, conoidPort string, connCh chan<- net.Conn) {
	// Serve the welcome page
	welcomeHost, welcomePort := s.ServeStatic("./assets/welcome/", 30000)
	// Set the welcome page as the werver's default page
	s.running[fmt.Sprintf("%s:%s", conoidHost, conoidPort)] = []string{fmt.Sprintf("%s:%s", welcomeHost, welcomePort)}

	// Serve registered running services
	dbPort := port.NewDomainPort(s.defaultDB)
	services, err := dbPort.ServiceProcesses().RetrieveRunning()
	if err != nil {
		utils.Log("Could not serve:", err)
		return
	}

	portNo := 30001
	for _, service := range services {
		// Addresses the service is running on
		var serverAddrs []string

		if service.Type == "static" {
			// Serve static
			host, port := s.ServeStatic(service.RootDirectory, portNo)
			addr := fmt.Sprintf("%s:%s", host, port)
			_, err := dbPort.ServiceProcesses().Update(service.Name, map[string]interface{}{
				"listeners": addr,
			})
			if err != nil {
				utils.Log("Could not update service state:", err)
			}
			serverAddrs = []string{addr}
			portNo++ // Increment the port number for nexe usage
		} else {
			servers := []string{}
			// Connect to all listening servers
			for _, addr := range service.Listeners {
				_, err := s.ConnectToServer(addr)
				if err != nil {
					utils.Logf("Could not connect to \"%s\" at: %s; Stopping...\n", service.Name, addr)
					// Update service state
					dbPort.ServiceProcesses().Update(service.Name, map[string]interface{}{
						"status": 0,
					})
					continue
				}
				// Append servers to listening servers
				servers = append(servers, addr)
			}
			serverAddrs = servers
		}

		// Tunnelling
		if service.Tunnelled {
			tunnel := tools.NewTunnel(service.Name, connCh)
			host, err := tunnel.AllocateHost()
			if err != nil {
				utils.Log("Error opening tunnel. Ensure your device is connected to the internet")
				continue
			}

			connectedAddressCh := make(chan string, 1)
			for i := 0; i < host.MaxConnectionCount(); i++ {
				go host.OpenTunnel(fmt.Sprintf("%s:%s", conoidHost, conoidPort), connectedAddressCh)

				// Block till the local address connected to the remote server is sent
				localConn := <-connectedAddressCh
				s.running[localConn] = serverAddrs
			}

			// Update service's remote_server
			_, err = dbPort.ServiceProcesses().Update(service.Name, map[string]interface{}{
				"remote_server": host.FullURL(),
			})
			if err != nil {
				panic(err)
			}
		}
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
func (s *Services) ServeStatic(dir string, port int) (string, string) {
	fs := http.FileServer(http.Dir(dir))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	// Get and listen on the next port number
	host := "localhost"
	// Dial the port number to see if it's available
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		// If it's not in use, serve
		go http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
	}
	return host, fmt.Sprintf("%d", port)
}
