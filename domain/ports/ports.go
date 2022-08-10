package port

import (
	"database/sql"
)

type IDomainPort interface {
	ServiceProcesses() ProcessesPort
}

type DomainPort struct {
	db *sql.DB
}

func NewDomainPort(db *sql.DB) IDomainPort {
	return DomainPort{
		db: db,
	}
}
