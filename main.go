package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/EreminDm/golang-basic-crud/controller"
	"github.com/EreminDm/golang-basic-crud/db/mariadb"
	"github.com/EreminDm/golang-basic-crud/db/mongo"
	netgrpc "github.com/EreminDm/golang-basic-crud/net/grpc"
	nethttp "github.com/EreminDm/golang-basic-crud/net/http"
)

// main initializes connection to database using timeout context,
// makes communication between database, controller and http layouts.
func main() {
	// envf parsing command line flags & returns database URI connection and database name,
	// connURI = "192.168.99.100:27017",
	// dbName = "information".
	connURI, dbName, dbtype := envf()
	// create context for db connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var c *controller.Personal
	// switch between database types
	switch dbtype {
	case "mongo":
		// returns mongo collection.
		m, err := mongo.Connect(ctx, connURI, dbName)
		if err != nil {
			log.Fatalf(`couldn't connect to database: %v`, err)
		}
		// returns controller provider.
		c = controller.New(m)

	case "mariadb":
		// returns mariadb collection.
		m, err := mariadb.Connect(ctx, connURI, dbName, 3, 30)
		if err != nil {
			log.Fatalf(`couldn't connect to database: %v`, err)
		}
		// returns controller provider.
		c = controller.New(m)
	}
	// returns handler provider.
	h := nethttp.New(c)
	// start listen grpc server on port 8888.
	l, srv, err := netgrpc.ConnectServer(c)
	if err != nil {
		log.Fatalf("could not init grpc servers port: %v", err)
	}
	go log.Fatal(srv.Serve(l))
	// port environment define to 8000.
	log.Fatalf(`server initialization fail: %v`, http.ListenAndServe(":8000", h))
}

// envf reades command line flags for database connection,
// connectURI flag returns database connection URI, example: localhost:27017,
// databes flag returns database name.
func envf() (string, string, string) {
	var conn, db, dbtype string
	flag.StringVar(
		&conn,
		"connectURI",
		"localhost:27017",
		"-connectURI flag, example: -connectURI=localhost:27017",
	)
	flag.StringVar(
		&db,
		"database",
		"person",
		"-database_name flag is a name of work database, example: -database_name=database_name_here",
	)
	flag.StringVar(
		&dbtype,
		"dbtype",
		"mongo",
		"-dbtype flag, example: -dbtype=mongo or -dbtype=mariadb",
	)
	flag.Parse()
	return conn, db, dbtype
}
