package port

import (
	"github.com/DeeStarks/conoid/domain/repository"
)

type ProcessesPort interface {
	RetrieveRunning() ([]repository.ServiceProcessModel, error)
	RetrieveAll() ([]repository.ServiceProcessModel, error)
	// Pass data in paramaters to create
	Create(map[string]interface{}) (repository.ServiceProcessModel, error)
	// Pass the service name and data in paramaters to update
	Update(string, map[string]interface{}) (repository.ServiceProcessModel, error)
	// Pass the service name to be retrieved
	Get(string) (repository.ServiceProcessModel, error)
}

func (p DomainPort) ServiceProcesses() ProcessesPort {
	return repository.ServiceProcess{
		DB: p.db,
	}
}
