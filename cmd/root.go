package cmd

import (
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
		Short: "Expose local development servers to the internet and serve static files",
		Long:  "CONOID:\nA simple HTTP server that can be used to serve static files. \nIt also provides TCP tunnelling through localtunnel to bypass a firewall or NAT,\nenabling local development servers be exposed to the internet.",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()

			// Show version
			if f, _ := flags.GetBool("version"); f {
				fmt.Println(config.CURRENT_VERSION)
				return
			}

			// Get the number of services that are tunnelled
			var tunnelled int
			rec, err := port.NewDomainPort().ServiceProcesses().RetrieveRunning()
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
				go app.NewServer(openConnsCh).Serve()
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
)

func Execute() {
	// Version
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show version")

	// Execute
	rootCmd.Execute()

}
