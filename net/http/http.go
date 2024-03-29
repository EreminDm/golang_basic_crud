package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/net"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Controller describes controller implementation.
type Controller struct {
	CTR net.Provider
}

// personalData is personal data filds description.
type personalData struct {
	DocumentID  string `json:"id,omitempty"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	YearOfBirth int    `json:"yaerOfBirth,omitempty"`
}

// New returns new controller provider
func New(c net.Provider) http.Handler {
	return Handler(&Controller{
		CTR: c,
	})
}

// receive returns http package personal data construction.
func receive(ep entity.PersonalData) personalData {
	return personalData{
		DocumentID:  ep.DocumentID,
		Name:        ep.Name,
		LastName:    ep.LastName,
		Phone:       ep.Phone,
		Email:       ep.Email,
		YearOfBirth: ep.YearOfBirth,
	}
}

// transmit returns entity data construction.
func (p *personalData) transmit() entity.PersonalData {
	return entity.PersonalData{
		DocumentID:  p.DocumentID,
		Name:        p.Name,
		LastName:    p.LastName,
		Phone:       p.Phone,
		Email:       p.Email,
		YearOfBirth: p.YearOfBirth,
	}
}

// Handler is function for routing map navigation.
func Handler(c *Controller) http.Handler {
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

// errRespons returns error to response.
func errRespons(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	w.WriteHeader(code)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Panic(errors.Wrap(err, "could not write bytes"))
	}
}

// successResponce returns success response.
func successResponce(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Panic(errors.Wrap(err, "could not write bytes"))
	}
}

// ByID returns personal data list by id.
func (c *Controller) ByID(w http.ResponseWriter, r *http.Request) {
	// reads parameters from url.
	params := mux.Vars(r)
	id := params["id"]

	// makes request to controller.
	ep, err := c.CTR.One(r.Context(), id)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// marshalls data from response.
	res, err := json.Marshal(receive(ep))
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// makes response to API client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		log.Panic(errors.Wrap(err, "could not write bytes"))
	}
}

// List returns a list of personaldata.
func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	var pd = make([]personalData, 0)
	// makes request to controller.
	usrs, err := c.CTR.All(r.Context())
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// converts data to http layout.
	for _, ep := range usrs {
		var epp = entity.PersonalData{
			DocumentID:  ep.DocumentID,
			Name:        ep.Name,
			LastName:    ep.LastName,
			Phone:       ep.Phone,
			Email:       ep.Email,
			YearOfBirth: ep.YearOfBirth,
		}
		received := receive(epp)
		pd = append(pd, received)
	}

	// marshalls data to hson byte format.
	res, err := json.Marshal(pd)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// makes response to API client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		log.Panic(errors.Wrap(err, "could not write bytes"))
	}
}

// Insert creates Personal Data by preparing to insert new data to database.
func (c *Controller) Insert(w http.ResponseWriter, r *http.Request) {
	// reads request body.
	var p personalData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	// unmarshal to json struct.
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}

	// makes insert request to controller.
	result, err := c.CTR.Insert(r.Context(), p.transmit())
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	epres := receive(result)
	res, err := json.Marshal(epres)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// makes response to API client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(res)
	if err != nil {
		log.Panic(errors.Wrap(err, "could not write bytes"))
	}
}

// Update adds changes to personal information using object ID.
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var p personalData
	// reading request body information.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}

	// unmarshal body to personalData object.
	err = json.Unmarshal(body, &p)
	if err != nil {
		errRespons(w, http.StatusBadRequest, err)
		return
	}

	// makes update request to controller.
	updateResult, err := c.CTR.Update(r.Context(), p.transmit())
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// makes response to API client.
	successResponce(w, http.StatusCreated, fmt.Sprintf("update %v document(s) successfully", updateResult))
}

// Remove deletes object using url param id which is objectID in database.
func (c *Controller) Remove(w http.ResponseWriter, r *http.Request) {
	// reads parameters from url.
	params := mux.Vars(r)
	idvalue := params["id"]

	// makes remove request to controller.
	result, err := c.CTR.Remove(r.Context(), idvalue)
	if err != nil {
		errRespons(w, http.StatusInternalServerError, err)
		return
	}

	// makes response to API client.
	successResponce(w, http.StatusOK, fmt.Sprintf("deleted %v document(s) successfully", result))
}
