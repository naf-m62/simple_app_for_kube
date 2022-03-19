package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type (
	DBManager interface {
		UserRepo() UserRepo
	}

	dbManager struct {
		user UserRepo
	}
)

func NewDBManager(db *sql.DB) DBManager {
	return &dbManager{user: &userRepo{db: db}}
}

func (d *dbManager) UserRepo() UserRepo {
	return d.user
}
