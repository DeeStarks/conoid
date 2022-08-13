package app

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/DeeStarks/conoid/app/tools"
	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/utils"
)

type Server struct {
	services  IServices
	host      string
	port      string
	openConns chan<- net.Conn
}

func NewServer(connCh chan<- net.Conn) *Server {
	// initialize and start running Services
	services := InitServices()

	return &Server{
		services:  services,
		openConns: connCh,
	}
}

func (s *Server) process(conn net.Conn) {
	// Get the servers
	addrs := s.services.GetServiceServers(conn.RemoteAddr().String())
	if len(addrs) <= 0 {
		// If the remote address is unknown, redirect to the welcome server
		addrs = s.services.GetServiceServers(fmt.Sprintf("%s:%s", s.host, s.port))
	}

	// TODO:
	// The idea for the load balancer isn't fully formed yet.
	// For now, it's always going to select the first server
	lb := tools.NewLoadBalancer(addrs)
	addr := lb.GetNextServer()

	// Connect to the available server
	localConn, err := s.services.ConnectToServer(addr)
	if err != nil {
		utils.Log(err)
		return
	}

	// Add local conn to open connections channel
	s.openConns <- localConn

	// Establish a point-to-point connection between conoid server and app's local server
	go func() {
		for {
			_, err = io.Copy(localConn, conn)
			if err != nil {
				break
			}
		}
	}()

	go func() {
		for {
			_, err = io.Copy(conn, localConn)
			if err != nil {
				break
			}
		}
	}()
}

func (s *Server) Serve() {
	// Start the server and wait for connections
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.TCP_PORT))
	if err != nil {
		utils.Log(err)
		os.Exit(1)
	}
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		utils.Log(err)
		return
	}
	s.host = host
	s.port = port

	// Start running services
	s.services.ServeServices(host, port, s.openConns)

	utils.Logf("Conoid listening on host: %s, port %s\n", host, port)
	// Record connections to ensure it doesn't exceed the max size
	connsCh := make(chan int, config.MAX_CONN_COUNT)

	for {
		// Block if connections count is full
		connsCh <- 1

		// Accept connections
		conn, err := listener.Accept()
		if err != nil {
			utils.Log("Connection failed:", err)
			// Remove record
			<-connsCh
			continue
		}

		// Add to open connections
		s.openConns <- conn

		// Handle connection in a new goroutine
		go s.process(conn)
	}
}
