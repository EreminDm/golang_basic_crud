package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	d "crud/db"
	gateway "crud/gateway"
)

// init function to initializ db connection.
func init() {
	err := d.MongodbURIConnection()
	if err != nil {
		log.Fatal(`Couldn't connect to DB`)
	}
}

// main function to init router.
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", gateway.PersonalDatas).Methods("GET")
	r.HandleFunc("/", gateway.PersonalDatas).Methods("POST")
	r.HandleFunc("/{id}", gateway.PersonalDataByID).Methods("GET")
	r.HandleFunc("/{id}", gateway.UpdatePersonalData).Methods("PUT")
	r.HandleFunc("/{id}", gateway.RemovePersonalData).Methods("DELETE")

	// port environment define to 8000.
	http.ListenAndServe(":8000", r)

}
