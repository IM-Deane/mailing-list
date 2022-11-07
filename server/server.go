package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/IM-Deane/mailing-list/jsonapi"
	"github.com/IM-Deane/mailing-list/mdb"
	"github.com/alexflint/go-arg"
)

var args struct {
	DBPath string `arg:"env"MAILINGLIST_DB"`
	BindJSON string `arg:"env"MAILINGLIST_BIND_JSON"`
}

func main() {
	arg.MustParse(&args)

	// set defaults if env not provided
	if args.DBPath == "" {
		args.DBPath = "list.db"
	}
	if args.BindJSON == "" {
		args.BindJSON = ":8080"
	}

	// connect to DB
	log.Printf("using database '%v'", args.DBPath)
	db, err := sql.Open("sqlite3", args.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// close once function finished
	defer db.Close()

	mdb.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)
	// run concurrently
	go func() {
		log.Printf("starting JSON API server...\n")
		jsonapi.Serve(db, args.BindJSON)
		wg.Done()
	}()

	wg.Wait()
}