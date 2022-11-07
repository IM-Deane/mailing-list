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

// emailEntryFromRow build an email entry from provided DB row
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


// CreateEmail adds new entry to email table
func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
		INSERT INTO
			emails(email, confirmed_at, opt_out)
		VALUES
			(?, 0, false)`, email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}


// GetEmail fetches email entry from DB
func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
	rows, err := db.Query(`
		SELECT
			id, email, confirmed_at, opt_out
		FROM
			emails
		WHERE
			email = ?`, email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// close DB connection if any error occurs
	defer rows.Close()

	// read new row from DB
	for rows.Next() {
		return emailEntryFromRow(rows)
	}

	return nil, nil
}


// UpdateEmail updates a given email entry or creates a new one if it doesn't exist
func UpdateEmail(db *sql.DB, entry EmailEntry) error {
	t := entry.ConfirmedAt.Unix()

	// UPSERT email (try to create new entry, if it exists update instead)
	_, err := db.Exec(`
		INSERT INTO
			emails(email, confirmed_at, opt_out)
		VALUES
			(?, ?, ?)
		ON CONFLICT(EMAIL) DO UPDATE SET
			confirmed_at=?
			opt_out=?`, entry.Email, t, entry.OptOut, t, entry.OptOut)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}


// DeleteEmail soft deletes email from mailing list
// NOTE: we keep the record to avoid edgecase where we
// send an email to someone that's already opted out (ie. spam).
func DeleteEmail(db *sql.DB, email string) error {
	// setting opt_out=true removes that email from the mailing list
	_, err := db.Exec(`
		UPDATE emails
		SET opt_out=true
		WHERE email=?`, email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type GetEmailBatchQueryParams struct {
	Page int
	Count int
}

// GetEmailBatch fetches all users currently subscribed to mailing list
func GetEmailBatch(db *sql.DB, params GetEmailBatchQueryParams) ([]EmailEntry, error) {
	var empty []EmailEntry

	// get current users offset by current page
	rows, err := db.Query(`
		SELECT
			id, email, confirmed_at, opt_out
		FROM
			emails
		WHERE
			opt_out = false
		ORDER BY id ASC
		LIMIT ? OFFSET ?`, params.Count, (params.Page-1)*params.Count)

	if err != nil {
		log.Println(err)
		return empty, err
	}

	// close DB connection on error or end of func
	defer rows.Close()

	// create slice of email list
	emails := make([]EmailEntry, 0, params.Count)

	// read rows from DB
	for rows.Next() {
		email, err := emailEntryFromRow(rows)
		if err != nil {
			// cancel iteration as we don't want a partial list
			return nil, err
		}
		emails = append(emails, *email)
	}

	return emails, nil
}