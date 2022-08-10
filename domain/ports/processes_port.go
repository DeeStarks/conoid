package port

import (
	"github.com/DeeStarks/conoid/domain/repository"
)

type ProcessesPort interface {
	RetrieveRunning() ([]repository.ServiceProcessModel, error)
	RetrieveAll() ([]repository.ServiceProcessModel, error)
	Create(map[string]interface{}) (repository.ServiceProcessModel, error)
	Update(string, map[string]interface{}) (repository.ServiceProcessModel, error)
}

func (p DomainPort) ServiceProcesses() ProcessesPort {
	return repository.ServiceProcess{
		DB: p.db,
	}
}
