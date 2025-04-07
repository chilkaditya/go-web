package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
func initDB()  {
	var err error
	db, err = sql.Open("sqlite3","./MusicAlbums.db")
	if err != nil {
		log.Fatal("Failed to connect to database:",err)

	}
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS MusicAlbums (
		id TEXT PRIMARY KEY,
		title TEXT,
		artist TEXT,
		movie_name TEXT,
		language TEXT

	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:",err)
	}
}