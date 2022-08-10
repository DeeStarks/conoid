package app

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/DeeStarks/conoid/app/tools"
	"github.com/DeeStarks/conoid/config"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	apps      IServices
	defaultDB *sql.DB
}

func NewServer() *Server {
	// Connect to the default db
	defaultDB, err := sql.Open("sqlite3", config.DEFAULT_DB)
	if err != nil {
		log.Panicln("Could not connect DB:", err)
	}

	// initialize and start running Services
	apps := InitServices(defaultDB)
	apps.ServeServices()

	return &Server{
		apps:      apps,
		defaultDB: defaultDB,
	}
}

func (s *Server) process(conn net.Conn) {
	// Get the servers
	addrs := s.apps.GetServiceServers(conn.RemoteAddr().String())
	if len(addrs) <= 0 {
		// If the remote address is unknown, redirect to the welcome server
		addrs = s.apps.GetServiceServers("[::1]:80")
	}

	// Get next server from load balancer
	lb := tools.NewLoadBalancer(addrs)
	addr := lb.GetNextServer()

	// Connect to the available server
	localConn, err := s.apps.ConnectToServer(addr)
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
