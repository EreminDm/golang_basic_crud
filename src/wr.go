package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func errRespons(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}
func successResponce(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

// personalData getes list of personal data info.
func personalData(w http.ResponseWriter, r *http.Request) {
	result, err := selectAllPersonalData(r.Context())
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	responceBody, err := json.Marshal(result)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responceBody)
}

// personalDataByID get personal data by id
func personalDataByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := selectPersonalData(r.Context(), "_id", idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	responceBody, err := json.Marshal(result)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responceBody)
}

// createPersonalData prepare to insert new data to DB
func createPersonalData(w http.ResponseWriter, r *http.Request) {
	var p PersonalData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	insertResult, err := insertPersonalData(r.Context(), p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("Created %v document(s)", insertResult.InsertedID))
}

// updatePersonalData
func updatePersonalData(w http.ResponseWriter, r *http.Request) {
	var p PersonalData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	updateResult, err := updatePersonalDataByID(r.Context(), p.DocumentID, p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("Update %v document(s) successfully", updateResult))
}

func removePersonalData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := deletePersonalData(r.Context(), idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("Deleted %v document(s) successfully", result))
}
