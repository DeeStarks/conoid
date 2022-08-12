package cmd

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/DeeStarks/conoid/app"
	"github.com/DeeStarks/conoid/config"
	port "github.com/DeeStarks/conoid/domain/ports"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "conoid",
		Short: "Conoid is a reverse proxy server that can also expose your localhost on the internet",
		Long:  fmt.Sprintf("CONOID (%s)\nA reverse proxy server with load balancer and caching feature.\nIt also uses localtunnel to help you easily expose web services on your\nlocal development machine without messing with DNS and firewall settings", config.VERSION),
		Run: func(cmd *cobra.Command, args []string) {
			// Print version
			if v, _ := cmd.Flags().GetBool("version"); v {
				fmt.Println(config.VERSION)
				return
			}

			// Connect to the default db
			defaultDB, err := sql.Open("sqlite3", config.DEFAULT_DB)
			if err != nil {
				panic(err)
			}

			// Get the number of services that are tunnelled
			var tunnelled int
			rec, err := port.NewDomainPort(defaultDB).ServiceProcesses().RetrieveRunning()
			if err != nil {
				panic(err)
			}
			for _, r := range rec {
				if r.Tunnelled {
					tunnelled++
				}
			}

			mltcc := 10       // Maximum localtunnel connections count
			mltcc = mltcc * 2 // Incoming and outgoing connections
			tunnelConnCount := tunnelled * mltcc
			// Record all open connections
			// This will be used to close all on shutdown
			openConnsCh := make(chan net.Conn, config.MAX_CONN_COUNT+tunnelConnCount)

			// Start server if no argumentss were passed or the first argument is "up"
			if len(args) <= 0 || args[0] == "start" {
				go app.NewServer(openConnsCh, defaultDB).Serve()
			}

			// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			// Block until signal is received
			<-c

			// Shutdown
			fmt.Println("Gracefully shutting down conoid...")
		L:
			for {
				select {
				case conn := <-openConnsCh:
					err := conn.Close()
					if err != nil {
						panic(err)
					}
				// Wait for all connections to close
				case <-time.After(time.Second * 5):
					break L
				}
			}
			fmt.Println("Done!")
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Current conoid version",
		Long:  `All software has versions. This is Conoid's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.VERSION)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().BoolP("version", "v", false, "conoid version")
}

func Execute() {
	rootCmd.Execute()
}
