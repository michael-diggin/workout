package database

import (
	"database/sql"
	"log"
)

// EnsureTableExists executes table creation query
func EnsureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS events
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT,
	sport TEXT,
	title TEXT,
	duration INTEGER
)`
