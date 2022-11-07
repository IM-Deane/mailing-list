package mdb

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

// EmailEntry schema for email entry database
type EmailEntry struct {
	ID int64
	Email string
	ConfirmedAt *time.Time
	OptOut bool
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE emails (
			id INTEGER PRIMARY KEY,
			email TEXT UNIQUE,
			confirmed_at INTEGER,
			opt_out INTEGER
		);
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 {
				// code 1 == table already exists
				log.Fatal(sqlError)
			}
		} else {
			// handle other errors
			log.Fatal(err)
		}
	}
}

func emailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool
	
	// get data from row
	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// convert time format
	t := time.Unix(confirmedAt, 0)

	return &EmailEntry{ID: id, Email: email, ConfirmedAt: &t, OptOut: optOut}, nil
}