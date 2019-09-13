package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	conn, db := envf()
	assert.Equal(t, "localhost:27017", conn, "Expected is not equals to actual, got %s", conn)
	assert.Equal(t, "information", db, "Expected is not equals to actual, got %s", db)
}

func TestMain(m *testing.M) {
	code := m.Run()
	fmt.Print("Code is", code)
	os.Exit(code)
}

// func TestMain(t *testing.T) {
// 	expectedHandler := new(mux.Router)
// 	expectedController := new(controller.Personal)
// 	expecteddatabaseProvier := new(mongo.Mongodatabase)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	tt := []struct {
// 		connectionURI string
// 		databaseName  string
// 	}{
// 		{connectionURI: "192.168.99.100:27017", databaseName: "test"},
// 	}

// 	for _, tc := range tt {
// 		m, err := mongo.Connect(ctx, tc.connectionURI, tc.databaseName)
// 		assert.NoError(t, err, "could not connect to mongo")
// 		assert.IsType(t, expecteddatabaseProvier, m, "Database provider type is not equals")
// 		c := controller.New(m)
// 		assert.IsType(t, expectedController, c, "Controllers type is not equals")
// 		h := httphandler.New(c)
// 		assert.IsType(t, expectedHandler, h, "Handler type is not equals")

// 	}
// }
