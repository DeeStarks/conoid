package port

import (
	"github.com/DeeStarks/conoid/domain/repository"
)

type ProcessesPort interface {
	RetrieveRunning() []repository.AppProcessModel
	RetrieveAll() []repository.AppProcessModel
}

func (p DomainPort) AppProcesses() ProcessesPort {
	return repository.AppProcess{
		DB: p.db,
	}
}
