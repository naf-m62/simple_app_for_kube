package database

import "database/sql"

type (
	UserRepo interface {
		GetName(userID int64) (string, error)
		Save(name string) (int64, error)
	}
	userRepo struct {
		db *sql.DB
	}
)

func (u *userRepo) GetName(userID int64) (name string, err error) {
	err = u.db.QueryRow("SELECT name FROM users WHERE id = $1", userID).Scan(&name)
	return name, err
}

func (u *userRepo) Save(name string) (userID int64, err error) {
	err = u.db.QueryRow("INSERT INTO users(name) VALUES ($1) RETURNING id", name).Scan(&userID)
	return userID, err
}
