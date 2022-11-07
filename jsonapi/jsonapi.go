package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
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