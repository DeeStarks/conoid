package cmd

import (
	"path/filepath"

	"github.com/DeeStarks/conoid/app/cli"
	"github.com/spf13/cobra"
)

var (
	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Manage applications",
		Long:  `Used to manage applications running on conoid`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get command flags
			flags := cmd.Flags()

			// Initialize "app" commands
			appCmds := cli.NewCLICommands().Apps()

			// Add new application
			if f, _ := flags.GetString("add"); f != "" {
				// Open the app's configuration file "conoid.yml"
				appFile := filepath.Join(f, "conoid.yml")
				appCmds.Add(appFile)
				return
			}
		},
	}

	// Application processes sub-command
	appPsCmd = &cobra.Command{
		Use:   "ps",
		Short: "List running applications",
		Long:  "List runnning applications",
		Run: func(cmd *cobra.Command, args []string) {
			// If the "all" flag is passed, list all processes
			flags := cmd.Flags()

			// Initialize "app" commands
			appCmds := cli.NewCLICommands().Apps()

			// List all applications
			if f, _ := flags.GetBool("all"); f {
				appCmds.ListAll()
				return
			}

			// List only running applications
			appCmds.ListRunning()
		},
	}
)

func init() {
	// appCmd
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringP("add", "a", "", "add new application. Pass application's root directory")

	// appPsCmd
	appCmd.AddCommand(appPsCmd)
	appPsCmd.Flags().BoolP("all", "a", false, "list all application")
}
