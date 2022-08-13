package cli

type (
	ICLICommands interface {
		Services() *ServiceCommand
	}
	CLICommands struct{}
)

// Accept and process CLI commands.
func NewCLICommands() ICLICommands {
	return &CLICommands{}
}
