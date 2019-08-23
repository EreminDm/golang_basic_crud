package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var objectid = new(primitive.ObjectID)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestInsertPersonalData(t *testing.T) {
	var err error
	gin.SetMode(gin.TestMode)
	var p PersonalData
	p.DocumentID = objectid
	if err != nil {
		t.Errorf(`Couldn't marshal interface %v to json`, &PersonalData{})
		t.Fail()
	}
	p.Name = `Tests`
	p.LastName = `Unit`
	body, err := json.Marshal(p)
	if err != nil {
		t.Errorf(`Couldn't marshal interface %v to json`, &PersonalData{})
		t.Fail()
	}
	// Grab our router
	router := SetupRouter()
	w := performRequest(router, "POST", "/persons/add", body)
	//Assert status code is correct
	assert.Equal(t, http.StatusCreated, w.Code)

}

func TestGetAllPersonalData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var actual interface{}
	var expected PersonalData
	// Grab our router
	router := SetupRouter()

	w := performRequest(router, "GET", "/persons/list", nil)

	//Assert status code is correct
	assert.Equal(t, http.StatusOK, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetPersonalDatabyIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// var actual interface{}
	// var expected PersonalData
	// Grab our router
	router := SetupRouter()
	w := performRequest(router, "GET", "/persons/list/000000000000000000000001", nil)
	//Assert status code is correct
	assert.Equal(t, http.StatusNotFound, w.Code)

	// err := json.Unmarshal([]byte(w.Body.String()), &actual)
	// assert.NoError(t, err)
	// assert.Equal(t, expected, actual)
}

//objectid
func TestGetPersonalDatabyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var actual interface{}
	// var expected PersonalData
	// Grab our router
	router := SetupRouter()
	w := performRequest(router, "GET", "/persons/list/000000000000000000000000", nil)
	//Assert status code is correct
	assert.Equal(t, http.StatusOK, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &actual)
	assert.NoError(t, err)
	// assert.Equal(t, expected, actual)
}
