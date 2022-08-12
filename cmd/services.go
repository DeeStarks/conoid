package cmd

import (
	"fmt"
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

			if f, _ := flags.GetBool("add"); f {
				// Add new service

				// Open the app's configuration file "conoid.yml"
				wd, _ := os.Getwd()
				pathToConf := filepath.Join(wd, "conoid.yml")
				// Check if such file exists
				if _, err := os.Stat(pathToConf); err != nil {
					fmt.Println("No file named \"conoid.yml\" in the current directory")
					return
				}
				services.Add(pathToConf, false)
				return
			} else if f, _ := flags.GetBool("update"); f {
				// Update existing service

				// Open the app's configuration file "conoid.yml"
				wd, _ := os.Getwd()
				pathToConf := filepath.Join(wd, "conoid.yml")
				// Check if such file exists
				if _, err := os.Stat(pathToConf); err != nil {
					fmt.Println("No file named \"conoid.yml\" in the current directory")
					return
				}
				services.Add(pathToConf, true)
				return
			}
			fmt.Println("No flags specified. Execute \"conoid service -h\" for help")
		},
	}

	// Services processes sub-command
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
			} else if f, _ := flags.GetBool("all"); f {
				// Retrieve a service
				services.ListAll()
			} else if f, _ := flags.GetString("name"); f != "" {
				services.Get(f)
			} else {
				// List running services
				services.ListRunning()
			}
		},
	}

	// Service restart
	serviceStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Restart a stopped service",
		Long:  "Restart a stopped service",
		Run: func(cmd *cobra.Command, args []string) {
			// If the "all" flag is passed, list all processes
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			if f, _ := flags.GetString("name"); f != "" {
				services.Start(f)
				return
			}
			fmt.Println("No flags specified. Execute \"conoid start -h\" for help")
		},
	}
)

func init() {
	// serviceCmd
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().BoolP("add", "", false, "lookup \"conoid.yml\" file in the current directory and add service")
	serviceCmd.Flags().BoolP("update", "", false, "modify service based on \"conoid.yml\"")

	// servicePsCmd
	serviceCmd.AddCommand(servicePsCmd)
	servicePsCmd.Flags().BoolP("all", "a", false, "list all services")
	servicePsCmd.Flags().StringP("name", "n", "", "show details of a service")

	// serviceStartCmd
	serviceCmd.AddCommand(serviceStartCmd)
	serviceStartCmd.Flags().StringP("name", "n", "", "name of service to restart")
}
