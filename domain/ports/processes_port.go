package port

import (
	"github.com/DeeStarks/conoid/domain/repository"
)

type ProcessesPort interface {
	RetrieveRunning() ([]repository.AppProcessModel, error)
	RetrieveAll() ([]repository.AppProcessModel, error)
}

func (p DomainPort) AppProcesses() ProcessesPort {
	return repository.AppProcess{
		DB: p.db,
	}
}
