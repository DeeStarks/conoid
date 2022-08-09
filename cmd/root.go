package cmd

import (
	"fmt"
	"log"

	"github.com/DeeStarks/conoid/app"
	"github.com/DeeStarks/conoid/config"
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
				log.Println(config.VERSION)
				return
			}

			// Start server if no argumentss were passed or the first argument is "up"
			if len(args) <= 0 || args[0] == "start" {
				app.NewServer().Serve()
			}
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Conoid",
		Long:  `All software has versions. This is Conoid's`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(config.VERSION)
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
