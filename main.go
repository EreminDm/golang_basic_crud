package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/EreminDm/golang_basic_crud/database"
	"github.com/EreminDm/golang_basic_crud/network"
	"github.com/gorilla/mux"
)

func main() {
	// parsing command line flags.
	conn, db := envf()

	// create context for db connection.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// geting db collection.
	err := database.Connect(ctx, conn, db)
	if err != nil {
		log.Fatalf(`couldn't connect to database: %v`, err)
	}
	// port environment define to 8000.
	log.Fatalf(`server initialization fail: %v`, http.ListenAndServe(":8000", handler()))
}

// handler for routing map navigation.
func handler() http.Handler {
	// making new router.
	r := mux.NewRouter()
	// handling urls API
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { network.List(w, r) }).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { network.Insert(w, r) }).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { network.Update(w, r) }).Methods("PUT")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { network.ByID(w, r) }).Methods("GET")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { network.Remove(w, r) }).Methods("DELETE")
	return r
}

// envf reades command line flags for db connection,
// connectURI flag returns db connection URI, example: localhost:27017 ,
// databes flag returns db name.
func envf() (string, string) {
	var conn, db string
	flag.StringVar(&conn, "connectURI", "localhost:27017", "a string var")
	flag.StringVar(&db, "database", "database_name", "a string var")
	flag.Parse()
	return conn, db
}
