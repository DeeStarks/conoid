package cmd

import (
	"fmt"
	"strings"

	"github.com/DeeStarks/conoid/app/cli"
	"github.com/DeeStarks/conoid/utils"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new service",
		Long:  `Add a new service or application`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get command flags
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			if name, _ := flags.GetString("name"); name != "" {
				// Add new service
				var conf utils.AppConf
				var t string
				if t, _ = flags.GetString("type"); t == "" {
					fmt.Println("Invalid command: \"--type\" not specified")
					return
				}
				conf.Name = name
				conf.Type = t
				conf.Tunnelled, _ = flags.GetBool("tunnel")
				if d, _ := flags.GetString("directory"); d != "" {
					conf.RootDirectory = d
				}
				if l, _ := flags.GetString("listener"); l != "" {
					conf.Listeners = strings.Split(l, ",")
				}

				services.Add(conf, false)
				return
			}
			fmt.Println("Invalid command: \"--name\" not specified")
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a service",
		Long:  `Update service/application`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get command flags
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			if name, _ := flags.GetString("name"); name != "" {
				// Add new service
				var conf utils.AppConf
				conf.Name = name
				conf.Type, _ = flags.GetString("type")
				listener, _ := flags.GetString("listener")
				conf.Listeners = []string{listener}
				conf.RootDirectory, _ = flags.GetString("directory")
				conf.Tunnelled, _ = flags.GetBool("tunnel")

				services.Add(conf, true)
				return
			}
			fmt.Println("Invalid command: \"--name\" not specified")
		},
	}

	// Services processes
	psCmd = &cobra.Command{
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
	startCmd = &cobra.Command{
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
			fmt.Println("Invalid command. Execute \"conoid start -h\" for help")
		},
	}

	// Service stop
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a running service",
		Long:  "Stop a running service",
		Run: func(cmd *cobra.Command, args []string) {
			// If the "all" flag is passed, list all processes
			flags := cmd.Flags()

			// Initialize "service" commands
			services := cli.NewCLICommands().Services()

			if f, _ := flags.GetString("name"); f != "" {
				services.Stop(f)
				return
			}
			fmt.Println("Invalid command. Execute \"conoid stop -h\" for help")
		},
	}
)

func init() {
	// Add service
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("name", "n", "", "name of application to add")
	addCmd.Flags().StringP("type", "t", "", "application render type (\"static\" or \"server\")")
	addCmd.Flags().StringP("listener", "l", "", "address to listen on (e.g. localhost:8080) - required if \"type\" is set to \"server\"")
	addCmd.Flags().StringP("directory", "d", "", "document root of static application to serve (e.g. use \".\" for the current directory) - required if \"type\" is set to \"static\"")
	addCmd.Flags().BoolP("tunnel", "", false, "share your app to the internet")

	// Update service
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("name", "n", "", "name of service to update")
	updateCmd.Flags().StringP("type", "t", "", "render type (\"static\" or \"server\")")
	updateCmd.Flags().StringP("listener", "l", "", "address to listen on (e.g. localhost:8080) - required if \"type\" is set to \"server\"")
	updateCmd.Flags().StringP("directory", "d", "", "document root of static application to serve (e.g. use \".\" for the current directory) - required if \"type\" is set to \"static\"")
	updateCmd.Flags().BoolP("tunnel", "", false, "share your app to the internet")

	// Start a stopped service
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("name", "n", "", "name of service to start")

	// Start a stopped service
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringP("name", "n", "", "name of running service to stop")

	// Processes
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolP("all", "a", false, "list all services")
	psCmd.Flags().StringP("name", "n", "", "show details of a service")
}
