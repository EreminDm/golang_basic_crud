package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/EreminDm/golang_basic_crud/database"
	"github.com/EreminDm/golang_basic_crud/gateway"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// parsing command line flags.
	conn, db := envf()

	// create context for db connection.
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	// geting db collection.
	collection, err := database.Connect(ctx, conn, db)
	if err != nil {
		log.Fatalf(`couldn't connect to database: %v`, err)
	}

	// port environment define to 8000.
	log.Fatalf(`server initialization fail: %v`, http.ListenAndServe(":8000", handler(collection)))
}

// handler for routing map navigation.
func handler(collection *mongo.Collection) http.Handler {
	// making new router.
	r := mux.NewRouter()

	// handling urls API
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { gateway.ShowList(w, r, collection) }).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { gateway.Insert(w, r, collection) }).Methods("POST")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { gateway.ShowListByID(w, r, collection) }).Methods("GET")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { gateway.Update(w, r, collection) }).Methods("PUT")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { gateway.Remove(w, r, collection) }).Methods("DELETE")
	return r
}

// envf reades command line flags for db connection,
// connectURI flag returns db connection URI, example: http://localhost:5432 ,
// databes flag returns db name.
func envf() (string, string) {
	var conn, db string
	flag.StringVar(&conn, "connectURI", "http://uri:port", "a striong var")
	flag.StringVar(&db, "database", "database_name", "a string var")
	flag.Parse()
	return conn, db
}
