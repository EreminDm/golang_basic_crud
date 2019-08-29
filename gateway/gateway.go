package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/EreminDm/golang_basic_crud/database"
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

// ShowList make a list of personaldata.
func ShowList(w http.ResponseWriter, r *http.Request, collection *database.Collection) {
	result, err := database.SelectAll(r.Context(), collection)
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

// Insert for create Personal Data by preparing to insert new data to DB.
func Insert(w http.ResponseWriter, r *http.Request, collection *database.Collection) {
	var p *database.PersonalData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	insertResult, err := database.Insert(r.Context(), collection, p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("created %v document(s)", insertResult.InsertedID))
}

// ShowListByID returns personal data list by id.
func ShowListByID(w http.ResponseWriter, r *http.Request, collection *database.Collection) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := database.SelectOne(r.Context(), collection, "_id", idvalue)
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

// Update function to add changes to personal information using object ID.
func Update(w http.ResponseWriter, r *http.Request, collection *database.Collection) {
	var p *database.PersonalData
	// reading request body information.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	// Unmarshal body to personalData object.
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}
	updateResult, err := database.Update(r.Context(), collection, p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("update %v document(s) successfully", updateResult))
}

// Remove using url param id which is objectID in DB.
func Remove(w http.ResponseWriter, r *http.Request, collection *database.Collection) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := database.Remove(r.Context(), collection, idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("deleted %v document(s) successfully", result))
}
