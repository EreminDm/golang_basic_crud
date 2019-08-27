package gateway

import (
	"crud/database"
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

// PersonalData include two methods working with data.
func PersonalDatas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case "POST" makes a list of personal data information.
	case "GET":
		result, err := database.SelectAllPersonalData(r.Context())
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
	// case "POST" for create Personal Data by preparing to insert new data to DB.
	case "POST":
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
		insertResult, err := database.InsertPersonalData(r.Context(), p)
		if err != nil {
			errRespons(w, http.StatusInternalServerError, err)
			return
		}
		successResponce(w, http.StatusCreated, fmt.Sprintf("Created %v document(s)", insertResult.InsertedID))
	}

}

// PersonalDataByID get personal data by id.
func PersonalDataByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := database.SelectPersonalData(r.Context(), "_id", idvalue)
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

// UpdatePersonalData function to add changes to personal information using object ID.
func UpdatePersonalData(w http.ResponseWriter, r *http.Request) {
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
	updateResult, err := database.UpdatePersonalDataByID(r.Context(), p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("Update %v document(s) successfully", updateResult))
}

// RemovePersonalData using url param id which is objectID in DB.
func RemovePersonalData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := database.DeletePersonalData(r.Context(), idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("Deleted %v document(s) successfully", result))
}
