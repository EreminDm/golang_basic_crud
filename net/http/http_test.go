package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/net"
	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewHandler(t *testing.T) {
	var p net.Provider
	var expected mux.Router

	tt := []struct {
		name     string
		provider net.Provider
		equal    bool
	}{
		{name: "Not nil interface", provider: p, equal: true},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := New(tc.provider)
			if tc.equal != assert.IsType(t, &expected, actual) {
				t.Fatalf("not equals interfaces, expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mocked read body error")
}

func TestReceive(t *testing.T) {
	tt := []struct {
		name      string
		enterT    entity.PersonalData
		expectedT personalData
		err       error
	}{
		{
			name: "Recive data from entity package type to mongo",
			enterT: entity.PersonalData{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: personalData{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualT := receive(tc.enterT)
			assert.Equal(
				t,
				tc.expectedT,
				actualT,
				fmt.Sprintf("expected type %v, actual %v", tc.expectedT, actualT),
			)
		})
	}
}

func TestTransmit(t *testing.T) {
	tt := []struct {
		name      string
		enterT    personalData
		expectedT entity.PersonalData
		err       error
	}{
		{
			name: "Recive data from entity package type to mongo",
			enterT: personalData{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: entity.PersonalData{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualT := tc.enterT.transmit()
			assert.Equal(
				t,
				tc.expectedT,
				actualT,
				fmt.Sprintf("expected type %v, actual %v", tc.expectedT, actualT),
			)
		})
	}
}

func TestErrRespons(t *testing.T) {
	tt := []struct {
		name           string
		httpStatus     int
		expectedStatus int
		httpError      error
	}{
		{
			name:           "Bad response",
			httpStatus:     http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
			httpError:      errors.New("bad request"),
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// New Recorder creates a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			errRespons(rr, tc.httpStatus, tc.httpError)

			assert.Equal(
				t,
				tc.expectedStatus,
				rr.Code,
				fmt.Sprintf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatus),
			)

			assert.Equal(
				t,
				tc.httpError.Error(),
				rr.Body.String(),
				fmt.Sprintf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.httpError.Error()),
			)
		})
	}
}

func TestSuccessResponce(t *testing.T) {
	tt := []struct {
		name           string
		httpStatus     int
		expectedStatus int
		message        string
	}{
		{
			name:           "Success response",
			httpStatus:     http.StatusOK,
			expectedStatus: http.StatusOK,
			message:        "OK",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// New Recorder creates a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			successResponce(rr, tc.httpStatus, tc.message)

			assert.Equal(
				t,
				tc.expectedStatus,
				rr.Code,
				fmt.Sprintf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatus),
			)

			assert.Equal(
				t,
				tc.message,
				rr.Body.String(),
				fmt.Sprintf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.message),
			)
		})
	}
}

type controllerMockedObject struct {
	mock.Mock
}

// Insert adds a personal information.
func (m *controllerMockedObject) Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

// One returns personal information.
func (m *controllerMockedObject) One(ctx context.Context, id string) (entity.PersonalData, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

// All returns an array of personal information.
func (m *controllerMockedObject) All(ctx context.Context) ([]entity.PersonalData, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.PersonalData), args.Error(1)
}

// Update changes information in collection.
func (m *controllerMockedObject) Update(ctx context.Context, document entity.PersonalData) (int64, error) {
	args := m.Called(ctx, document)
	return int64(args.Int(0)), args.Error(1)
}

// Remove deletes information from collection.
func (m *controllerMockedObject) Remove(ctx context.Context, id string) (int64, error) {
	args := m.Called(ctx, id)
	return int64(1), args.Error(1)
}

