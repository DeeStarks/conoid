package port

import (
	"github.com/DeeStarks/conoid/domain/repository"
)

func (p DomainPort) AppProcesses() repository.AppProcess {
	return repository.AppProcess{
		DB: p.db,
	}
}