package api

import "database/sql"

type (
	API interface {
	}

	qontoAPI struct {
		db *sql.DB
	}
)

func NewAPI(db *sql.DB) *qontoAPI {
	return &qontoAPI{
		db: db,
	}
}
