package app

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/DeeStarks/conoid/config"
)

type Server struct {
	apps IApplications
}

func NewServer() *Server {
	apps := InitApplications() // initialize applications
	apps.RegisterApps()        // Register all apps

	return &Server{
		apps: apps,
	}
}

func (s *Server) process(conn net.Conn) {
	// Get the local port number of the server serving the remote client
	addrs := s.apps.GetAppServers(conn.RemoteAddr().String())
	if len(addrs) <= 0 {
		// If the remote address is unknown, redirect to the welcome server
		addrs = s.apps.GetAppServers("[::1]:80")
	}
	// Connect to the available server
	localConn, err := s.apps.ConnectToServer(addrs[0])
	if err != nil {
		log.Println(err)
		return
	}

	// Establish a point-to-point connection between conoid server and app's local server
	go func() {
		for {
			_, err = io.Copy(localConn, conn)
			if err != nil {
				log.Println("Failed to read from remote connection:", err)
				break
			}
		}
	}()

	go func() {
		for {
			_, err = io.Copy(conn, localConn)
			if err != nil {
				log.Println("Failed to write to remote connection:", err)
				break
			}
		}
	}()
}

func (s *Server) Serve() {
	// Start the server and wait for connections
	listener, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", config.TCP_PORT))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Conoid started and listening on port %d\n", config.TCP_PORT)

	for {
		// Establish a point-to-point connection between the client and server
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection failed:", err)
			continue
		}

		// Handle connection in a new goroutine
		go s.process(conn)
	}
}
