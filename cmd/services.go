package cmd

import (
	"os"
	"path/filepath"

	"github.com/DeeStarks/conoid/app/cli"
	"github.com/spf13/cobra"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage services",
		Long:  `Used to manage services running on conoid`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get command flags
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			// Add new service
			if f, _ := flags.GetBool("add"); f {
				// Open the app's configuration file "conoid.yml"
				wd, _ := os.Getwd()
				appFile := filepath.Join(wd, "conoid.yml")
				services.Add(appFile)
				return
			}
		},
	}

	// Application processes sub-command
	servicePsCmd = &cobra.Command{
		Use:   "ps",
		Short: "List running services",
		Long:  "List runnning services",
		Run: func(cmd *cobra.Command, args []string) {
			// If the "all" flag is passed, list all processes
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			// List all Services
			if f, _ := flags.GetBool("all"); f {
				services.ListAll()
				return
			}

			// List only running services
			services.ListRunning()
		},
	}
)

func init() {
	// serviceCmd
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().BoolP("add", "", false, "lookup \"conoid.yml\" file in the current directory and add service")

	// servicePsCmd
	serviceCmd.AddCommand(servicePsCmd)
	servicePsCmd.Flags().BoolP("all", "a", false, "list all services")
}
