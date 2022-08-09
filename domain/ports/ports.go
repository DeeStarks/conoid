package port

import (
	"database/sql"

	"github.com/DeeStarks/conoid/domain/repository"
)

type IDomainPort interface {
	AppProcesses() repository.AppProcess
}

type DomainPort struct {
	db *sql.DB
}

func NewDomainPort(db *sql.DB) IDomainPort {
	return DomainPort{
		db: db,
	}
}
