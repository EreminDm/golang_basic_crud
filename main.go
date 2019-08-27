package main

import (
	"crud/database"
	"crud/gateway"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// init function to initializ db connection.
func init() {
	err := database.MongodbURIConnection()
	if err != nil {
		log.Fatal(`Couldn't connect to DB`)
	}
}

// main function to init router.
func main() {

	// port environment define to 8000.
	err := http.ListenAndServe(":8000", handler())
	if err != nil {
		log.Fatal(err)
	}

}

func handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", gateway.PersonalDatas).Methods("GET")
	r.HandleFunc("/", gateway.PersonalDatas).Methods("POST")
	r.HandleFunc("/{id}", gateway.PersonalDataByID).Methods("GET")
	r.HandleFunc("/{id}", gateway.UpdatePersonalData).Methods("PUT")
	r.HandleFunc("/{id}", gateway.RemovePersonalData).Methods("DELETE")
	return r
}
