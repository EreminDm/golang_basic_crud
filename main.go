package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/httphandler"
	"github.com/EreminDm/golang_basic_crud/mongo"
)

// main initializes connection to database using timeout context,
// makes communication between database, controller and http layouts.
func main() {
	// envf parsing command line flags & returns database URI connection and database name,
	// connURI = "192.168.99.100:27017",
	// dbName = "information".
	connURI, dbName := envf()
	// create context for db connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// returns mongo collection.
	m, err := mongo.Connect(ctx, connURI, dbName)
	if err != nil {
		log.Fatalf(`couldn't connect to database: %v`, err)
	}
	// returns controller provider.
	c := controller.New(m)
	// returns handler provider.
	h := httphandler.New(c)
	// port environment define to 8000.
	log.Fatalf(`server initialization fail: %v`, http.ListenAndServe(":8000", h))
}

// envf reades command line flags for database connection,
// connectURI flag returns database connection URI, example: localhost:27017,
// databes flag returns database name.
func envf() (string, string) {
	var conn, db string
	flag.StringVar(
		&conn,
		"connectURI",
		"localhost:27017",
		"-connectURI flag, example: -connectURI=localhost:27017",
	)
	flag.StringVar(
		&db,
		"database",
		"information",
		"-database_name flag is a name of work database, example: -database_name=database_name_here",
	)
	flag.Parse()
	return conn, db
}