func TestInsert(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	tt := []struct {
		name           string
		method         string
		body           []byte
		object         personalData
		expectedObject entity.PersonalData
		expectedError  error
		status         int
		err            string
	}{
		{
			name:   "Success request",
			method: "POST",
			body:   nil,
			object: personalData{
				DocumentID:  "",
				Name:        "firstName",
				LastName:    "secondName",
				Phone:       "",
				Email:       "",
				YearOfBirth: 1980,
			},
			expectedError: nil,
			status:        201,
		},
		{
			name:          "Wrong request",
			method:        "POST",
			body:          []byte("e"),
			object:        personalData{},
			expectedError: nil,
			status:        400,
			err:           "invalid character 'e' looking for beginning of value",
		},
		{
			name:          "HTTP -> read body error",
			method:        "POST",
			body:          []byte("e"),
			object:        personalData{},
			expectedError: nil,
			status:        400,
			err:           "could not read request body\n",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			var err error
			if tc.body == nil {
				tc.body, err = json.Marshal(tc.object)
				assert.NoError(
					t,
					err,
					fmt.Sprintf("couldn't marshal request body: %v", err),
				)
			}
			var req *http.Request
			if tc.name == "HTTP -> read body error" {
				req = httptest.NewRequest(tc.method, "http://localhost:8000/", errReader(0))
			} else {
				req, err = http.NewRequest(tc.method, "http://localhost:8000/", bytes.NewReader(tc.body))
				assert.NoError(t, err, fmt.Sprintf("couldn't create requset: %v", err))
			}
			rec := httptest.NewRecorder()

			ctr.On("Insert", mock.Anything, tc.object.transmit()).Return(tc.expectedObject, tc.expectedError)
			c.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if tc.err != "" {
				b, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, fmt.Sprintf("could not read response body: %v", err))
				bString := string(b)
				assert.Equal(t,
					tc.err,
					bString,
					"not equals",
				)
			}

			if tc.expectedError != nil {
				b, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, fmt.Sprintf("could not read response body: %v", err))
				bString := string(b)
				assert.Equal(t,
					tc.expectedError,
					bString,
					"not equals",
				)
			}

			assert.Equal(t,
				tc.status,
				res.StatusCode,
				fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode),
			)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	tt := []struct {
		name          string
		method        string
		body          []byte
		object        personalData
		expectedError error
		expectedCount int
		status        int
		err           string
	}{
		{
			name:   "HTTP -> Success request",
			method: "PUT",
			body:   nil,
			object: personalData{
				DocumentID:  "",
				Name:        "firstName",
				LastName:    "secondName",
				Phone:       "",
				Email:       "",
				YearOfBirth: 1980,
			},
			expectedError: nil,
			expectedCount: 1,
			status:        201,
		},
		{
			name:          "HTTP -> Wrong request",
			method:        "PUT",
			body:          []byte("e"),
			object:        personalData{},
			expectedError: nil,
			expectedCount: 0,
			status:        400,
			err:           "invalid character 'e' looking for beginning of value",
		},
		{
			name:          "HTTP -> wrong read body error",
			method:        "PUT",
			body:          []byte("e"),
			object:        personalData{},
			expectedError: nil,
			expectedCount: 0,
			status:        400,
			err:           "mocked read body error",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			var err error
			if tc.body == nil {
				tc.body, err = json.Marshal(tc.object)
				assert.NoError(
					t,
					err,
					fmt.Sprintf("couldn't marshal request body: %v", err),
				)
			}
			var req *http.Request

			if tc.name == "HTTP -> wrong read body error" {
				req = httptest.NewRequest(tc.method, "http://localhost:8000/", errReader(0))
			} else {
				req, err = http.NewRequest(tc.method, "http://localhost:8000/", bytes.NewReader(tc.body))
				assert.NoError(t, err, fmt.Sprintf("couldn't create requset: %v", err))
			}

			rec := httptest.NewRecorder()

			ctr.On("Update", mock.Anything, tc.object.transmit()).Return(tc.expectedCount, tc.expectedError)
			c.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if tc.err != "" {
				b, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, fmt.Sprintf("could not read response body: %v", err))
				bString := string(b)
				assert.Equal(t,
					tc.err,
					bString,
					"not equals",
				)
				return
			}

			if tc.expectedError != nil {
				b, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, fmt.Sprintf("could not read response body: %v", err))
				bString := string(b)
				assert.Equal(t,
					tc.expectedError.Error(),
					bString,
					"not equals",
				)
				return
			}

			assert.Equal(t,
				tc.status,
				res.StatusCode,
				fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode),
			)
		})
	}
}

