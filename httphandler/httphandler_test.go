package httphandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	var p Provider
	var expected mux.Router

	tt := []struct {
		name     string
		provider Provider
		equal    bool
	}{
		{name: "Not nil interface", provider: p, equal: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := New(tc.provider)
			if tc.equal != assert.IsType(t, &expected, actual) {
				t.Fatalf("not equals interfaces, expected: %v, actual: %v", expected, actual)
			}
		})
	}
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
			// Using the variable on range scope `tc` in function literal (scopelint)
			actualT := receive(&tc.enterT)
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
		expectedT *entity.PersonalData
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
			expectedT: &entity.PersonalData{
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
			// Using the variable on range scope `tc` in function literal (scopelint)
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
			httpError:      errors.New("Bad request"),
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

// func TestHandler(t *testing.T) {
// 	var c *httphandler.Controller

// 	t.Run("", func(t *testing.T) {
// 		srv := httptest.NewServer(httphandler.Handler(c))
// 		defer srv.Close()
// 		fmt.Println(fmt.Sprintf("%s/", srv.URL))
// 		res, err := http.Get(fmt.Sprintf("%s/", srv.URL))
// 		if err != nil {
// 			t.Fatalf("couldn't send GET request: %v", err)
// 		}

// 		defer errors.Wrap(res.Body.Close(), "could not close response body")
// 		_, err = ioutil.ReadAll(res.Body)
// 		if err != nil {
// 			t.Fatalf("couldn't read response body: %v", err)
// 		}
// 		assert.Equal(
// t,
// http.StatusOK,
//  res.StatusCode,
//   fmt.Sprintf("expected status %v; got %v", http.StatusOK, res.StatusCode),
// )
// 	})
// }

// func TestInsert(t *testing.T) {
// 	tt := []struct {
// 		name   string
// 		method string
// 		body   []byte
// 		status int
// 		err    string
// 	}{
// 		{name: "post request", method: "POST", body: nil, status: http.StatusCreated},
// 	}
// 	// addedObject which will used in request body.
// 	addedObject := &network.PersonalData{
// DocumentID: "",
//  Name: "firstName",
//   LastName: "secondName",
//    Phone: "",
//    Email: "",
// 	YearOfBirth: 1980,
// }

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var err error
// 			tc.body, err = json.Marshal(addedObject)
// 			assert.NoError(
// t,
//  err,
//   fmt.Sprintf("couldn't marshal request body: %v", err),
// )

// 			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
// 			assert.NoError(t, err, fmt.Sprintf("couldn't create requset: %v", err))

// 			rec := httptest.NewRecorder()
// 			network.Insert(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			if tc.err != "" {
// 				assert.Equal(
// 	t,
// 	http.StatusBadRequest,
// 	res.StatusCode,
// 	fmt.Sprintf("expected status Bad Request; got: %v", res.StatusCode),
// )
// 				return
// 			}
// 			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode))
// 		})
// 	}
// }

// func TestList(t *testing.T) {
// 	var expectedObject *network.PersonalData

// 	tt := []struct {
// 		name   string
// 		method string
// 		body   []byte
// 		status int
// 		err    string
// 	}{
// 		{name: "get request", method: "GET", body: nil, status: http.StatusOK},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
// 			assert.NoError(t, err, fmt.Sprintf("couldn't create requset: %v", err))

// 			rec := httptest.NewRecorder()
// 			network.List(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			body, err := ioutil.ReadAll(res.Body)
// 			assert.NoError(t, err, fmt.Sprintf("couldn't read response body: %v", err))

// 			if tc.err != "" {
// 				assert.Equal(
// t,
// http.StatusBadRequest,
// res.StatusCode,
//  fmt.Sprintf("expected status Bad Request; got: %v", res.StatusCode),
// )
// 				assert.Equal(
// t,
//  tc.err,
//   string(bytes.TrimSpace(body)),
//    fmt.Sprintf("expected message %q; got %q", tc.err, string(bytes.TrimSpace(body))),
// )
// 				return
// 			}
// 			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode))
// 			assert.NoError(t, json.Unmarshal(body, &expectedObject), err.Error())
// 		})
// 	}
// }

// func TestShowListByID(t *testing.T) {
// 	var expectedObject *network.PersonalData
// 	tt := []struct {
// 		name   string
// 		method string
// 		param  string
// 		value  string
// 		body   []byte
// 		status int
// 		err    string
// 	}{
// 		{name: "get request by id", method: "GET", param: "id", value: "1", body: nil, status: http.StatusOK},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
// 			assert.NoError(t, err, fmt.Sprintf("couldn't create requset: %v", err))

// 			// add url params to request.
// 			q := req.URL.Query()
// 			q.Add(tc.param, tc.value)
// 			req.URL.RawQuery = q.Encode()

// 			rec := httptest.NewRecorder()
// 			network.ByID(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			body, err := ioutil.ReadAll(res.Body)
// 			assert.NoError(t, err, fmt.Sprintf("couldn't read response body: %v", err))

// 			if tc.err != "" {
// 				assert.Equal(t,
//  http.StatusBadRequest,
//   res.StatusCode,
//    fmt.Sprintf("expected status Bad Request; got: %v", res.StatusCode),
// )
// 				assert.Equal(
// 	t,
// 	 tc.err, string(bytes.TrimSpace(body)),
//  fmt.Sprintf(
// 	 "expected message %q; got %q",
// 	  tc.err, string(bytes.TrimSpace(body)),
// 	)
// )
// 				return
// 			}

// 			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode))
// 			assert.NoError(t, json.Unmarshal(body, &expectedObject), err.Error())

// 		})
// 	}
// }

// func TestUpdate(t *testing.T) {
// 	tt := []struct {
// 		name   string
// 		method string
// 		body   []byte
// 		status int
// 		err    string
// 	}{
// 		{name: "update request", method: "PUT", body: nil, status: http.StatusOK},
// 	}

// 	// updateO which will used in request body.
// 	updateO := &network.PersonalData{
// DocumentID: "",
//  Name: "firstName",
//   LastName: "secondName",
//    Phone: "12345678",
//    Email: "test@test.com",
// 	YearOfBirth: 1980,
// }

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var err error
// 			tc.body, err = json.Marshal(updateO)
// 			assert.NoError(
// t,
//  err,
//   fmt.Sprintf("couldn't marshal request body: %v", err),
// )

// 			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
// 			assert.NoError(
// t,
//  err,
//   fmt.Sprintf("couldn't create requset: %v", err),
// )

// 			rec := httptest.NewRecorder()
// 			network.Insert(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			if tc.err != "" {
// 				assert.Equal(
// t,
//  http.StatusBadRequest,
//   res.StatusCode,
//   fmt.Sprintf("expected status Bad Request; got: %v", res.StatusCode),
// )
// 				return
// 			}
// 			assert.Equal(
// t,
// tc.status,
//  res.StatusCode,
//   fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode),
// )
// 		})
// 	}
// }

// func TestRemove(t *testing.T) {
// 	tt := []struct {
// 		name   string
// 		method string
// 		param  string
// 		value  string
// 		body   []byte
// 		status int
// 		err    string
// 	}{
// 		{name: "delete data by id", method: "DELETE", param: "id", value: "1", body: nil, status: http.StatusOK},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req, err := http.NewRequest(tc.method, "localhost:8000/", bytes.NewReader(tc.body))
// 			assert.NoError(
// t,
//  err,
//   fmt.Sprintf("couldn't create requset: %v", err),
// )

// 			// add url params to request.
// 			q := req.URL.Query()
// 			q.Add(tc.param, tc.value)
// 			req.URL.RawQuery = q.Encode()

// 			rec := httptest.NewRecorder()
// 			network.Remove(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()
// 			if tc.err != "" {
// 				assert.Equal(t,
//  http.StatusBadRequest,
//   res.StatusCode,
//   fmt.Sprintf("expected status Bad Request; got: %v",
//   res.StatusCode),
// )
// 				return
// 			}
// 			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("expected status %v; got %v", tc.status, res.StatusCode))
// 		})
// 	}
// }
