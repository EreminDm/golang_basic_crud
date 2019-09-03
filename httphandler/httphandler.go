package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CTRL describes controller implementation.
type CTRL struct {
	CTR CTRprovider
}

// CTRprovider describes provider methods.
type CTRprovider interface {
	One(ctx context.Context, value string) (*PersonalData, error)
	All(ctx context.Context) (*[]PersonalData, error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p *PersonalData) (int64, error)
	Insert(ctx context.Context, document *PersonalData) (interface{}, error)
}

// NewCTRL returns new controller provider
func NewCTRL(c CTRprovider) (*CTRL, error) {
	return &CTRL{
		CTR: c,
	}, nil
}

// PersonalData is Personal Data filds description.
type PersonalData struct {
	DocumentID  string `json:"id,omitempty"` // as *primitive.ObjectID.
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	YearOfBirth int    `json:"yaerOfBirth,omitempty"`
}

// Handler is function for routing map navigation.
func Handler(c *CTRL) http.Handler {
	// making new router.
	r := mux.NewRouter()
	// handling urls API.
	r.HandleFunc("/", c.List).Methods("GET")
	r.HandleFunc("/", c.Insert).Methods("POST")
	r.HandleFunc("/", c.Update).Methods("PUT")
	r.HandleFunc("/{id}", c.ByID).Methods("GET")
	r.HandleFunc("/{id}", c.Remove).Methods("DELETE")
	return r
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
func (c *CTRL) ByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	result, err := c.CTR.One(r.Context(), id)
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
func (c *CTRL) List(w http.ResponseWriter, r *http.Request) {
	usrs, err := c.CTR.All(r.Context())
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
func (c *CTRL) Insert(w http.ResponseWriter, r *http.Request) {
	var p PersonalData
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

	insertResult, err := c.CTR.Insert(r.Context(), &p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("created %v document(s)", insertResult))
}

// Update adds changes to personal information using object ID.
func (c *CTRL) Update(w http.ResponseWriter, r *http.Request) {
	var p PersonalData
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
	updateResult, err := c.CTR.Update(r.Context(), &p)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("update %v document(s) successfully", updateResult))
}

// Remove deletes object using url param id which is objectID in database.
func (c *CTRL) Remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idvalue := params["id"]
	result, err := c.CTR.Remove(r.Context(), idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}
	successResponce(w, http.StatusCreated, fmt.Sprintf("deleted %v document(s) successfully", result))
}
