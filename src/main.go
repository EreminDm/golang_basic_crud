package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	err := MongodbURIConnection()
	if err != nil {
		log.Fatal(`Couldn't connect to DB`)
	}

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", personalData).Methods("GET")
	r.HandleFunc("/", createPersonalData).Methods("POST")
	r.HandleFunc("/{id}", personalDataByID).Methods("GET")
	r.HandleFunc("/{id}", updatePersonalData).Methods("PUT")
	r.HandleFunc("/{id}", removePersonalData).Methods("DELETE")

	// PORT environment define to 8000
	http.ListenAndServe(":8000", r)

}
