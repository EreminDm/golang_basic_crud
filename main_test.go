package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouting(t *testing.T) {
	t.Run("", func(t *testing.T) {
		srv := httptest.NewServer(handler())
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/", srv.URL))
		if err != nil {
			t.Fatalf("couldn't send GET request: %v", err)
		}

		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("couldn't read responce body: %v", err)
		}
		assert.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("expected status %v; got %v", http.StatusOK, res.StatusCode))
	})
}
