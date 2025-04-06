package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
func initDB()  {
	var err error
	db, err = sql.Open("sqlite3","./albums.db")
	if err != nil {
		log.Fatal("Failed to connect to database:",err)

	}
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS albums (
		id TEXT PRIMARY KEY,
		title TEXT,
		artist TEXT,
		price REAL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:",err)
	}
}