func TestList(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	server := httptest.NewServer(c)
	defer server.Close()
	tt := []struct {
		name           string
		schema         string
		object         personalData
		status         int
		expectedObject []entity.PersonalData
		expectedError  error
		err            string
	}{
		{name: "get request",
			schema: `{
				"type":"array",
				"items":{
					"type":"object",
					"properties": {
						"documentID":{"type":"string"},
						"name":{"type":"string"},
						"lastName":{"type":"string"},
						"phone":{"type":"string"},
						"email":{"type":"string"},
						"yearOfBirth":{"type":"integer"}
					},
					"required": ["name", "description"]
				}	 
		 }
		 `,
			object: personalData{
				DocumentID:  "",
				Name:        "firstName",
				LastName:    "secondName",
				Phone:       "",
				Email:       "",
				YearOfBirth: 1980,
			},
			status: http.StatusOK,
		},
		// {
		// 	name:           "WRONG reques",
		// 	status:         http.StatusInternalServerError,
		// 	expectedObject: nil,
		// 	expectedError:  errors.New("404 Not Found"),
		// 	err:            "404 Not Found",
		// },
	}

	for _, tc := range tt {
		tc := tc
		// create httpexpect instance
		ctr.On("All", mock.Anything).Return(tc.expectedObject, tc.expectedError)
		t.Run(tc.name, func(t *testing.T) {

			e := httpexpect.New(t, server.URL)
			if tc.err != "" {
				res := e.GET("/").
					Expect().
					Status(tc.status)
				res.Body().Equal(tc.err)
				return
			}
			e.GET("/").
				Expect().
				Status(http.StatusOK).JSON().Array().Empty()

			repos := e.GET("/").
				Expect().Status(http.StatusOK).JSON().Array()
			repos.Schema(tc.schema)

		})
	}
}

func TestByID(t *testing.T) {
	var expectedObject entity.PersonalData
	ctr := new(controllerMockedObject)
	c := New(ctr)
	server := httptest.NewServer(c)
	defer server.Close()
	oid := primitive.NewObjectID().Hex()

	tt := []struct {
		name          string
		schema        string
		param         string
		status        int
		expectedError error
		err           string
	}{
		{
			name: "get request by id",
			schema: `{
			"type":"object",
					"properties": {
						"documentID":{"type":"string"},
						"name":{"type":"string"},
						"lastName":{"type":"string"},
						"phone":{"type":"string"},
						"email":{"type":"string"},
						"yearOfBirth":{"type":"integer"}
					},
					"required": ["name", "description"]
		}`,

			param:  oid,
			status: http.StatusOK,
			err:    "",
		},
		{
			name: "WRONG request by id",
			schema: `{
				
			}`,
			param:         "oid",
			status:        http.StatusInternalServerError,
			expectedError: errors.New("404 Not Found"),
			err:           "404 Not Found",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("One", mock.Anything, tc.param).Return(expectedObject, tc.expectedError)
			e := httpexpect.New(t, server.URL)
			if tc.err != "" {
				res := e.GET("/" + tc.param).
					Expect().
					Status(tc.status)
				res.Body().Equal(tc.err)
				return
			}
			e.GET("/" + tc.param).
				Expect().
				Status(tc.status)
		})
	}
}

func TestRemove(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	server := httptest.NewServer(c)
	defer server.Close()
	oid := primitive.NewObjectID().Hex()
	tt := []struct {
		name          string
		param         string
		value         int64
		status        int
		expectedError error
		err           string
	}{
		{name: "Sucsecc deleting", param: oid, value: 1, status: http.StatusOK, expectedError: nil, err: ""},
		{
			name: "False deleting", param: "oid", value: 0, status: http.StatusInternalServerError,
			expectedError: errors.New("404 Not Found"), err: "404 Not Found",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("Remove", mock.Anything, tc.param).Return(tc.value, tc.expectedError)
			e := httpexpect.New(t, server.URL)
			if tc.err != "" {
				res := e.DELETE("/" + tc.param).
					Expect().Status(tc.status)
				res.Body().Equal(tc.err)
				//assert.Equal(t, tc.err, res.Body(), "not equals, want %v, got %v", tc.err, res.Body().Equal())
				return
			}
			e.DELETE("/" + tc.param).
				Expect().
				Status(tc.status)

		})
	}
}
