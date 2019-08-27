package gateway_test

import (
	"bytes"
	"crud/database"
	"crud/gateway"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPersonalDatas(t *testing.T) {
	var expectedObject *database.PersonalData
	var addedObject = &database.PersonalData{nil, "firstName", "secondName", "", "", 1980}
	tTable := []struct {
		name   string
		method string
		body   []byte
		status int
		err    string
	}{
		{name: "get request", method: "GET", body: nil, status: http.StatusOK},
		{name: "post request", method: "POST", body: nil, status: http.StatusCreated},
	}

	for _, tCase := range tTable {
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.method == "POST" {
				var err error
				tCase.body, err = json.Marshal(addedObject)
				if err != nil {
					t.Fatalf("Couldn't marshal request body: %v", err)
				}

			}
			req, err := http.NewRequest(tCase.method, "localhost:8000/", bytes.NewReader(tCase.body))
			if err != nil {
				t.Fatalf("Couldn't create requset: %v", err)
			}
			rec := httptest.NewRecorder()
			gateway.PersonalDatas(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("couldn't read responce body: %v", err)
			}

			if tCase.err != "" {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status Bad Request; got: %v", res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(body)); msg != tCase.err {
					t.Errorf("expected message %q; got %q", tCase.err, msg)
				}
				return
			}
			if res.StatusCode != tCase.status {
				t.Errorf("expected status %v; got %v", tCase.status, res.StatusCode)
			}

			if tCase.method == "GET" {
				if json.Unmarshal(body, &expectedObject) != nil {
					t.Fatalf("expected body must be %v; got %v", expectedObject, string(body))
				}
			}
		})
	}
}
