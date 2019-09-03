package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/gorilla/mux"
)

// PersonalData description.
type PersonalData struct {
	DocumentID  string `json:"id,omitempty"` // as *primitive.ObjectID.
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	YearOfBirth int    `json:"yaerOfBirth,omitempty"`
}

// errRespons returns error to responce.
func errRespons(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

// successResponce returns succsess responce.
func successResponce(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

// ByID returns personal data list by id.
func ByID(w http.ResponseWriter, r *http.Request) {
	var u controller.UsersPersonalData
	params := mux.Vars(r)
	id := params["id"]
	result, err := u.One(r.Context(), id)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	res, err := json.Marshal(result)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// List returns a list of personaldata.
func List(w http.ResponseWriter, r *http.Request) {
	var u controller.UsersPersonalData
	usrs, err := u.All(r.Context())
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	responceBody, err := json.Marshal(usrs)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responceBody)
}

// Insert creates Personal Data by preparing to insert new data to database.
func Insert(w http.ResponseWriter, r *http.Request) {
	var (
		u controller.UsersPersonalData
		p PersonalData
	)
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

	insertResult, err := u.Insert(r.Context(), transformToController(p))
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("created %v document(s)", insertResult))
}

// Update adds changes to personal information using object ID.
func Update(w http.ResponseWriter, r *http.Request) {
	var (
		u controller.UsersPersonalData
		p PersonalData
	)
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
	updateResult, err := u.Update(r.Context(), transformToController(p))
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("update %v document(s) successfully", updateResult))
}

// Remove deletes object using url param id which is objectID in database.
func Remove(w http.ResponseWriter, r *http.Request) {
	var u controller.UsersPersonalData
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := u.Remove(r.Context(), idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("deleted %v document(s) successfully", result))
}

// transformToController returns controllers personal data.
func transformToController(p PersonalData) *controller.PersonalData {
	var up controller.PersonalData
	up.DocumentID = p.DocumentID
	up.Email = p.Email
	up.LastName = p.LastName
	up.Name = p.Name
	up.Phone = p.Phone
	up.YearOfBirth = p.YearOfBirth
	return &up
}
