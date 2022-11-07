package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/IM-Deane/mailing-list/mdb"
)

// setJSONHeader adds a json header to response writer
func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// fromJSON is a generic function that converts provided JSON to GO struct
func fromJSON[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	// converted bytes to target strcuct
	json.Unmarshal(buf.Bytes(), &target)
}

// returnJSON is a generic function that converts data encapsulated by 'withData'
// and returns a JSON response
func returnJSON[T any](w http.ResponseWriter, withData func() (T, error)) {
	setJSONHeader(w)

	data, serverErr := withData()

	if serverErr != nil {
		// server encountered
		w.WriteHeader(500)
		serverErrJSON, err := json.Marshal(&serverErr)
		if err != nil {
			// something went wrong
			log.Println(err)
			return
		}
		// return our custom server
		w.Write(serverErrJSON)
		return
	}

	dataJSON, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// return JSON response
	w.Write(dataJSON)
}

// returnErr returns an error response
func returnErr(w http.ResponseWriter, err error, httpCode int) {
	// can pass any type
	returnJSON(w, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(), // convert to string
		}
		w.WriteHeader(httpCode) // 4XX, 5XX, etc.
		return errorMessage, nil
	})
}

// CreateEmail adds email to DB and and returns a JSON response object
func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(r.Body, &entry)

		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}

		// get email as JSON
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON CreateEmail: %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})
	})
}

// GetEmail fetches an email from the DB as a JSON response
func GetEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(r.Body, &entry)

		// get email as JSON
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON GetEmail: %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})
	})
}

// UpdateEmail updates email in DB and and returns a JSON response object
func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(r.Body, &entry)

		if err := mdb.UpdateEmail(db, entry); err != nil {
			returnErr(w, err, 400)
			return
		}

		// get email as JSON
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON UpdateEmail: %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})
	})
}

// DeleteEmail removes email from mailing list and returns a JSON response object
func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(r.Body, &entry)

		if err := mdb.DeleteEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}

		// get email as JSON
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON DeleteEmail: %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})
	})
}


// GetEmailBatch fetches all emails in list as a JSON response
func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		queryOptions := mdb.GetEmailBatchQueryParams{}
		fromJSON(r.Body, &queryOptions)

		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("page and Count fields are required and must be > 0"), 400)	
			return
		}

		// return email list
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON GetEmailBatch: %v\n", queryOptions)
			return mdb.GetEmailBatch(db, queryOptions)
		})
	})
}


// Serve serves JSON handler functions
func Serve(db *sql.DB, bind string) {
	
	// handlers
	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/get_batch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))

	log.Printf("JSON API server listening on: %v", bind)
	
	// init server
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("JSON server error: %v", err)
	}
}