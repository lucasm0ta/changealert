package core

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type WatcherRepository struct {
}

func NewWatcherRepository() *WatcherRepository {
	repository := new(WatcherRepository)
	db, err := sql.Open("sqlite3", "./trex.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return repository
}
