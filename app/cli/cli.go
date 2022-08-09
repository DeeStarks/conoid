package cli

type ICLICommands interface {
	Apps() *AppCommand
}

type CLICommands struct{}

// Accept and process CLI commands.
func NewCLICommands() ICLICommands {
	return &CLICommands{}
}
