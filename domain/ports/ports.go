package port

type (
	IDomainPort interface {
		ServiceProcesses() ProcessesPort
	}
	DomainPort struct{}
)

func NewDomainPort() IDomainPort {
	return DomainPort{}
}
