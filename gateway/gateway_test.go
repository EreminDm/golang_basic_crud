package gateway_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EreminDm/golang_basic_crud/database"
	"github.com/EreminDm/golang_basic_crud/gateway"
	"github.com/stretchr/testify/assert"
)

func TestShowList(t *testing.T) {
	var expectedObject *database.PersonalData
	var mockCollecton *database.Collection
	tt := []struct {
		name   string
		method string
		body   []byte
		status int
		err    string
	}{
		{name: "get request", method: "GET", body: nil, status: http.StatusOK},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
			assert.Error(t, err, fmt.Sprintf("couldn't create requset: %v", err))

			rec := httptest.NewRecorder()
			gateway.ShowList(rec, req, mockCollecton)

			res := rec.Result()
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			assert.Error(t, err, fmt.Sprintf("couldn't read responce body: %v", err))

			if tc.err != "" {
				assert.Equal(t, http.StatusBadRequest, res.StatusCode, fmt.Sprintf("expected status Bad Request; got: %v", res.StatusCode))
				assert.Equal(t, tc.err, string(bytes.TrimSpace(body)), fmt.Sprintf("expected message %q; got %q", tc.err, string(bytes.TrimSpace(body))))
				return
			}
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode))
			assert.Error(t, json.Unmarshal(body, &expectedObject), err.Error())
		})
	}
}

// for post test
// var addedObject = &database.PersonalData{DocumentID: nil, Name: "firstName", LastName: "secondName", Phone: "", Email: "", YearOfBirth: 1980}
//tt: {name: "post request", method: "POST", body: nil, status: http.StatusCreated},
// if tc.method == "POST" {
// 	var err error
// 	tc.body, err = json.Marshal(addedObject)
// 	if err != nil {
// 		t.Fatalf("couldn't marshal request body: %v", err)
// 	}

// }
