package repository

import (
	"database/sql"
	"fmt"
	"pr-reviewer-assignment-service/internal/config"

	_ "github.com/lib/pq"
)

type repository struct {
	Db *sql.DB
}

func New(configDB config.ConfigDatabase) (*repository, error) {
	db, err := createDatabaseObject(configDB)
	if err != nil {
		return nil, err
	}

	repo := repository{Db: db}
	return &repo, nil
}

func createDatabaseObject(configDB config.ConfigDatabase) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.BuildDB_URL(configDB))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping() doesn't work. err = %v", err)
	}

	return db, nil
}
