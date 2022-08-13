package port

import (
	"os"

	"github.com/DeeStarks/conoid/config"
	"github.com/DeeStarks/conoid/domain/repository"
	"github.com/DeeStarks/conoid/utils"
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
	db, err := os.OpenFile(config.DEFAULT_DB, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		utils.Log("Could not connect to db:", err)
	}
	return repository.ServiceProcess{
		DB: db,
	}
}
