package mongo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/db/mongo"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	expected := &mongo.Mongodatabase{}
	tt := []struct {
		name          string
		context       context.Context
		connectionURI string
		dbName        string
		err           string
	}{
		{
			name:          "Mongo connection",
			context:       ctx,
			connectionURI: "localhost:27017",
			dbName:        "test",
			err:           "",
		},
		{
			name:          "Mongo wrong Ping",
			context:       ctx,
			connectionURI: "notlocalhost:27017",
			dbName:        "test",
			err:           "couldn't ping database after connection using uri: context deadline exceeded",
		},
		{
			name:          "Mongo wrong Connection",
			context:       ctx,
			connectionURI: "//notlocalhost:27017",
			dbName:        "test",
			err: "couldn't connect to database using uri:" +
				" error parsing uri: must have at least 1 host",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mc, err := mongo.Connect(tc.context, tc.connectionURI, tc.dbName)
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("expected error is not equals; want: %v, got: %v", tc.err, err.Error()),
				)
				return
			}
			assert.NoError(
				t,
				err,
				fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err),
			)
			assert.IsType(t, expected, mc, fmt.Sprintf("expected type %v, actual %v", expected, mc))
		})
	}
}
