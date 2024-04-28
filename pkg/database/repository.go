package database

import "database/sql"

type repositoryDB struct {
	db *sql.DB
}

func NewRepositoryDB(db *sql.DB) Repository {
	return repositoryDB{db: db}
}

type Repository interface {
	GetMasterTaxDeduction(cycleYear int) (MasterTaxDeduction, error)
}
